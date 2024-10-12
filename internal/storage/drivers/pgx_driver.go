package drivers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	storageaccrual "github.com/mi4r/gophermart/internal/storage/accrual"
	storagedefault "github.com/mi4r/gophermart/internal/storage/default"
	storagemart "github.com/mi4r/gophermart/internal/storage/gophermart"
)

const (
	migrDefaultPath = "default"
)

var (
	errNotFoundOrder = errors.New("order not found")
)

type pgxDriver struct {
	dbURL    string
	connPool *pgxpool.Pool
	isTest   bool
}

func (d *pgxDriver) exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	slog.Debug(sql, slog.Any("args", args))
	return d.connPool.Exec(ctx, sql, args...)
}

func (d *pgxDriver) queryRows(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	slog.Debug(sql, slog.Any("args", args))
	return d.connPool.Query(ctx, sql, args...)
}

func (d *pgxDriver) queryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	slog.Debug(sql, slog.Any("args", args))
	return d.connPool.QueryRow(ctx, sql, args...)
}

func NewPgxDriver(path string) *pgxDriver {
	return &pgxDriver{
		dbURL: path,
	}
}

func (d *pgxDriver) Open(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, d.dbURL)
	if err != nil {
		return err
	}
	d.connPool = pool

	// Try connect
	if err := d.Ping(); err != nil {
		return err
	}

	return nil
}

func (d *pgxDriver) Migrate(migrDirName string) {
	// Try auto-migration
	if err := d.autoDefaultMigrate(); err != nil {
		slog.Warn("migration error", slog.String("err", err.Error()))
	} else {
		slog.Debug("migration OK")
	}
}

func (d *pgxDriver) Ping() error {
	return d.connPool.Ping(context.Background())
}

func (d *pgxDriver) autoDefaultMigrate() error {
	mpath, err := filepath.Abs(
		filepath.Join("internal", "storage", migrDefaultPath, "migrations"))
	if err != nil {
		return err
	}

	slog.Debug("migration init", slog.String("path", mpath))
	migr, err := migrate.New(
		fmt.Sprintf("file://%s", mpath),
		d.dbURL,
	)
	if err != nil {
		return err
	}
	return migr.Up()
}

func (d *pgxDriver) Close() {
	d.connPool.Close()
}

func (d *pgxDriver) UserCreate(ctx context.Context, user storagemart.User) error {
	_, err := d.exec(ctx, `
	INSERT INTO users (login, password)
	VALUES ($1, $2)
	`, user.Login, user.Password,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *pgxDriver) UserReadOne(ctx context.Context, login string) (storagemart.User, error) {
	var user storagemart.User
	if err := d.queryRow(ctx, `
		SELECT login, password, current, withdrawn FROM users WHERE login=$1
	`, login).Scan(&user.Login, &user.Password, &user.Current, &user.Withdrawn); err != nil {
		return user, err
	}
	return user, nil
}

func (d *pgxDriver) UserOrderCreate(ctx context.Context, login, number string) error {
	_, err := d.exec(ctx, `
	INSERT INTO user_orders (number, user_login)
	VALUES ($1, $2)
	`, number, login,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *pgxDriver) UserOrderReadOne(ctx context.Context, number string) (storagemart.Order, error) {
	var o storagemart.Order
	if err := d.queryRow(ctx, `
	SELECT number, status, accrual, uploaded_at, user_login
	FROM user_orders
		WHERE number = $1
		LIMIT 1
	`, number).Scan(
		&o.Number, &o.Status, &o.Accrual,
		&o.UploadedAt, &o.UserLogin,
	); err != nil {
		return o, err
	}
	return o, nil
}

func (d *pgxDriver) UserOrdersReadByLogin(ctx context.Context, login string) ([]storagemart.Order, error) {
	var orders []storagemart.Order
	rows, err := d.queryRows(ctx, `
	SELECT number, status, accrual, uploaded_at, user_login
	FROM user_orders
		WHERE user_login = $1
		ORDER BY uploaded_at ASC
	`, login)
	if err != nil {
		return orders, err
	}
	defer rows.Close()

	for rows.Next() {
		var o storagemart.Order
		if err := rows.Scan(
			&o.Number, &o.Status, &o.Accrual,
			&o.UploadedAt, &o.UserLogin,
		); err != nil {
			slog.Error("scan error", slog.String("err", err.Error()))
			return []storagemart.Order{}, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func (d *pgxDriver) UserOrderReadAllNumbers(ctx context.Context) ([]string, error) {
	var orders []string
	rows, err := d.queryRows(ctx, `
	SELECT number FROM user_orders
		WHERE status != 'INVALID' AND status != 'PROCESSED'`)
	if err != nil {
		return orders, err
	}
	defer rows.Close()

	for rows.Next() {
		var o string
		if err := rows.Scan(
			&o,
		); err != nil {
			slog.Error("scan error", slog.String("err", err.Error()))
			return []string{}, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func (d *pgxDriver) UserOrderUpdateStatus(ctx context.Context, number string, status storagedefault.OrderStatus) error {
	if _, err := d.exec(ctx, `
	UPDATE user_orders SET status = $1 
		WHERE number = $2`, status, number); err != nil {
		return err
	}
	return nil
}

func (d *pgxDriver) UserOrderUpdateAll(ctx context.Context, orders []storagedefault.Order) error {
	tx, err := d.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, o := range orders {
		if _, err := tx.Exec(ctx, `
		UPDATE user_orders SET status = $1, accrual=$2 processed_at = NOW() WHERE number = $3
		`, o.Status, o.Accrual, o.Number); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (d *pgxDriver) RewardCreate(ctx context.Context, r storageaccrual.Reward) error {
	_, err := d.exec(context.Background(), `
	INSERT INTO rewards (match, reward, reward_type)
	VALUES ($1, $2, $3)
	`, r.Match, r.Reward, r.RewardType,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *pgxDriver) RewardReadAll(ctx context.Context) ([]storageaccrual.Reward, error) {
	var rewards []storageaccrual.Reward
	rows, err := d.queryRows(ctx, "SELECT match, reward, reward_type FROM rewards")
	if err != nil {
		return rewards, err
	}
	defer rows.Close()
	for rows.Next() {
		var r storageaccrual.Reward
		if err := rows.Scan(
			&r.Match, &r.Reward, &r.RewardType,
		); err != nil {
			slog.Error("scan error", slog.String("err", err.Error()))
			return rewards, err
		}
		rewards = append(rewards, r)
	}
	return rewards, nil
}

func (d *pgxDriver) OrderRegCreate(ctx context.Context, o storageaccrual.Order) error {
	tx, err := d.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	var orderID int64
	defer tx.Rollback(ctx)
	if err := tx.QueryRow(ctx, `
	INSERT INTO orders (order_number, status) VALUES ($1, $2) RETURNING id`, o.Order, storagedefault.StatusRegistered).Scan(&orderID); err != nil {
		return err
	}
	slog.Debug("order id is fetch", slog.Int64("id", orderID), slog.String("order", o.Order))

	sqlScriptCreateGoods := `INSERT INTO goods (description, price) VALUES ($1, $2) RETURNING id`
	sqlScriptGoodInOrder := `INSERT INTO order_goods (order_id, good_id) VALUES ($1, $2)`

	for _, good := range o.Goods {
		slog.Debug("add good of order",
			slog.String("description", good.Description),
			slog.Float64("price", good.Price),
		)
		var goodID int64
		if err := tx.QueryRow(
			ctx, sqlScriptCreateGoods, good.Description, good.Price).
			Scan(&goodID); err != nil {
			return err
		}
		if _, err := tx.Exec(
			ctx, sqlScriptGoodInOrder, orderID, goodID); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (d *pgxDriver) OrderRegReadOne(ctx context.Context, number string) (storagedefault.Order, error) {
	var o storagedefault.Order
	if err := d.queryRow(ctx, `
	SELECT order_number, status, accrual
	FROM orders
		WHERE order_number = $1
		LIMIT 1
	`, number).Scan(
		&o.Number, &o.Status, &o.Accrual,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return o, errNotFoundOrder
		}
		return o, err
	}
	return o, nil
}

func (d *pgxDriver) OrderRegUpdateStatus(ctx context.Context, status storagedefault.OrderStatus, number string) error {
	if _, err := d.exec(ctx, `
	UPDATE orders SET status=$1 WHERE order_number=$2
	`, status, number); err != nil {
		return err
	}
	return nil
}

func (d *pgxDriver) OrderRegUpdateOne(ctx context.Context, order storagedefault.Order) error {
	if _, err := d.exec(ctx, `
	UPDATE orders SET status=$1, accrual=$2 WHERE order_number=$3
	`, order.Status, order.Accrual, order.Number); err != nil {
		return err
	}
	return nil
}

func (d *pgxDriver) WithdrawBalance(ctx context.Context, login, order string, sum, curBalance float64) error {
	tx, err := d.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queryUpdate := `UPDATE users SET current = current - $1, withdrawn = withdrawn + $1 WHERE login = $2;`
	_, err = tx.Exec(ctx, queryUpdate, sum, login)
	if err != nil {
		return err
	}

	queryInsertOrder := `INSERT INTO orders (number, user_login, sum, is_withdrawn, processed_at) VALUES ($1, $2, $3, $4, $5);`
	_, err = tx.Exec(ctx, queryInsertOrder, order, login, sum, true, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (d *pgxDriver) GetUserWithdrawals(ctx context.Context, login string) ([]storagemart.Order, error) {
	query := `SELECT number, sum, processed_at
			FROM orders
			WHERE user_login = $1 AND is_withdrawn = $2
			ORDER BY processed_at ASC`

	rows, err := d.queryRows(ctx, query, login, true)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []storagemart.Order
	for rows.Next() {
		var w storagemart.Order
		if err := rows.Scan(&w.Number, &w.Sum, &w.ProcessedAt); err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, w)
	}

	return withdrawals, nil
}

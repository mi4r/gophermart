package drivers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	storageaccrual "github.com/mi4r/gophermart/internal/storage/accrual"
	storagemart "github.com/mi4r/gophermart/internal/storage/gophermart"
)

const (
	migrDefaultPath = "migrations"
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

func (d *pgxDriver) Open() error {
	pool, err := pgxpool.New(context.Background(), d.dbURL)
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
	if err := d.autoMigrate(migrDirName, d.isTest); err != nil {
		slog.Warn("migration error", slog.String("err", err.Error()))
	} else {
		slog.Debug("migration OK")
	}
}

func (d *pgxDriver) Ping() error {
	return d.connPool.Ping(context.Background())
}

// Автоматическая миграция базы. Думаю прикрутить ключ при запуске
func (d *pgxDriver) autoMigrate(migrDirName string, isTest bool) (err error) {

	mpath := migrDefaultPath
	if !isTest {
		mpath, err = filepath.Abs(
			filepath.Join("internal", "storage", migrDirName, "migrations"))
		if err != nil {
			return err
		}
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

func (d *pgxDriver) UserCreate(user storagemart.User) error {
	_, err := d.exec(context.Background(), `
	INSERT INTO users (login, password)
	VALUES ($1, $2)
	`, user.Login, user.Password,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *pgxDriver) UserReadOne(login string) (storagemart.User, error) {
	ctx := context.Background()
	var user storagemart.User
	if err := d.queryRow(ctx, `
		SELECT login, password, current, withdrawn FROM users WHERE login=$1
	`, login).Scan(&user.Login, &user.Password, &user.Current, &user.Withdrawn); err != nil {
		return user, err
	}
	return user, nil
}

func (d *pgxDriver) UserReadAll() ([]storagemart.User, error) {
	var users []storagemart.User
	ctx := context.Background()
	rows, err := d.queryRows(ctx, `
	SELECT login, password, current, withdrawn FROM users
	`)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	var errs []error
	for rows.Next() {
		var user storagemart.User
		if err := rows.Scan(&user.Login, &user.Password, &user.Current, &user.Withdrawn); err != nil {
			slog.Error("scan error", slog.String("err", err.Error()))
			errs = append(errs, err)
		}
		users = append(users, user)
	}

	return users, errors.Join(errs...)
}

func (d *pgxDriver) OrderCreate(login, number string) error {
	_, err := d.exec(context.Background(), `
	INSERT INTO orders (number, user_login)
	VALUES ($1, $2)
	`, number, login,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *pgxDriver) OrderReadOne(number string) (storagemart.Order, error) {
	ctx := context.Background()
	var o storagemart.Order
	if err := d.queryRow(ctx, `
	SELECT number, status, accrual, uploaded_at, user_login
	FROM orders
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

func (d *pgxDriver) OrdersReadByLogin(login string) ([]storagemart.Order, error) {
	var orders []storagemart.Order
	ctx := context.Background()
	rows, err := d.queryRows(ctx, `
	SELECT number, status, accrual, uploaded_at, user_login
	FROM orders
		WHERE user_login = $1
		ORDER BY uploaded_at ASC
	`, login)
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	var errs []error
	for rows.Next() {
		var o storagemart.Order
		if err := rows.Scan(
			&o.Number, &o.Status, &o.Accrual,
			&o.UploadedAt, &o.UserLogin,
		); err != nil {
			slog.Error("scan error", slog.String("err", err.Error()))
			errs = append(errs, err)
		}
		orders = append(orders, o)
	}

	return orders, errors.Join(errs...)
}

func (d *pgxDriver) RewardCreate(r storageaccrual.Reward) error {
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

func (d *pgxDriver) OrderRegCreate(o storageaccrual.Order) error {
	ctx := context.Background()
	tx, err := d.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	var orderID int64
	defer tx.Rollback(ctx)
	if err := tx.QueryRow(ctx, `
	INSERT INTO orders (order_number) VALUES ($1) RETURNING id`, o.Order).Scan(&orderID); err != nil {
		return err
	}
	slog.Debug("order id is fetch", slog.Int64("id", orderID), slog.String("order", o.Order))

	sqlScriptCreateGoods := `INSERT INTO goods (description, price) VALUES ($1, $2) RETURNING id`
	sqlScriptGoodInOrder := `INSERT INTO order_goods (order_id, good_id) VALUES ($1, $2)`
	var errs []error
	for _, good := range o.Goods {
		slog.Debug("add good of order",
			slog.String("description", good.Description),
			slog.Float64("price", good.Price),
		)
		var goodID int64
		if err := tx.QueryRow(
			ctx, sqlScriptCreateGoods, good.Description, good.Price).
			Scan(&goodID); err != nil {
			errs = append(errs, err)
		}
		if _, err := tx.Exec(
			ctx, sqlScriptGoodInOrder, orderID, goodID); err != nil {
			errs = append(errs, err)
		}
	}
	errs = append(errs, tx.Commit(ctx))
	return errors.Join(errs...)
}

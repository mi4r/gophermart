package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	MigrDirNameTest string = "migrations"
	MigrDirNameProd string = "internal/storage/migrations"
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

	// Try auto-migration
	if err := d.autoMigrate(d.isTest); err != nil {
		slog.Warn("migration error", slog.String("err", err.Error()))
	}
	return nil
}

func (d *pgxDriver) Ping() error {
	return d.connPool.Ping(context.Background())
}

// Автоматическая миграция базы. Думаю прикрутить ключ при запуске
func (d *pgxDriver) autoMigrate(isTest bool) error {
	curDirAbs, err := os.Getwd()
	if err != nil {
		return err
	}
	MigrDirName := MigrDirNameProd
	if isTest {
		MigrDirName = MigrDirNameTest
	}
	migrDirAbsPath := path.Join(curDirAbs, MigrDirName)
	slog.Debug("migration init", slog.String("path", migrDirAbsPath))
	migr, err := migrate.New(
		fmt.Sprintf("file://%s", migrDirAbsPath),
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

func (d *pgxDriver) UserCreate(user User) error {
	_, err := d.exec(context.Background(),
		`
	INSERT INTO users (login, password)
	VALUES ($1, $2)
	`, user.Login, user.Password,
	)
	if err != nil {
		return err
	}
	return nil
}

func (d *pgxDriver) UserReadOne(login string) (User, error) {
	ctx := context.Background()
	var user User
	if err := d.queryRow(ctx, `
		SELECT login, password, balance FROM users WHERE login=$1
		`, login).Scan(&user.Login, &user.Password, &user.Balance); err != nil {
		return user, err
	}
	return user, nil
}

func (d *pgxDriver) UserReadAll() ([]User, error) {
	var users []User
	ctx := context.Background()
	rows, err := d.queryRows(ctx, `
	SELECT login, password, balance FROM users
	`)
	if err != nil {
		return users, err
	}
	defer rows.Close()
	var errs []error
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Login, &user.Password, &user.Balance); err != nil {
			slog.Error("scan error", slog.String("err", err.Error()))
			errs = append(errs, err)
		}
		users = append(users, user)
	}

	return users, errors.Join(errs...)
}

func (d *pgxDriver) UserUpdate(user User) error {
	return nil
}

func (d *pgxDriver) UserDelete(user User) error {
	return nil
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

func (d *pgxDriver) OrderReadOne(number string) (Order, error) {
	ctx := context.Background()
	var o Order
	var id int64
	// if err := d.queryRow(ctx, `
	// 	SELECT number, status, accrual, uploaded_at, user_login FROM orders WHERE number=$1
	// 	`, number).Scan(
	if err := d.queryRow(ctx, `
		SELECT * FROM orders WHERE number=$1
	`, number).Scan(&id,
		&o.Number, &o.Status, &o.Accrual,
		&o.UploadedAt, &o.UserLogin,
	); err != nil {
		return o, err
	}
	return o, nil
}

func (d *pgxDriver) OrdersReadByLogin(login string) ([]Order, error) {
	var orders []Order
	ctx := context.Background()
	rows, err := d.queryRows(ctx, `
	SELECT * FROM orders
	`)
	if err != nil {
		return orders, err
	}
	defer rows.Close()
	var errs []error
	for rows.Next() {
		var o Order
		var id int64
		if err := rows.Scan(&id,
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

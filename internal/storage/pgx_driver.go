package storage

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxDriver struct {
	dbURL    string
	connPool *pgxpool.Pool
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
	return nil
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
			slog.Error("user not found", slog.Any("user", user))
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

package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ginhub.com/Aller101/sso/internal/domain/models"
	"ginhub.com/Aller101/sso/internal/storage"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return &Storage{db: db}, err
}

func (s *Storage) UserSave(ctx context.Context, email string, passHash []byte) (uid int64, err error) {
	const op = "storage.sqlite.UserSave"

	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %s", op, err)
	}
	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %s", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %s", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %s", op, err)
	}
	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email=?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %s", op, err)
	}

	var user models.User
	if err := stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %s", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %s", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id=?")
	if err != nil {
		return false, fmt.Errorf("%s: %s", op, err)
	}

	var isAdmin bool

	if err := stmt.QueryRowContext(ctx, userID).Scan(&isAdmin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %s", op, storage.ErrAppNotFound)
		}
		return false, fmt.Errorf("%s: %s", op, err)
	}
	return isAdmin, nil

}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.sqlite.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id=?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %s", op, err)
	}
	var app models.App

	err = stmt.QueryRowContext(ctx, appID).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %s", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %s", op, err)
	}

	return app, nil

}

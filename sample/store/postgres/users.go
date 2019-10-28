package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/guzenok/go-sqltest/sample/store"
)

func (s *postgresStore) sqlCreateUser(u *store.User) (string, []interface{}, error) {
	return s.db.BindNamed(`
INSERT INTO users (id, login, password, is_super) 
VALUES (:id, :login, :password, :is_super)
RETURNING created_at;`,
		u)
}

func (s *postgresStore) CreateUser(ctx context.Context, in *store.User) (out *store.User, err error) {
	out = new(store.User)
	*out = *in

	err = s.ExecuteInTransaction(ctx, func(tx *sqlx.Tx) error {
		query, args, err := s.sqlCreateUser(in)
		if err != nil {
			return err
		}

		return tx.GetContext(ctx, out, query, args...)
	})

	return
}

func (s *postgresStore) sqlDeleteUser(u *store.User) (string, []interface{}, error) {
	return s.db.BindNamed(`
DELETE FROM users 
WHERE id = :id;`,
		u)
}

func (s *postgresStore) DeleteUser(ctx context.Context, id int32) (err error) {
	user := &store.User{
		ID: id,
	}

	err = s.ExecuteInTransaction(ctx, func(tx *sqlx.Tx) error {
		query, args, err := s.sqlDeleteUser(user)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, query, args...)
		return err
	})

	return
}

func (s *postgresStore) sqlGetUserByID(u *store.User) (string, []interface{}, error) {
	return s.db.BindNamed(`
SELECT login, password, created_at, is_super
FROM users
WHERE id = :id
LIMIT 1;`,
		u)
}

func (s *postgresStore) GetUserByID(ctx context.Context, id int32) (out *store.User, err error) {
	user := &store.User{
		ID: id,
	}

	query, args, err := s.sqlGetUserByID(user)
	if err != nil {
		return
	}

	err = s.db.GetContext(ctx, user, query, args...)
	if err != nil {
		return
	}

	out = user
	return
}

func (s *postgresStore) sqlGetUserByLogin(u *store.User) (string, []interface{}, error) {
	return s.db.BindNamed(`
SELECT id, password, created_at, is_super
FROM users
WHERE login = :login
LIMIT 1;`,
		u)
}

func (s *postgresStore) GetUserByLogin(ctx context.Context, login string) (out *store.User, err error) {
	user := &store.User{
		Login: login,
	}

	query, args, err := s.sqlGetUserByLogin(user)
	if err != nil {
		return
	}

	err = s.db.GetContext(ctx, user, query, args...)
	if err != nil {
		return
	}

	out = user
	return
}

func (s *postgresStore) sqlSetPassword(u *store.User) (string, []interface{}, error) {
	return s.db.BindNamed(`
UPDATE users SET password = :password 
WHERE id = :id;`,
		u)
}

func (s *postgresStore) SetPassword(ctx context.Context, userID int32, password string) (err error) {
	user := &store.User{
		ID:       userID,
		Password: password,
	}

	err = s.ExecuteInTransaction(ctx, func(tx *sqlx.Tx) error {
		query, args, err := s.sqlSetPassword(user)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, query, args...)
		return err
	})

	return
}

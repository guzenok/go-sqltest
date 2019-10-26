package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guzenok/go-sqltest/sample/store"
)

func TestConnection_CreateUser(t *testing.T) {
	s := testConnection(t)
	ctx := context.Background()

	user := &store.User{
		Login:    "user01",
		Password: "123456",
	}
	newPassword := "654321"

	t.Run("create", func(t *testing.T) {
		var err error
		user, err = s.CreateUser(ctx, user)
		assert.NoError(t, err)
	})

	t.Run("exists", func(t *testing.T) {
		var err error
		user, err = s.CreateUser(ctx, user)
		assert.Error(t, err)
	})

	t.Run("set password", func(t *testing.T) {
		err := s.SetPassword(ctx, user.ID, newPassword)
		assert.NoError(t, err)
	})

	t.Run("get by login", func(t *testing.T) {
		var err error
		user, err = s.GetUserByLogin(ctx, user.Login)
		assert.NoError(t, err)
		assert.Equal(t, newPassword, user.Password)
	})

	t.Run("delete", func(t *testing.T) {
		err := s.DeleteUser(ctx, user.ID)
		assert.NoError(t, err)
	})

	t.Run("get by id", func(t *testing.T) {
		var err error
		user, err = s.GetUserByID(ctx, user.ID)
		assert.Error(t, err)
	})
}

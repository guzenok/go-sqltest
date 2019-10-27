package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/guzenok/go-sqltest/sample/store"
)

func TestConnection_CreateUser(t *testing.T) {
	ctx := context.Background()

	db, _, err := UsersTestDb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	s := Wrap(db)

	user := &store.User{
		ID:       1,
		Login:    "user01",
		Password: "123456",
	}
	newPassword := "654321"

	_ = true &&

		t.Run("exists", func(t *testing.T) {
			var err error
			_, err = s.CreateUser(ctx, user)
			assert.Error(t, err)
		}) &&

		t.Run("delete", func(t *testing.T) {
			err := s.DeleteUser(ctx, user.ID)
			assert.NoError(t, err)
		}) &&

		t.Run("create", func(t *testing.T) {
			var err error
			user, err = s.CreateUser(ctx, user)
			assert.NoError(t, err)
		}) &&

		t.Run("set password", func(t *testing.T) {
			err := s.SetPassword(ctx, user.ID, newPassword)
			assert.NoError(t, err)
		}) &&

		t.Run("get by login", func(t *testing.T) {
			var err error
			user, err = s.GetUserByLogin(ctx, user.Login)
			_ = assert.NoError(t, err) &&
				assert.NotNil(t, user) &&
				assert.Equal(t, newPassword, user.Password)
		}) &&

		t.Run("get by id", func(t *testing.T) {
			var err error
			user, err = s.GetUserByID(ctx, user.ID)
			_ = assert.NoError(t, err) &&
				assert.NotNil(t, user) &&
				assert.Equal(t, newPassword, user.Password)
		})
}

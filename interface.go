package db

import (
	"context"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id    int
	Name  string
	Email string
}

var ErrNotFound = errors.New("not found")

type Database interface {
	UserById(ctx context.Context, Id int) (*User, error)
	CreateUser(ctx context.Context, u *User) (*User, error)
	UpdateUser(ctx context.Context, u *User) error
	DeleteUser(ctx context.Context, Id int) error
	ListAllUsers(ctx context.Context) ([]*User, error)
}

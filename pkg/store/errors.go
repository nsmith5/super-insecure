package store

import "errors"

var (
	ErrorItemNotFound      = errors.New("item not found")
	ErrorUserNotFound      = errors.New("user not found")
	ErrorUserAlreadyExists = errors.New("user already exists")
)

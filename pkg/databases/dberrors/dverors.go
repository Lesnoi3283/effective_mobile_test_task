package dberrors

import "errors"

var errNotFound error = errors.New("no content found")

func NewNotFoundErr() error {
	return errNotFound
}

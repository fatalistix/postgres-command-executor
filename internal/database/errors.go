package database

import "errors"

// command

var ErrCommandExists = errors.New("command exists")
var ErrCommandNotFound = errors.New("command not found")

// process

var ErrProcessNotFound = errors.New("process not found")

// unknown

var ErrInternal = errors.New("internal error")

package common

import (
	"errors"
)

var UncreatedUserError = errors.New("user does not exist")
var UserPasswordError = errors.New("wrong password")

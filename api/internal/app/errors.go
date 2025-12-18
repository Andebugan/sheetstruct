package app

import "errors"

var ErrUserAuthFailed = errors.New("Unable to find user with matching credentials")
var ErrUserNotFound = errors.New("Unable to find requested user")
var ErrUserNameAlredyExists = errors.New("User with provided name already exists")
var ErrUserEmailAlredyExists = errors.New("User with provided email already exists")

var ErrSheetNotFound = errors.New("Unable to find requested user")
var ErrSheetAlredyExists = errors.New("User already exists")

var ErrComponentNotFound = errors.New("Unable to find requested user")
var ErrComponentAlredyExists = errors.New("User already exists")


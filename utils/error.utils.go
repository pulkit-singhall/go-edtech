package utils

import "errors"

var QueryParamMissing = errors.New("Query parameter missing")
var QueryBodyMissing = errors.New("Query body missing")
var ValidationFailed = errors.New("Validation failed")
var UserAlreadyExists = errors.New("User already exists with same email")
var InternalServerError = errors.New("Internal server error")
var HexIdError = errors.New("Conversion into HEX Error")
var UserNotFound = errors.New("User not found")
var PasswordWrong = errors.New("Wrong password provided")
var TokenError = errors.New("Token generation error")
var AuthorizeError = errors.New("User not authorize")

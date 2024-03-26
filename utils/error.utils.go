package utils

import "errors"

var QueryParamMissing = errors.New("Query parameter missing")
var QueryBodyMissing = errors.New("Query body missing")
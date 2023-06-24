package api

import "errors"

var SessionClosedError = errors.New("ISession Closed")
var SessionBlockedError = errors.New("ISession Blocked")

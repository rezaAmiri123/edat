package edathttp

import (

_	"github.com/go-chi/render"
_	"github.com/stackus/errors"
)
type ErrorResponse  struct{
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Err        error  `json:"-"`
}

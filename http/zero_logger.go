package edathttp

import (
	"net/http"
	_"github.com/go-chi/chi/v5/middleware"
	_"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

func ZeroLogger(logger zerolog.Logger)func(next http.Handler)http.Handler{
	return func(next http.Handler) http.Handler {
		
	}
}


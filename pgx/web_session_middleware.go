package edatpgx

import (
	"net/http"
	_"github.com/jackc/pgx/v4/pgxpool"
)

func WebSessionMiddleware(conn *pgxp)


type statusWriter struct{
	http.ResponseWriter
	status int
}

func newStatusWriter(w http.ResponseWriter)*statusWriter{
	return &statusWriter{
		ResponseWriter: w,
	}
}

func(s *statusWriter)WriteHeader(status int){
	s.status = status
	s.ResponseWriter.WriteHeader(status)
}

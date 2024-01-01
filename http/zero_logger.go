package edathttp

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"github.com/rezaAmiri123/edat/core"
	"github.com/rs/zerolog"
)

func ZeroLogger(logger zerolog.Logger)func(next http.Handler)http.Handler{
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ww:= middleware.NewWrapResponseWriter(writer, request.ProtoMajor)

			start := time.Now()

			requestID := core.GetRequestID(request.Context())
			correlationID := core.GetCorrelationID(request.Context())
			causationID := core.GetCausationID(request.Context())

			defer func(){
				var err error
				var logFn func()*zerolog.Event

				p:= recover()

				switch{
				case p!= nil:
					logFn = logger.Error().Stack
					// ensure the status code reflects this panic
					if ww.Status() <500{
						ww.WriteHeader(http.StatusInternalServerError)
					}
					err = errors.Errorf("%s", p)
				case ww.Status()<400:
					logFn = logger.Info
				case ww.Status()<500:
					logFn = logger.Warn
				default:
					logFn = logger.Error
				}
				log := logFn()

				if err!=nil{
					log = log.Err(err)
				}
				log = log.Str("RemoteAddr", request.RemoteAddr).
					Int("ContextLength", ww.BytesWritten()).
					Dur("ResponseTime", time.Since(start))

				if requestID!= ""{
					log.Str("RequestID", requestID).
					Str("CorrelationID", correlationID).
					Str("CausationID", causationID)
				}

				log.Msgf("[%d] %s %s", ww.Status(), request.Method, request.RequestURI)
			}()

			next.ServeHTTP(ww, request)
		})		
	}
}


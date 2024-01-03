package edatpgx

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rezaAmiri123/edat/log"
	"google.golang.org/grpc"
)

func RpcSessionUnrayInterceptor(conn *pgxpool.Pool, logger edatlog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		tx, err := conn.Begin(ctx)
		if err != nil {
			logger.Error("error while starting the request transaction", edatlog.Error(err))
			return
		}

		newCtx := context.WithValue(ctx, pgxTxKey, tx)

		defer func() {
			p := recover()
			switch {
			case p != nil:
				txErr := tx.Rollback(ctx)
				if txErr != nil {
					logger.Error("error while rolling back the rpc request transaction during panic", edatlog.Error(txErr))
				}
				panic(p)
			case err != nil:
				txErr := tx.Rollback(ctx)
				if txErr != nil {
					logger.Error("error while rollng back the rpc request transaction", edatlog.Error(txErr))
				}
			default:
				txErr := tx.Commit(ctx)
				if txErr != nil {
					logger.Error("error while commiting the rpc request transaction", edatlog.Error(err))
				}
			}
		}()

		return handler(newCtx, req)
	}
}

func RpcSessionStreamInterceptor(_ *pgxpool.Pool, logger edatlog.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		logger.Error("outbox pattern not yet implemented for streaming connections")
		return handler(srv, ss)
	}
}

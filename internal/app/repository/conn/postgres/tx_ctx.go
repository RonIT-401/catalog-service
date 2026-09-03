package rcpostgres

import (
	"context"

	"github.com/uptrace/bun"
)

type contextKeyTx struct{}

func getTxFromContext(ctx context.Context) bun.Tx {
	value, ok := ctx.Value(contextKeyTx{}).(bun.Tx)
	if !ok {
		return bun.Tx{}
	}

	return value
}

func setTxToContext(ctx context.Context, tx bun.Tx) context.Context {
	return context.WithValue(ctx, contextKeyTx{}, tx)
}

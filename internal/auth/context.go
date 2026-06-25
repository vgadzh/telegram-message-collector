package auth

import (
	"context"
)

type principalKey struct{}

var ctxPrincipalKey = principalKey{}

type Principal struct {
	UserID int64
	Role   string
}

func WithPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, ctxPrincipalKey, p)
}

func PrincipalFromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(ctxPrincipalKey).(Principal)
	return p, ok
}

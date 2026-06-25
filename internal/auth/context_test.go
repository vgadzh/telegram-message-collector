package auth

import (
	"context"
	"testing"
)

func TestWithPrincipal(t *testing.T) {
	ctx := context.Background()

	expected := Principal{
		UserID: 123,
		Role:   "ADMIN",
	}

	ctx = WithPrincipal(ctx, expected)
	actual, ok := PrincipalFromContext(ctx)
	if !ok {
		t.Fatal("principal not found")
	}
	if actual != expected {
		t.Fatalf("got %+v want %+v", actual, expected)
	}
}

func TestPrincipalFromContext_NotFound(t *testing.T) {
	ctx := context.Background()
	p, ok := PrincipalFromContext(ctx)
	if ok {
		t.Fatal("expected false")
	}
	if p != (Principal{}) {
		t.Fatal("expected zero principal")
	}
}

func TestPrincipalIsolation(t *testing.T) {
	type otherKey struct{}

	ctx := context.WithValue(context.Background(), otherKey{}, "hello")
	_, ok := PrincipalFromContext(ctx)
	if ok {
		t.Fatal("unexpected principal")
	}
}

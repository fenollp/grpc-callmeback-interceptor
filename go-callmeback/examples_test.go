package callmeback_test

import (
	"context"
	"time"

	"github.com/fenollp/grpc-callmeback-interceptor/go-callmeback"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

type everySecond struct{}

func (*everySecond) PleaseComeAgain(ctx context.Context, c callmeback.Context) (time.Duration, error) {
	return time.Second, nil
}

// Simple example of server initialization code.
func Example() {
	// Create unary/stream rateLimiters, based on token bucket here.
	cb := &everySecond{}
	_ = grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			callmeback.UnaryServerInterceptor(cb),
		),
	)
}

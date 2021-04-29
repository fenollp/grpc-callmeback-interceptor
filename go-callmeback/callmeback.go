package callmeback

import (
	"context"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Context can be used to make more accurate predictions of the duration.
type Context struct {
	Elapsed    time.Duration
	FullMethod string
}

// CallMeBacker defines the interface to instruct client when to call back.
type CallMeBacker interface {
	// If PleaseComeAgain returns a non-nil error, the request will have been processed for nothing.
	// If it returns a positive duration, a trailer is added to the response.
	PleaseComeAgain(ctx context.Context, c Context) (time.Duration, error)
}

// UnaryServerInterceptor returns a new unary server interceptor that instructs the client to call again.
func UnaryServerInterceptor(callmebacker CallMeBacker) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		rep, err := handler(ctx, req)
		end := time.Since(start)
		if err != nil {
			return nil, err
		}

		c := Context{
			Elapsed:    end,
			FullMethod: info.FullMethod,
		}
		var d time.Duration
		if d, err = callmebacker.PleaseComeAgain(ctx, c); err != nil {
			return nil, err
		}

		var dstr string
		switch {
		case strings.HasSuffix(Trailer, "-ms"):
			dstr = strconv.FormatInt(int64(d/time.Millisecond), 10)
		default:
			dstr = d.String()
		}

		md := metadata.Pairs(Trailer, dstr)
		if err = grpc.SetTrailer(ctx, md); err != nil {
			return nil, err
		}
		return rep, nil
	}
}

// ErrBadTrailer is a constant error
const ErrBadTrailer = errStr("bad callmeback trailer")

type errStr string

func (es errStr) Error() string { return string(es) }

// In reads and parses the duration trailer or fails with an error.
func In(md metadata.MD) (d time.Duration, err error) {
	vs := md.Get(Trailer)
	if len(vs) != 1 {
		err = ErrBadTrailer
		return
	}
	d, err = time.ParseDuration(vs[0])
	return
}

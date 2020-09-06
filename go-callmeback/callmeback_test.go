package callmeback_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fenollp/grpc-callmeback-interceptor/go-callmeback"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
	errMsgFake = "fake error"
	pause      = 1 * time.Second
)

var wentThroughHere = false

type mockCMBer struct{}

var _ callmeback.CallMeBacker = &mockCMBer{}

func (*mockCMBer) PleaseComeAgain(ctx context.Context, c callmeback.Context) (time.Duration, error) {
	wentThroughHere = true
	return pause, nil
}

func TestUnaryServerInterceptor_CallMeBackFail(t *testing.T) {
	require.False(t, wentThroughHere)

	interceptor := callmeback.UnaryServerInterceptor(&mockCMBer{})
	handler := func(ctx context.Context, req interface{}) (rep interface{}, err error) {
		require.False(t, wentThroughHere)
		err = errors.New(errMsgFake)
		return
	}

	info := &grpc.UnaryServerInfo{FullMethod: "/pkg.Svc/FakeMethod/"}
	ctx := context.Background()
	rep, err := interceptor(ctx, nil, info, handler)
	require.EqualError(t, err, errMsgFake)
	require.Nil(t, rep)

	require.False(t, wentThroughHere)
}

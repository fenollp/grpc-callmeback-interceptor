package callmeback

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

const (
	errMsgFake = "fake error"
	pause      = 1 * time.Second
)

type mockCMBer struct{}

func (*mockCMBer) PleaseComeAgain(ctx context.Context, c Context) (time.Duration, error) {
	return pause, nil
}

func TestUnaryServerInterceptor_CallMeBackFail(t *testing.T) {
	interceptor := UnaryServerInterceptor(&mockCMBer{})
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New(errMsgFake)
	}
	info := &grpc.UnaryServerInfo{FullMethod: "/pkg.Svc/FakeMethod/"}
	ctx := context.Background()
	rep, err := interceptor(ctx, nil, info, handler)
	require.EqualError(t, err, errMsgFake)
	require.Nil(t, rep)
}

// func TestUnaryServerInterceptor_CallMeBackPass(t *testing.T) {
// 	interceptor := UnaryServerInterceptor(&mockCMBer{})
// 	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
// 		return &struct{}{}, nil
// 	}
// 	info := &grpc.UnaryServerInfo{FullMethod: "/pkg.Svc/FakeMethod/"}
// 	ctx := context.Background()
// 	rep, err := interceptor(ctx, nil, info, handler)
// 	require.NoError(t, err)
// 	require.NotNil(t, rep)
// 	// md, ok := metadata.FromIncomingContext(ctx)
// 	md, ok := metadata.FromOutgoingContext(ctx)
// 	require.True(t, ok)
// 	pauseFor, err := In(md)
// 	require.NoError(t, err)
// 	require.Equal(t, pauseFor, pause)
// }

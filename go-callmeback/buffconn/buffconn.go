package buffconn

import (
	"context"
	"net"
	"time"

	"github.com/fenollp/grpc-callmeback-interceptor/go-callmeback"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const (
	pause = 1 * time.Second

	bufSize = 1024 * 1024
)

type mockCMBer struct{}

func (*mockCMBer) PleaseComeAgain(ctx context.Context, c callmeback.Context) (time.Duration, error) {
	return pause, nil
}

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			callmeback.UnaryServerInterceptor(&mockCMBer{}),
		),
	)
	RegisterEchoServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type server struct{}

var _ EchoServer = &server{}

func (s *server) UnaryEcho(ctx context.Context, req *EchoRequest) (rep *EchoResponse, err error) {
	// defer func() {
	//     trailer := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	//     grpc.SetTrailer(ctx, trailer)
	// }()
	// header := metadata.New(map[string]string{"location": "MTV", "timestamp": time.Now().Format(timestampFormat)})
	// grpc.SendHeader(ctx, header)
	rep = &EchoResponse{Message: req.Message}
	return
}

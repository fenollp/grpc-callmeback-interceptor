package buffconn

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/stretchr/testify/require"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/test/bufconn"
)

const (
    pause = 1 * time.Second

    bufSize = 1024 * 1024
)

type mockCMBer struct{}

func (*mockCMBer) PleaseComeAgain(ctx context.Context, c Context) (time.Duration, error) {
    return pause, nil
}

var lis *bufconn.Listener

func init() {
    lis = bufconn.Listen(bufSize)
    s := grpc.NewServer()
    pb.RegisterEchoServer(s, &server{})
    go func() {
        if err := s.Serve(lis); err != nil {
            log.Fatalf("Server exited with error: %v", err)
        }
    }()
}

func bufDialer(context.Context, string) (net.Conn, error) { return lis.Dial() }

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

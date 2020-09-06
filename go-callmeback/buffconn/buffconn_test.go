package buffconn

import (
    "context"
    "testing"
    "time"

    "github.com/fenollp/grpc-callmeback-interceptor/go-callmeback"
    "github.com/stretchr/testify/require"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
)

func TestUnaryServerInterceptor_CallMeBackPass(t *testing.T) {
    ctx := context.Background()
    conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
    require.NoError(t, err)
    defer conn.Close()

    client := NewEchoClient(conn)
    const msg = "Hello Joe"

    var trailer metadata.MD
    req := &EchoRequest{Message: msg}
    rep, err := client.UnaryEcho(ctx, req, grpc.Trailer(&trailer))
    require.NoError(t, err)
    require.Equal(t, rep.Message, msg)

    pauseFor, err := callmeback.In(trailer)
    require.NoError(t, err)
    require.Equal(t, pauseFor, 1*time.Second)
}

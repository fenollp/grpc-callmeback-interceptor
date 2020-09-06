package buffconn_test

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/fenollp/grpc-callmeback-interceptor/go-callmeback/bufconn"
    "github.com/stretchr/testify/require"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
)

func TestUnaryServerInterceptor_CallMeBackPass(t *testing.T) {
    ctx := context.Background()
    conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
    require.NoError(t, err)
    defer conn.Close()
    client := buffconn.NewEchoClient(conn)
    var trailer metadata.MD
    rep, err := client.UnaryEcho(ctx, &EchoRequest{"Dr. Seuss"}, grpc.Trailer(&trailer))
    require.NoError(t, err)
}

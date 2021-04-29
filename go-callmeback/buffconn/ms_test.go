package buffconn_test

import (
	"context"
	"testing"

	"github.com/fenollp/grpc-callmeback-interceptor/go-callmeback"
	"github.com/fenollp/grpc-callmeback-interceptor/go-callmeback/buffconn"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestUnaryServerInterceptor_CallMeBackMSPass(t *testing.T) {
	callmeback.Trailer = "x-pleasecomeagain-ms"

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(buffconn.NewDialer), grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	client := buffconn.NewEchoClient(conn)
	const msg = "Hello Joe"

	var trailer metadata.MD
	req := &buffconn.EchoRequest{Message: msg}
	rep, err := client.UnaryEcho(ctx, req, grpc.Trailer(&trailer))
	require.NoError(t, err)
	require.Equal(t, rep.Message, msg)

	require.Equal(t, []string{"1000"}, trailer.Get(callmeback.Trailer))
}

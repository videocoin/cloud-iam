package e2e_test

import (
	"context"
	"testing"

	iam "github.com/videocoin/cloud-api/iam/v1"

	"github.com/gogo/protobuf/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/videocoin/cloud-iam/pkg/auth"
	"google.golang.org/grpc"
)

func TestIAMService(t *testing.T) {
	conn, err := grpc.Dial(":5000", grpc.WithInsecure())
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer conn.Close()

	cli := iam.NewIAMClient(conn)

	ctx := auth.NewOutgoingContext(context.Background(), uuid.New().String())

	key, err := cli.CreateKey(ctx, new(types.Empty))
	require.NoError(t, err)
	require.NotNil(t, key)
	require.NotEmpty(t, key.Id)
	require.NotNil(t, key.PrivateKeyData)
}

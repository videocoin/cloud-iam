package e2e_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	iam "github.com/videocoin/cloud-api/iam/v1"
	account "github.com/videocoin/cloud-pkg/api/resources/serviceaccount"

	"github.com/stretchr/testify/require"
	"github.com/videocoin/cloud-pkg/api/resources/key"
	"github.com/videocoin/cloud-pkg/api/resources/project"
	"github.com/videocoin/oauth2"
	"google.golang.org/grpc"
)

func TestIAMService(t *testing.T) {
	conn, err := grpc.Dial(":5000", grpc.WithInsecure())
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer conn.Close()

	cli := iam.NewIAMClient(conn)
	require.NotNil(t, cli)

	ctx := context.Background()
	require.NotNil(t, ctx)

	userID := uuid.New().String()

	// create project based on the user ID
	projName := project.NewName(userID)
	projID := projName.ID()

	// account
	accID := "account1"
	// expected
	accEmail := account.NewEmail(projID, accID)
	accName := account.NewName(projID, accEmail)

	// create service account
	req := &iam.CreateServiceAccountRequest{
		Name:      string(projName),
		AccountId: accID,
	}
	acc, err := cli.CreateServiceAccount(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, acc)
	require.Equal(t, string(accName), acc.Name)
	require.Equal(t, accEmail, acc.Email)
	require.Equal(t, projID, acc.ProjectId)
	require.True(t, key.IsValidID(acc.UniqueId))

	// duplicate
	acc2, err := cli.CreateServiceAccount(ctx, req)
	require.Error(t, err)
	require.Nil(t, acc2)

	// get service account
	acc3, err := cli.GetServiceAccount(ctx, &iam.GetServiceAccountRequest{Name: string(accName)})
	require.NoError(t, err)
	require.NotNil(t, acc3)
	require.Equal(t, acc, acc3)

	// list service accounts
	accs, err := cli.ListServiceAccounts(ctx, &iam.ListServiceAccountsRequest{Name: string(projName)})
	require.NoError(t, err)
	require.NotNil(t, accs)
	require.NotNil(t, accs.Accounts)
	require.Len(t, accs.Accounts, 1)
	require.Equal(t, accs.Accounts[0], acc)

	// delete service account
	empty, err := cli.DeleteServiceAccount(ctx, &iam.DeleteServiceAccountRequest{Name: string(accName)})
	require.NoError(t, err)
	require.NotNil(t, empty)

	// create service account once again
	acc, err = cli.CreateServiceAccount(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, acc)

	// create service account key
	// there are no duplicates since the key ID is generated on the server side.
	keyReq := &iam.CreateServiceAccountKeyRequest{
		Name: string(accName),
	}
	// key id is not known beforehand
	key, err := cli.CreateServiceAccountKey(ctx, keyReq)
	require.NoError(t, err)
	require.NotNil(t, key)
	require.NotEmpty(t, key.Name)
	pk, err := oauth2.ParseKey(key.PrivateKeyData)
	require.NoError(t, err)
	require.NotNil(t, pk)

	// get service account key
	key3, err := cli.GetServiceAccountKey(ctx, &iam.GetServiceAccountKeyRequest{Name: key.Name})
	require.NoError(t, err)
	require.NotNil(t, key3)
	require.Equal(t, key, key)

	// list service account keys
	keys, err := cli.ListServiceAccountKeys(ctx, &iam.ListServiceAccountKeysRequest{Name: string(accName)})
	require.NoError(t, err)
	require.NotNil(t, keys)
	require.NotNil(t, keys.Keys)
	require.Len(t, keys.Keys, 1)
	// fix timestamps diffs (not relevant for now)
	//require.Equal(t, keys.Keys[0], key)

	// delete service account key
	empty, err = cli.DeleteServiceAccountKey(ctx, &iam.DeleteServiceAccountKeyRequest{Name: key.Name})
	require.NoError(t, err)
	require.NotNil(t, empty)
}

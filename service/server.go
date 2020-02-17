package service

import (
	"context"

	guuid "github.com/google/uuid"
	iam "github.com/videocoin/cloud-api/iam/v1"
	keyspec "github.com/videocoin/cloud-pkg/api/resources/key"
	accspec "github.com/videocoin/cloud-pkg/api/resources/serviceaccount"

	"github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/cloud-iam/datastore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the IAMServer interface.
type Server struct {
	logger     *logrus.Entry
	ds         datastore.DataStore
	passphrase string
	roleDs     datastore.RoleDatastore
}

// NewServer creates an IAM server.
func NewServer(logger *logrus.Entry, ds datastore.DataStore, passphrase string) *Server {
	return &Server{
		logger:     logger,
		ds:         ds,
		passphrase: passphrase,
	}
}

// CreateServiceAccount creates a service account.
func (srv *Server) CreateServiceAccount(ctx context.Context, req *iam.CreateServiceAccountRequest) (*iam.ServiceAccount, error) {
	principal := ctx.Value("principal").(string)

	if ok := accspec.IsValidID(req.AccountId); !ok {
		return nil, status.Error(codes.InvalidArgument, accspec.ErrInvalidID.Error())
	}

	acc, err := srv.ds.CreateServiceAccount(&datastore.ServiceAccount{
		ID:     guuid.New().String(),
		UserID: principal,
		Email:  accspec.NewEmail(principal, req.AccountId),
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return acc.Proto(), nil
}

// GetServiceAccount gets a service account.
func (srv *Server) GetServiceAccount(ctx context.Context, req *iam.GetServiceAccountRequest) (*iam.ServiceAccount, error) {
	principal := ctx.Value("principal").(string)

	ok := accspec.IsValidEmail(req.Email)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, accspec.ErrInvalidEmail.Error())
	}

	acc, err := srv.ds.GetServiceAccountByEmail(principal, req.Email)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return acc.Proto(), nil
}

// ListServiceAccounts lists service accounts.
func (srv *Server) ListServiceAccounts(ctx context.Context, req *iam.ListServiceAccountsRequest) (*iam.ListServiceAccountsResponse, error) {
	principal := ctx.Value("principal").(string)

	accs, err := srv.ds.ListServiceAccounts(principal)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	accsPB := make([]*iam.ServiceAccount, 0, len(accs))
	for _, acc := range accs {
		accsPB = append(accsPB, acc.Proto())
	}

	return &iam.ListServiceAccountsResponse{Accounts: accsPB}, nil
}

// DeleteServiceAccount deletes a service account.
func (srv *Server) DeleteServiceAccount(ctx context.Context, req *iam.DeleteServiceAccountRequest) (*types.Empty, error) {
	ok := accspec.IsValidEmail(req.Email)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, accspec.ErrInvalidEmail.Error())
	}
	if err := srv.ds.DeleteServiceAccount(req.Email); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return new(types.Empty), nil
}

// CreateServiceAccountKey creates a service account key.
func (srv *Server) CreateServiceAccountKey(ctx context.Context, req *iam.CreateServiceAccountKeyRequest) (*iam.ServiceAccountKey, error) {
	ok := accspec.IsValidEmail(req.ServiceAccountEmail)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, accspec.ErrInvalidEmail.Error())
	}

	key, err := srv.ds.CreateServiceAccountKey(req.ServiceAccountEmail, srv.passphrase)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	keyPB, err := key.CreationProto(srv.passphrase)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return keyPB, nil
}

// GetServiceAccountKey gets a service account key.
func (srv *Server) GetServiceAccountKey(ctx context.Context, req *iam.GetServiceAccountKeyRequest) (*iam.ServiceAccountKey, error) {
	ok := keyspec.IsValidID(req.KeyId)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, keyspec.ErrInvalidID.Error())
	}

	key, err := srv.ds.GetServiceAccountKey(req.KeyId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	keyPB, err := key.Proto()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return keyPB, nil
}

// ListServiceAccountKeys lists service account keys.
func (srv *Server) ListServiceAccountKeys(ctx context.Context, req *iam.ListServiceAccountKeysRequest) (*iam.ListServiceAccountKeysResponse, error) {
	if ok := accspec.IsValidEmail(req.ServiceAccountEmail); !ok {
		return nil, status.Error(codes.InvalidArgument, accspec.ErrInvalidEmail.Error())
	}

	keys, err := srv.ds.ListServiceAccountKeysByEmail(req.ServiceAccountEmail)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	keysPB := make([]*iam.ServiceAccountKey, 0, len(keys))
	for _, key := range keys {
		keyPB, err := key.Proto()
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		keysPB = append(keysPB, keyPB)
	}

	return &iam.ListServiceAccountKeysResponse{Keys: keysPB}, nil
}

// DeleteServiceAccountKey deletes a service account key.
func (srv *Server) DeleteServiceAccountKey(ctx context.Context, req *iam.DeleteServiceAccountKeyRequest) (*types.Empty, error) {
	ok := keyspec.IsValidID(req.KeyId)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, keyspec.ErrInvalidID.Error())
	}

	if err := srv.ds.DeleteServiceAccountKey(req.KeyId); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return new(types.Empty), nil
}

// GetIamPolicy gets an IAM policy.
func (srv *Server) GetIamPolicy(ctx context.Context, req *iam.GetIamPolicyRequest) (*iam.Policy, error) {
	return nil, nil
}

// SetIamPolicy sets an IAM policy.
func (srv *Server) SetIamPolicy(ctx context.Context, req *iam.SetIamPolicyRequest) (*iam.Policy, error) {
	return nil, nil
}

// GetRole gets a predefined role.
func (srv *Server) GetRole(ctx context.Context, req *iam.GetRoleRequest) (*iam.Role, error) {
	// TODO validate role name

	role, err := srv.roleDs.GetRole(req.Name)
	if err != nil {
		return nil, err
	}

	return role.Proto(), nil
}

// ListRoles lists the predefined roles.
func (srv *Server) ListRoles(ctx context.Context, req *iam.ListRolesRequest) (*iam.ListRolesResponse, error) {
	roles, err := srv.roleDs.ListRoles()
	if err != nil {
		return nil, err
	}

	rolesPB := make([]*iam.Role, 0, len(roles))
	for _, role := range roles {
		rolesPB = append(rolesPB, role.Proto())
	}

	return &iam.ListRolesResponse{Roles: rolesPB}, nil
}

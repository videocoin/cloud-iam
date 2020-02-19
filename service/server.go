package service

import (
	"context"

	iam "github.com/videocoin/cloud-api/iam/v1"
	keyspec "github.com/videocoin/cloud-pkg/api/resources/key"

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

// CreateKey creates an user key.
func (srv *Server) CreateKey(ctx context.Context, empty *types.Empty) (*iam.Key, error) {
	value := ctx.Value("subject")
	principal, ok := ctx.Value("subject").(string)
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "principal required")
	}
	if !keyspec.IsValidID(principal) {
		return nil, status.Error(codes.FailedPrecondition, "invalid principal")
	}

	return srv.createKey(principal)
}

func (srv *Server) createKey(userID string) return (*iam.Key, error) {
	key, err = generateKey(rand.Reader, srv.passphrase, userID)
	if err != nil {
		// TODO
		return nil, err
	}

	key, err := srv.ds.CreateKey(userID, srv.passphrase)
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

// CreateRoleBinding creates a role binding.
func (srv *Server) CreateRoleBinding(ctx context.Context, req *iam.RoleBinding) (*types.Empty, error) {
	return nil, nil
}

// DeleteRoleBinding deletes a role binding.
func (srv *Server) DeleteRoleBinding(ctx context.Context, req *iam.RoleBinding) (*types.Empty, error) {
	return nil, nil
}

// ListRoleBindings lists role bindings.
func (srv *Server) ListRoleBindings(context.Context, *iam.ListRoleBindingsRequest) (*iam.ListRoleBindingsResponse, error) {
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

package service

import (
	"context"
	"crypto/rand"

	iam "github.com/videocoin/cloud-api/iam/v1"

	"github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/cloud-iam/datastore"
	"github.com/videocoin/cloud-iam/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the IAMServer interface.
type Server struct {
	logger     *logrus.Entry
	ds         datastore.DataStore
	passphrase string
}

// NewServer creates an IAM server.
func NewServer(logger *logrus.Entry, ds datastore.DataStore, passphrase string) *Server {
	return &Server{
		logger:     logger,
		ds:         ds,
		passphrase: passphrase,
	}
}

// CreateKey creates a key for an authenticated user.
func (srv *Server) CreateKey(ctx context.Context, empty *types.Empty) (*iam.Key, error) {
	subject, err := auth.FromIncomingContext(ctx)
	if err != nil {
		srv.logger.Error(err)
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	key, err := srv.createUserKey(subject)
	if err != nil {
		srv.logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return key, nil
}

func (srv *Server) createUserKey(userID string) (*iam.Key, error) {
	key, err := generateKey(rand.Reader, srv.passphrase, userID)
	if err != nil {
		return nil, err
	}

	if err := srv.ds.CreateUserKey(key); err != nil {
		return nil, err
	}

	keyPB, err := key.CreationProto(srv.passphrase)
	if err != nil {
		return nil, err
	}

	return keyPB, nil
}

// GetKey gets a key for an authenticated user.
func (srv *Server) GetKey(ctx context.Context, req *iam.GetKeyRequest) (*iam.Key, error) {
	subject, err := auth.FromIncomingContext(ctx)
	if err != nil {
		srv.logger.Error(err)
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	key, err := srv.getUserKey(subject, req.KeyId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			srv.logger.Debug(err)
			return nil, status.Error(codes.NotFound, err.Error())
		}
		srv.logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return key, nil
}

func (srv *Server) getUserKey(userID string, keyID string) (*iam.Key, error) {
	key, err := srv.ds.GetUserKey(userID, keyID)
	if err != nil {
		return nil, err
	}

	keyPB, err := key.Proto()
	if err != nil {
		return nil, err
	}

	return keyPB, nil
}

// ListKeys lists keys for an authenticated user.
func (srv *Server) ListKeys(ctx context.Context, req *iam.ListKeysRequest) (*iam.ListKeysResponse, error) {
	subject, err := auth.FromIncomingContext(ctx)
	if err != nil {
		srv.logger.Error(err)
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	keys, err := srv.listUserKeys(subject)
	if err != nil {
		srv.logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &iam.ListKeysResponse{Keys: keys}, nil
}

func (srv *Server) listUserKeys(userID string) ([]*iam.Key, error) {
	keys, err := srv.ds.ListUserKeys(userID)
	if err != nil {
		return nil, err
	}

	keysPB := make([]*iam.Key, 0, len(keys))
	for _, key := range keys {
		keyPB, err := key.Proto()
		if err != nil {
			return nil, err
		}
		keysPB = append(keysPB, keyPB)
	}

	return keysPB, nil
}

// DeleteKey deletes an user key.
func (srv *Server) DeleteKey(ctx context.Context, req *iam.DeleteKeyRequest) (*types.Empty, error) {
	subject, err := auth.FromIncomingContext(ctx)
	if err != nil {
		srv.logger.Error(err)
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	if err := srv.deleteUserKey(subject, req.KeyId); err != nil {
		srv.logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return new(types.Empty), nil
}

func (srv *Server) deleteUserKey(userID string, keyID string) error {
	return srv.ds.DeleteUserKey(userID, keyID)
}

// ListRoleBindings lists role bindings.
func (srv *Server) ListRoleBindings(context.Context, *iam.ListRoleBindingsRequest) (*iam.ListRoleBindingsResponse, error) {
	// TODO
	return nil, nil
}

// CreateRoleBinding binds a role to an user.
func (srv *Server) CreateRoleBinding(ctx context.Context, req *iam.RoleBinding) (*types.Empty, error) {
	// TODO
	return nil, nil
}

// DeleteRoleBinding deletes a role binding.
func (srv *Server) DeleteRoleBinding(ctx context.Context, req *iam.RoleBinding) (*types.Empty, error) {
	// TODO
	return nil, nil
}

// GetRole gets a predefined role.
func (srv *Server) GetRole(ctx context.Context, req *iam.GetRoleRequest) (*iam.Role, error) {
	// TODO
	return nil, nil
}

// ListRoles lists the predefined roles.
func (srv *Server) ListRoles(ctx context.Context, req *iam.ListRolesRequest) (*iam.ListRolesResponse, error) {
	// TODO
	return nil, nil
}

// ListUserRoles lists the user roles.
func (srv *Server) ListUserRoles(ctx context.Context, req *iam.ListRolesRequest) (*iam.ListRolesResponse, error) {
	// TODO
	return nil, nil
}

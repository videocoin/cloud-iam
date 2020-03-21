package service

import (
	"context"
	"crypto/rand"

	iam "github.com/videocoin/videocoinapis/videocoin/iam/v1"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"github.com/videocoin/cloud-iam/datastore"
	"github.com/videocoin/common/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server implements the IAMServer interface.
type Server struct {
	ds         datastore.DataStore
	passphrase string
}

// NewServer creates an IAM server.
func NewServer(ds datastore.DataStore, passphrase string) *Server {
	return &Server{
		ds:         ds,
		passphrase: passphrase,
	}
}

// CreateKey creates a key for an authenticated user.
func (srv *Server) CreateKey(ctx context.Context, empty *empty.Empty) (*iam.Key, error) {
	sub, err := subjectFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "invalid auth info")
	}

	key, err := srv.createUserKey(sub)
	if err != nil {
		log.Errorln(err)
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
	sub, err := subjectFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "invalid auth info")
	}

	key, err := srv.getUserKey(sub, req.KeyId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Debugln(err)
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Errorln(err)
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
	sub, err := subjectFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "invalid auth info")
	}

	keys, err := srv.listUserKeys(sub)
	if err != nil {
		log.Errorln(err)
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
func (srv *Server) DeleteKey(ctx context.Context, req *iam.DeleteKeyRequest) (*empty.Empty, error) {
	sub, err := subjectFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "invalid auth info")
	}

	if err := srv.deleteUserKey(sub, req.KeyId); err != nil {
		log.Errorln(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return new(empty.Empty), nil
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
func (srv *Server) CreateRoleBinding(ctx context.Context, req *iam.RoleBinding) (*empty.Empty, error) {
	// TODO
	return nil, nil
}

// DeleteRoleBinding deletes a role binding.
func (srv *Server) DeleteRoleBinding(ctx context.Context, req *iam.RoleBinding) (*empty.Empty, error) {
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

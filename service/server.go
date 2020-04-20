package service

import (
	"context"
	"crypto/rand"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jinzhu/gorm"
	"github.com/videocoin/cloud-iam/datastore"
	"github.com/videocoin/common/log"
	"github.com/videocoin/runtime/grpc/middleware/auth"
	iam "github.com/videocoin/videocoinapis/videocoin/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{ ds datastore.DataStore }

// New creates an IAM server.
func New(ds datastore.DataStore) *server {
	return &server{
		ds: ds,
	}
}

// CreateKey creates a key for an authenticated user.
func (s *server) CreateKey(ctx context.Context, empty *empty.Empty) (*iam.Key, error) {
	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	key, err := s.createUserKey(user)
	if err != nil {
		log.Errorln(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return key, nil
}

func (s *server) createUserKey(userID string) (*iam.Key, error) {
	priv, userKey, err := generateKey(rand.Reader, userID)
	if err != nil {
		return nil, err
	}

	if err := s.ds.CreateUserKey(userKey); err != nil {
		return nil, err
	}

	keyPB, err := userKey.Proto()
	if err != nil {
		return nil, err
	}
	keyPB.PrivateKeyData = priv

	return keyPB, nil
}

// GetKey gets a key for an authenticated user.
func (s *server) GetKey(ctx context.Context, req *iam.GetKeyRequest) (*iam.Key, error) {
	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	key, err := s.getUserKey(user, req.KeyId)
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

func (s *server) getUserKey(userID string, keyID string) (*iam.Key, error) {
	key, err := s.ds.GetUserKey(userID, keyID)
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
func (s *server) ListKeys(ctx context.Context, req *iam.ListKeysRequest) (*iam.ListKeysResponse, error) {
	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	keys, err := s.listUserKeys(user)
	if err != nil {
		log.Errorln(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &iam.ListKeysResponse{Keys: keys}, nil
}

func (s *server) listUserKeys(userID string) ([]*iam.Key, error) {
	keys, err := s.ds.ListUserKeys(userID)
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
func (s *server) DeleteKey(ctx context.Context, req *iam.DeleteKeyRequest) (*empty.Empty, error) {
	user, err := auth.UserFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.deleteUserKey(user, req.KeyId); err != nil {
		log.Errorln(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return new(empty.Empty), nil
}

func (s *server) deleteUserKey(userID string, keyID string) error {
	return s.ds.DeleteUserKey(userID, keyID)
}

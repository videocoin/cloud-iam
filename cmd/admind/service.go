package main

import (
	"context"

	iam "github.com/videocoin/cloud-iam/api/admin/v1"

	"github.com/jinzhu/gorm"
	"github.com/videocoin/cloud-iam/datastore"
	"github.com/videocoin/common/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	ds datastore.DataStore
}

// NewServer ...
func NewServer(ds datastore.DataStore) *server { return &server{ds: ds} }

// GetUserPublicKey
func (s *server) GetUserPublicKey(ctx context.Context, req *iam.GetUserPublicKeyRequest) (*iam.PublicKey, error) {
	key, err := s.ds.GetUserKey(req.UserId, req.KeyId)
	if err != nil {
		return nil, err
	}
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Debugln(err)
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Errorln(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &iam.PublicKey{PublicKeyData: key.PublicKeyData}, nil
}

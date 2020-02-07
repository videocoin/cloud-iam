package service

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"

	guuid "github.com/google/uuid"
	iam "github.com/videocoin/cloud-api/iam/v1"
	key "github.com/videocoin/cloud-pkg/api/resources/key"
	project "github.com/videocoin/cloud-pkg/api/resources/project"
	account "github.com/videocoin/cloud-pkg/api/resources/serviceaccount"

	"github.com/gogo/protobuf/types"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/cloud-iam/datastore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrPEMDataNotFound is returned when no PEM data is found.
var ErrPEMDataNotFound = errors.New("pem: PEM data not found")

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

// CreateServiceAccount creates a service account.
func (srv *Server) CreateServiceAccount(ctx context.Context, req *iam.CreateServiceAccountRequest) (*iam.ServiceAccount, error) {
	projName, err := project.ParseName(req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if ok := account.IsValidID(req.AccountId); !ok {
		return nil, status.Error(codes.InvalidArgument, account.ErrInvalidID.Error())
	}

	projID := projName.ID()
	acc, err := srv.ds.CreateServiceAccount(&datastore.ServiceAccount{
		ID:        guuid.New().String(),
		ProjectID: projID,
		Email:     account.NewEmail(projID, req.AccountId),
	})
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &iam.ServiceAccount{
		Name:      string(account.NewName(acc.ProjectID, acc.Email)),
		ProjectId: acc.ProjectID,
		UniqueId:  acc.ID,
		Email:     acc.Email,
	}, nil
}

// GetServiceAccount gets a service account.
func (srv *Server) GetServiceAccount(ctx context.Context, req *iam.GetServiceAccountRequest) (*iam.ServiceAccount, error) {
	name, err := account.ParseName(req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	acc, err := srv.ds.GetServiceAccountByEmail(name.Email())
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &iam.ServiceAccount{
		Name:      string(account.NewName(acc.ProjectID, acc.Email)),
		ProjectId: acc.ProjectID,
		UniqueId:  acc.ID,
		Email:     acc.Email,
	}, nil
}

// ListServiceAccounts lists service accounts.
func (srv *Server) ListServiceAccounts(ctx context.Context, req *iam.ListServiceAccountsRequest) (*iam.ListServiceAccountsResponse, error) {
	name, err := project.ParseName(req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accs, err := srv.ds.ListServiceAccounts(name.ID())
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	accsPB := make([]*iam.ServiceAccount, 0, len(accs))
	for _, acc := range accs {
		accsPB = append(accsPB, &iam.ServiceAccount{
			Name:      string(account.NewName(acc.ProjectID, acc.Email)),
			ProjectId: acc.ProjectID,
			UniqueId:  acc.ID,
			Email:     acc.Email,
		})
	}

	return &iam.ListServiceAccountsResponse{Accounts: accsPB}, nil
}

// DeleteServiceAccount deletes a service account.
func (srv *Server) DeleteServiceAccount(ctx context.Context, req *iam.DeleteServiceAccountRequest) (*types.Empty, error) {
	name, err := account.ParseName(req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := srv.ds.DeleteServiceAccount(name.Email()); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return new(types.Empty), nil
}

// CreateServiceAccountKey creates a service account key.
func (srv *Server) CreateServiceAccountKey(ctx context.Context, req *iam.CreateServiceAccountKeyRequest) (*iam.ServiceAccountKey, error) {
	name, err := account.ParseName(req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accEmail := name.Email()
	accKey, projID, err := srv.ds.CreateServiceAccountKey(name.Email(), srv.passphrase)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	validAfterTimePB, err := types.TimestampProto(accKey.ValidAfterTime)
	if err != nil {
		srv.logger.Error(accKey.ID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	validBeforeTimePB, err := types.TimestampProto(accKey.ValidBeforeTime)
	if err != nil {
		srv.logger.Error(accKey.ID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	block, _ := pem.Decode(accKey.PrivateKeyData)
	if block == nil {
		srv.logger.Error(accKey.ID, ErrPEMDataNotFound)
		return nil, status.Error(codes.Internal, ErrPEMDataNotFound.Error())
	}

	decrypted, err := x509.DecryptPEMBlock(block, []byte(srv.passphrase))
	if err != nil {
		srv.logger.Error(accKey.ID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	block.Bytes = decrypted

	return &iam.ServiceAccountKey{
		Name:            string(key.NewName(projID, accEmail, accKey.ID)),
		PrivateKeyData:  pem.EncodeToMemory(block),
		ValidAfterTime:  validBeforeTimePB,
		ValidBeforeTime: validAfterTimePB,
	}, nil
}

// GetServiceAccountKey gets a service account key.
func (srv *Server) GetServiceAccountKey(ctx context.Context, req *iam.GetServiceAccountKeyRequest) (*iam.ServiceAccountKey, error) {
	name, err := key.ParseName(req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accKey, err := srv.ds.GetServiceAccountKey(name.ID())
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	validAfterTimePB, err := types.TimestampProto(accKey.ValidAfterTime)
	if err != nil {
		srv.logger.Error(accKey.ID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	validBeforeTimePB, err := types.TimestampProto(accKey.ValidBeforeTime)
	if err != nil {
		srv.logger.Error(accKey.ID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	block, _ := pem.Decode(accKey.PrivateKeyData)
	if block == nil {
		srv.logger.Error(accKey.ID, ErrPEMDataNotFound)
		return nil, status.Error(codes.Internal, ErrPEMDataNotFound.Error())
	}

	decrypted, err := x509.DecryptPEMBlock(block, []byte(srv.passphrase))
	if err != nil {
		srv.logger.Error(accKey.ID, err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	block.Bytes = decrypted

	return &iam.ServiceAccountKey{
		Name:            req.Name,
		PrivateKeyData:  pem.EncodeToMemory(block),
		ValidAfterTime:  validAfterTimePB,
		ValidBeforeTime: validBeforeTimePB,
	}, nil
}

// ListServiceAccountKeys lists service account keys.
func (srv *Server) ListServiceAccountKeys(ctx context.Context, req *iam.ListServiceAccountKeysRequest) (*iam.ListServiceAccountKeysResponse, error) {
	accName, err := account.ParseName(req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	accEmail := accName.Email()
	accKeys, projID, err := srv.ds.ListServiceAccountKeysByEmail(accEmail)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Unknown, err.Error())
	}

	accKeysPB := make([]*iam.ServiceAccountKey, 0, len(accKeys))
	for _, accKey := range accKeys {
		validAfterTimePB, err := types.TimestampProto(accKey.ValidAfterTime)
		if err != nil {
			srv.logger.Error(accKey.ID, err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		validBeforeTimePB, err := types.TimestampProto(accKey.ValidBeforeTime)
		if err != nil {
			srv.logger.Error(accKey.ID, err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		block, _ := pem.Decode(accKey.PrivateKeyData)
		if block == nil {
			srv.logger.Error(accKey.ID, ErrPEMDataNotFound)
			return nil, status.Error(codes.Internal, ErrPEMDataNotFound.Error())
		}

		decrypted, err := x509.DecryptPEMBlock(block, []byte(srv.passphrase))
		if err != nil {
			srv.logger.Error(accKey.ID, err)
			return nil, status.Error(codes.Internal, err.Error())
		}
		block.Bytes = decrypted

		accKeysPB = append(accKeysPB, &iam.ServiceAccountKey{
			Name:            string(key.NewName(projID, accEmail, accKey.ID)),
			PrivateKeyData:  pem.EncodeToMemory(block),
			ValidAfterTime:  validBeforeTimePB,
			ValidBeforeTime: validAfterTimePB,
		})
	}

	return &iam.ListServiceAccountKeysResponse{Keys: accKeysPB}, nil
}

// DeleteServiceAccountKey deletes a service account key.
func (srv *Server) DeleteServiceAccountKey(ctx context.Context, req *iam.DeleteServiceAccountKeyRequest) (*types.Empty, error) {
	name, err := key.ParseName(req.Name)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := srv.ds.DeleteServiceAccountKey(name.ID()); err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return new(types.Empty), nil
}

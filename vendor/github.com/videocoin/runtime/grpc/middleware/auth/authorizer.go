package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

// RBAC handles role based access authorization.
func RBAC() AuthorizerFunc {
	authzCache := NewAuthzCache(AuthzCacheSize)

	// fix me: mock data
	mapMethodToPermission := map[string]string{
		"/videocoin.iam.v1.IAM/CreateKey": "iam.serviceAccountKeys.create",
		"/videocoin.iam.v1.IAM/ListKeys":  "iam.serviceAccountKeys.list",
		"/videocoin.iam.v1.IAM/GetKey":    "iam.serviceAccountKeys.get",
		"/videocoin.iam.v1.IAM/DeleteKey": "iam.serviceAccountKeys.delete",
	}
	roles := map[string]struct {
		IncludedPermissions []string
	}{
		"USER_ROLE_MINER": {
			IncludedPermissions: []string{
				"iam.serviceAccountKeys.create",
				"iam.serviceAccountKeys.get",
				"iam.serviceAccountKeys.list",
				"iam.serviceAccountKeys.delete",
			},
		},
		"USER_ROLE_SUPER": {
			IncludedPermissions: []string{
				"iam.serviceAccountKeys.create",
				"iam.serviceAccountKeys.get",
				"iam.serviceAccountKeys.list",
				"iam.serviceAccountKeys.delete",
			},
		},
	}

	return func(ctx context.Context, principal interface{}, fullMethod string) error {
		userInfo, ok := principal.(*UserInfo)
		if !ok {
			return errors.New("invalid principal")
		}

		var (
			tokenStr string
			err      error
		)
		if userInfo.HMACToken == "" {
			tokenStr, err = grpc_auth.AuthFromMD(ctx, "Bearer")
			if err != nil {
				return err
			}
		} else {
			tokenStr = userInfo.HMACToken
		}

		key := authzCache.ComposeKey(tokenStr, fullMethod)
		val, found := authzCache.Get(key)
		if found {
			if val.Success {
				return nil
			}
			return fmt.Errorf("permission %s is required to perform this operation on account %s", val.RequiredPermission, principal)
		}

		requiredPermission, found := mapMethodToPermission[fullMethod]
		if !found {
			return fmt.Errorf("permission not found")
		}

		userRole, err := getUserRole(tokenStr)
		if err != nil {
			return err
		}

		role, found := roles[userRole]
		if !found {
			return fmt.Errorf("role not found")
		}

		for _, permission := range role.IncludedPermissions {
			if permission == requiredPermission {
				authzCache.Add(key, &AuthzValue{Success: true})
				return nil
			}
		}

		authzCache.Add(key, &AuthzValue{Success: false, RequiredPermission: requiredPermission})

		return fmt.Errorf("permission %s is required to perform this operation on account %s", requiredPermission, principal)
	}
}

func getUserRole(tokenStr string) (string, error) {
	req, err := http.NewRequest("GET", "https://studio.dev.videocoin.network/api/v1/user", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenStr))

	res, err := cleanhttp.DefaultClient().Do(req)
	if err != nil {
		return "", err
	}

	props := map[string]interface{}{}
	if err := json.NewDecoder(res.Body).Decode(&props); err != nil {
		return "", fmt.Errorf("unable to decode JSON response: %v", err)
	}

	role, ok := props["role"]
	if !ok {
		return "", errors.New("role not available")
	}

	return role.(string), nil
}

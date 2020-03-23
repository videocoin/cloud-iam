package auth

import (
	"context"
	"fmt"
)

// AuthorizerFunc ...
type AuthorizerFunc func(ctx context.Context, principal interface{}, fullMethod string) error

// Authorize ...
func (f AuthorizerFunc) Authorize(ctx context.Context, principal interface{}, fullMethod string) error {
	return f(ctx, principal, fullMethod)
}

// Authorizer authorizes requests.
type Authorizer interface {
	Authorize(ctx context.Context, principal interface{}, fullMethod string) error
}

// RBACAuthZ ...
func RBACAuthZ() AuthorizerFunc {
	var (
		pstore PermissionStore = NewPermissionCache(NewPermissionRepo())
		rstore RoleStore       = NewRoleCache(NewRoleRepo())
	)

	return func(ctx context.Context, principal interface{}, fullMethod string) error {
		// maps method to permission
		requiredPermission, err := pstore.GetPermission(ctx, fullMethod)
		if err != nil {
			return err
		}

		// get user role
		userRole, err := rstore.GetUserRole(ctx, principal.(string))
		if err != nil {
			return err
		}

		// load role definition
		role, err := rstore.GetRole(ctx, userRole)
		if err != nil {
			return err
		}

		// verify if role has the required permission
		for _, permission := range role.IncludedPermissions {
			if permission == requiredPermission {
				return nil
			}
		}

		return fmt.Errorf("Permission %s is required to perform this operation on account %s", requiredPermission, principal)
	}
}

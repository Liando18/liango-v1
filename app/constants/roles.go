package constants

// Role constants — extend as needed.
const (
	RoleAdmin  = "admin"
	RoleUser   = "user"
	RoleEditor = "editor"
)

// Permission constants — extend as needed.
const (
	PermissionCreateAny = "create:any"
	PermissionReadAny   = "read:any"
	PermissionUpdateAny = "update:any"
	PermissionDeleteAny = "delete:any"
	PermissionReadOwn   = "read:own"
	PermissionUpdateOwn = "update:own"
)

// RolePermissions maps roles to their allowed permissions.
var RolePermissions = map[string][]string{
	RoleAdmin: {
		PermissionCreateAny,
		PermissionReadAny,
		PermissionUpdateAny,
		PermissionDeleteAny,
	},
	RoleEditor: {
		PermissionCreateAny,
		PermissionReadAny,
		PermissionUpdateOwn,
	},
	RoleUser: {
		PermissionReadOwn,
		PermissionUpdateOwn,
	},
}

// HasPermission checks if a role has a given permission.
func HasPermission(role, permission string) bool {
	perms, ok := RolePermissions[role]
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}

package models

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	RoleStaff Role = "staff"
)

// String returns the string representation of the Role
func (r Role) ToString() string {
	return string(r)
}

// IsValidRole checks if the role name is valid
func IsValidRole(roleName string) bool {
	validRoles := []Role{RoleAdmin, RoleUser, RoleStaff}
	for _, role := range validRoles {
		if role.ToString() == roleName {
			return true
		}
	}
	return false
}

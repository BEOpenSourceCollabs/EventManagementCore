package types

// Role representing the permission level of a user within the system.
type Role string

func (role Role) IsValid() bool {
	switch role {
	case UserRole, AdminRole, OrganizerRole:
		return true
	default:
		return false
	}
}

const (
	UserRole      Role = "user"
	AdminRole     Role = "admin"
	OrganizerRole Role = "organizer"
)

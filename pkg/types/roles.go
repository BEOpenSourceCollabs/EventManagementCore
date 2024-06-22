package types

// Role representing the permission level of a user within the system.
type Role string

const (
	UserRole      Role = "user"
	AdminRole     Role = "admin"
	OrganizerRole Role = "organizer"
)

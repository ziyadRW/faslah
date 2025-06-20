package auth

type UserType string

const (
	UserTypeCreator UserType = "creator"
	UserTypeViewer  UserType = "viewer"
	UserTypeAdmin   UserType = "admin"
)

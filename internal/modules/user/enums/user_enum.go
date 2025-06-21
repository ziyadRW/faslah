package user

type Type string

const (
	TypeCreator Type = "creator"
	TypeViewer  Type = "viewer"
	TypeAdmin   Type = "admin"
)

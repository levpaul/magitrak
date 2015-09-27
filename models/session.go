package models

const (
	SESSION_NAME = "magitrak"
)

type MagiSession struct {
	Authenticated bool
	UserId        int
}

package entities

type User struct {
	ID           int
	Mail         string
	Phone        string
	PasswordHash []byte
	Role         string
	Token        string
}

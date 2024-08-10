package entities

type User struct {
	ID       uint64
	Mail     string
	Phone    string
	PasswordHash []byte
	Role     string
	Token 	 string
}

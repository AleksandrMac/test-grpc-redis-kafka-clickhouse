package model

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type UserDB struct {
	ID    uint64
	Email string
}

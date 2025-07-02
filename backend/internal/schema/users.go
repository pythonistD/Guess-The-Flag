package schema

type Register struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

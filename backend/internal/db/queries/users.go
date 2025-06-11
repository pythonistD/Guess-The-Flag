package queries

var UserQueries = struct {
	Create     string
	GetByID    string
	GetByEmail string
}{
	Create: `
		INSERT INTO users (user_id, username, email, created_at)
		VALUES (:user_id, :username, :email, :created_at)
	`,
	GetByID: `
		SELECT user_id, username, email, created_at
		FROM users
		WHERE user_id = $1
	`,
	GetByEmail: `
		SELECT user_id, username, email, created_at
		FROM users
		WHERE email = $1
	`,
}

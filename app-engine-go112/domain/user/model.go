package user

// User ...
type User struct {
	ID    int      `db:"qwerty"`
	Name  string   `db:"asdfgh"`
	Foods []string `db:"zxcvbn"`
}

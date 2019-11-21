package user

// User ...
type User struct {
	ID    int64  `db:"qwerty"`
	Name  string `db:"asdfgh"`
	Foods []Food `db:"zxcvbn"`
}

// Food ...
type Food struct {
	ID   int64  `db:"edcrfv"`
	Name string `db:"qazwsx"`
}

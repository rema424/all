package handler

import "github.com/labstack/echo/v4"

// --------------------------------------------------
// LoginPage
// --------------------------------------------------

// LoginPage ...
func LoginPage(c echo.Context) error {
	return render(c, "signin.html", map[string]interface{}{})
}

// --------------------------------------------------
// LoginExec
// --------------------------------------------------

// LoginExec ...
func LoginExec(c echo.Context) error {
	return nil
}

// --------------------------------------------------
// IsLoggedIn
// --------------------------------------------------

// IsLoggedIn ...
func IsLoggedIn(c echo.Context) error {
	return nil
}

// --------------------------------------------------
// SignUpPage
// --------------------------------------------------

// SignUpPage ...
func SignUpPage(c echo.Context) error {
	return render(c, "signup.html", map[string]interface{}{})
}

// --------------------------------------------------
// SignUpExec
// --------------------------------------------------

// SignUpExec ...
func SignUpExec(c echo.Context) error {
	return nil
}

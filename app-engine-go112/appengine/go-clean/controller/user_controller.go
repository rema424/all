package controller

import (
	"fmt"
	"net/http"

	"myproject/domain/user"

	"github.com/labstack/echo/v4"
)

// UserController ...
type UserController struct {
	ui *user.Interactor
}

// NewUserController ...
func NewUserController(ui *user.Interactor) *UserController {
	return &UserController{ui}
}

// Register ...
func (uc *UserController) Register(c echo.Context) error {
	// curl -X POST localhost:8080/users -H 'Content-type: application/json' -d '{"name":"Alice", "foods":["apple", "banana"]}'
	in := struct {
		Name  string   `json:"name"`
		Foods []string `json:"foods"`
	}{}

	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println(in.Name, in.Foods)

	// ctx := c.Request().Context()
	u := user.User{
		Name:  in.Name,
		Foods: in.Foods,
	}
	// var err error

	// u, err = uc.ui.Register(ctx, u)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	return c.JSON(http.StatusOK, u)
}

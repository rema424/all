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
	fmt.Println("start user controller Register")

	// curl -X POST localhost:8080/users -H 'Content-type: application/json' -d '{"name":"Alice", "foods":["apple", "banana"]}'
	in := struct {
		Name  string   `json:"name"`
		Foods []string `json:"foods"`
	}{}

	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	fmt.Println(in.Name, in.Foods)

	foods := make([]user.Food, len(in.Foods))
	for i, food := range in.Foods {
		foods[i] = user.Food{Name: food}
	}

	u := user.User{
		Name:  in.Name,
		Foods: foods,
	}

	var err error
	ctx := c.Request().Context()
	u, err = uc.ui.Register(ctx, u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, u)
}

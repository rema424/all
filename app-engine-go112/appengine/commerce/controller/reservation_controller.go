package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HandleMakeReservation ...
func HandleMakeReservation(c echo.Context) error {
	type OrderDetail struct {
		ItemID   int `json:"itemId"`
		Quantity int `json:"quantity"`
	}

	type Customer struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
	}

	type Employee struct {
		ID int `json:"id"`
	}

	type In struct {
		OrderDetails []*OrderDetail `json:"orderDetails"`
		Customer     *Customer      `json:"customer"`
		Employee     *Employee      `json:"employee"`
		EmployeeID   int            `json:"employeeId"`
	}

	var in In

	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, in)
}

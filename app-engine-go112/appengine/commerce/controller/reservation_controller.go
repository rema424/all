package controller

import (
	"net/http"

	"myproject/domain/object"

	"github.com/labstack/echo/v4"
)

// HandleMakeReservation ...
func HandleMakeReservation(c echo.Context) error {
	var in struct {
		OrderDetails []struct {
			ItemID   int `json:"itemId"`
			Quantity int `json:"quantity"`
		} `json:"orderDetails"`
		Customer struct {
			Name        string `json:"name"`
			PhoneNumber string `json:"phoneNumber"`
		} `json:"customer"`
		Employee struct {
			ID int `json:"id"`
		} `json:"employee"`
	}

	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	orderDetails := make([]object.OrderDetail, len(in.OrderDetails))
	for i, d := range in.OrderDetails {
		orderDetails[i] = object.OrderDetail{
			Item:     object.Item{ID: d.ItemID},
			Quantity: d.Quantity,
		}
	}

	order := object.Order{
		OrderDetails: orderDetails,
		Customer: object.Customer{
			Name:        in.Customer.Name,
			PhoneNumber: in.Customer.PhoneNumber,
		},
		Employee: object.Employee{ID: in.Employee.ID},
	}

	return c.JSONPretty(http.StatusOK, order, "  ")
}

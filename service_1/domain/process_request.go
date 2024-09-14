package domain

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleRequest(c echo.Context, req interface{}) (bool, error) {
	if err := c.Bind(req); err != nil {
		return false, c.JSON(http.StatusBadRequest, ErrorRes{
			Status:  false,
			Message: "Binding failed. Invalid request",
			Error:   err.Error(),
		})
	}
	if err := c.Validate(req); err != nil {
		return false, c.JSON(http.StatusBadRequest, ErrorRes{
			Status:  false,
			Message: "Validation failed. Invalid request",
			Error:   err.Error(),
		})
	}
	return true, nil
}

package domain

import (
	"github.com/labstack/echo/v4"
)

func ErrorResponse(c echo.Context, statusCode int, msg string, err error) error {
	if err == nil {
		return c.JSON(statusCode, ErrorRes{
			Status:  false,
			Message: msg,
		})
	}
	return c.JSON(statusCode, ErrorRes{
		Status:  false,
		Message: msg,
		Error:   err.Error(),
	})
}

func DbErrorResponse(c echo.Context, err error) error {
	return ErrorResponse(c, 500, "Database error", err)
}
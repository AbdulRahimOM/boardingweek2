package main

import (
	"boarding-week2/service_1/config"
	"boarding-week2/service_1/controller"
	"boarding-week2/utils/validate"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()
	e.Validator = &validate.CustomValidator{Validator: validator.New()}
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	userGroup := e.Group("/user")
	{
		userGroup.GET("/:id", controller.GetUser)
		userGroup.POST("/create", controller.CreateUser)
		userGroup.PATCH("/update", controller.UpdateUser)
		userGroup.DELETE("/delete", controller.DeleteUser)
		userGroup.POST("/names", controller.GetNamesList)
	}

	e.GET("/ready", func(c echo.Context) error {
		return c.String(200, "Ready")
	})

	e.Logger.Fatal(e.Start(":" + config.EnvValues.Svc1Port))

}

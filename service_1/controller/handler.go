package controller

import (
	"boarding-week2/pb"
	"boarding-week2/service_1/config"
	"boarding-week2/service_1/domain"
	"fmt"
	"strconv"

	"log"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var handler struct {
	svc2Client pb.Svc2Client
}

func init() {
	svc2ClientConn, err := grpc.Dial(config.EnvValues.Svc2Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Grpc dial failed, error: %v", err)
	}

	handler.svc2Client = pb.NewSvc2Client(svc2ClientConn)
}

// Get names list
func GetNamesList(c echo.Context) error {
	req := new(domain.GetNamesReq)
	if ok, err := domain.HandleRequest(c, req); err != nil || !ok {
		return err
	}

	res, err := handler.svc2Client.Methods(c.Request().Context(), &pb.GetUserReq{
		Method:   req.Method,
		WaitTime: req.WaitTime,
	})
	if err != nil {
		return domain.ErrorResponse(c, 400, "Error", err)
	}

	return c.JSON(200, domain.GetUserNamesRes{
		Status:  true,
		Message: "Names list",
		Names:   res.Names,
	})
}

// GetUser
func GetUser(c echo.Context) error {
	idStr := c.Param("id")
	user := new(domain.User)
	ok, _ := retrieveUserFromCache(idStr, user)
	if ok {
		return c.JSON(200, domain.GetUserRes{
			Status:  true,
			Message: "User taken from cache",
			User:    *user,
		})
	}

	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return domain.ErrorResponse(c, 400, "Invalid id", err)
	}

	result := config.DB.First(user, idInt)
	if result.Error != nil {
		return domain.ErrorResponse(c, 400, "Db Error", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrorResponse(c, 400, "Error", fmt.Errorf("user not found"))
	}

	//store user in redis
	err = storeUserInCache(idStr, *user)
	if err != nil {
		fmt.Println("Error storing user in redis:", err) //continuing even if redis fails
	}

	return c.JSON(200, domain.GetUserRes{
		Status:  true,
		Message: "User taken from db",
		User:    *user,
	})
}

// CreateUser
func CreateUser(c echo.Context) error {
	req := new(domain.NewUserReq)
	if ok, err := domain.HandleRequest(c, req); err != nil || !ok {
		return err
	}

	user := domain.User{
		Name:  req.Name,
		Email: req.Email,
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		return domain.DbErrorResponse(c, result.Error)
	}

	//store user in redis
	err := storeUserInCache(strconv.Itoa(user.Id), user)
	if err != nil {
		fmt.Println("Error storing user in redis", err) //continuing even if redis fails
	}

	return c.JSON(200, domain.UserCreatedRes{
		Status:  true,
		Message: "User created successfully",
		UserId:  user.Id,
	})
}

// UpdateUser
func UpdateUser(c echo.Context) error {
	req := new(domain.UpdateUserReq)
	if ok, err := domain.HandleRequest(c, req); err != nil || !ok {
		return err
	}

	user := domain.User{
		Id:    req.Id,
		Name:  req.Name,
		Email: req.Email,
	}

	//check if user exists
	var count int64
	result := config.DB.Model(&domain.User{}).Where("id = ?", user.Id).Count(&count)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return domain.ErrorResponse(c, 400, "Error", fmt.Errorf("user not found"))
		} else {
			return domain.DbErrorResponse(c, result.Error)
		}
	}

	// initiate transaction
	tx := config.DB.Begin()
	if tx.Error != nil {
		return domain.DbErrorResponse(c, tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//update user
	result = tx.Model(&domain.User{}).Where("id = ?", user.Id).Updates(&user)
	if result.Error != nil {
		tx.Rollback()
		return domain.DbErrorResponse(c, result.Error)
	}

	//update/set user in redis
	err := storeUserInCache(strconv.Itoa(user.Id), user)
	if err != nil {
		tx.Rollback()
		return domain.ErrorResponse(c, 500, "Error storing user in redis. Hence, transaction rolled back in db", err)
	}

	tx.Commit()

	return c.JSON(200, domain.UserUpdatedRes{
		Status:  true,
		Message: "User updated successfully",
		User:    user,
	})
}

// DeleteUser
func DeleteUser(c echo.Context) error {
	req := new(domain.DeleteUserReq)
	if ok, err := domain.HandleRequest(c, req); err != nil || !ok {
		return err
	}

	user := domain.User{
		Id: req.Id,
	}

	//check if user exists
	var count int64
	result := config.DB.Model(&domain.User{}).Where("id = ?", user.Id).Count(&count)
	if result.Error != nil {
		return domain.DbErrorResponse(c, result.Error)
	}
	if count == 0 {
		return domain.ErrorResponse(c, 400, "Error", fmt.Errorf("user not found"))
	}

	// initiate transaction
	tx := config.DB.Begin()
	if tx.Error != nil {
		return domain.DbErrorResponse(c, tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//delete user from db
	result = tx.Delete(&user)
	if result.Error != nil {
		tx.Rollback()
		return domain.DbErrorResponse(c, result.Error)
	}

	//delete user from redis
	err := deleteUserFromCache(strconv.Itoa(user.Id))
	if err != nil {
		tx.Rollback()
		return domain.ErrorResponse(c, 500, "Error deleting user from redis. Hence, transaction rolled back in db", err)
	}

	tx.Commit()

	return c.JSON(200, domain.SuccessRes{
		Status:  true,
		Message: "User deleted successfully",
	})
}

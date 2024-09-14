package handler

import (
	"boarding-week2/pb"
	"boarding-week2/service_1/config"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type response struct {
	names []string
	err   error
}

type Handler struct {
	method1chan  chan int32
	responseChan chan response
	pb.UnimplementedSvc2Server
}

func NewHandler() *Handler {
	method1chan := make(chan int32)
	responseChan := make(chan response)
	db := config.DB
	go method1(method1chan, db, responseChan)

	return &Handler{
		method1chan:  method1chan,
		responseChan: responseChan,
	}
}

func (h *Handler) Methods(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserNamesResponse, error) {
	if req.Method == 1 {
		h.method1chan <- req.WaitTime
	} else if req.Method == 2 {
		res, err := method2(req.WaitTime, config.DB)
		if err != nil {
			return nil, err
		}
		return &pb.GetUserNamesResponse{
			Names: res,
		}, nil
	}

	res := <-h.responseChan
	if res.err != nil {
		return nil, res.err
	}
	return &pb.GetUserNamesResponse{
		Names: res.names,
	}, nil
}

func method1(method1chan chan int32, db *gorm.DB, responseChan chan response) {
	for {
		waitTime := <-method1chan
		var userNames []string
		result := db.Table("users").Select("name").Find(&userNames)

		time.Sleep(time.Duration(waitTime) * time.Second)
		if result.Error != nil {
			responseChan <- response{err: result.Error}
		} else {
			responseChan <- response{names: userNames}
		}
	}
}

func method2(waitTime int32, db *gorm.DB) ([]string, error) {
	for {
		fmt.Println("Method2 called with waitTime: ", waitTime)

		//get all user names
		var userNames []string
		result := db.Table("users").Select("name").Find(&userNames)

		time.Sleep(time.Duration(waitTime) * time.Second)
		if result.Error != nil {
			return nil, result.Error
		} else {
			return userNames, nil
		}

	}
}

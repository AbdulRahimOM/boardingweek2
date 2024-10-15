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
		fmt.Println(time.Now().Format("15:04:05"),"Method2 returned after WT:",req.WaitTime)
		return &pb.GetUserNamesResponse{
			Names: res,
		}, nil
	}else{
		return nil, fmt.Errorf("method-%d not supported", req.Method)
	}

	res := <-h.responseChan
	if res.err != nil {
		return nil, res.err
	}
	fmt.Println(time.Now().Format("15:04:05"),"Method1 returned after WT:",req.WaitTime)
	return &pb.GetUserNamesResponse{
		Names: res.names,
	}, nil
}

func method1(method1chan chan int32, db *gorm.DB, responseChan chan response) {
	for {
		waitTime := <-method1chan
		fmt.Println(time.Now().Format("15:04:05"),"Method1 called, WT: ", waitTime)
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
		fmt.Println(time.Now().Format("15:04:05"),"Method2 called, WT: ", waitTime)

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

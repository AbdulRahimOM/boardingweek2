package handler

import (
	"boarding-week2/pb"
	"boarding-week2/service_2/config"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type response struct {
	names []string
	err   error
}

type method1Req struct {
	waitTime     int32
	responseChan chan response
}

type Handler struct {
	method1chan chan method1Req
	pb.UnimplementedSvc2Server
}

func NewHandler() *Handler {
	method1chan := make(chan method1Req)
	db := config.DB
	go method1(method1chan, db)

	return &Handler{
		method1chan: method1chan,
	}
}

var (
	logCall = func(methodNum, waitTime int32) {
		fmt.Println(time.Now().Format("15:04:05"), "Method", methodNum, "called, WT: ", waitTime)
	}

	logReturn = func(methodNum, waitTime int32) {
		fmt.Println(time.Now().Format("15:04:05"), "Method", methodNum, "returned after WT:", waitTime)
	}
)

func (h *Handler) Methods(ctx context.Context, req *pb.GetUserReq) (*pb.GetUserNamesResponse, error) {
	logCall(req.Method, req.WaitTime)
	defer logReturn(req.Method, req.WaitTime)

	var res response

	if req.Method == 1 {
		responseChan := make(chan response)
		h.method1chan <- method1Req{
			waitTime:     req.WaitTime,
			responseChan: responseChan,
		}
		res = <-responseChan
	} else {
		res = method2(req.WaitTime, config.DB)
	}

	if res.err != nil {
		return nil, res.err
	} else {
		return &pb.GetUserNamesResponse{
			Names: res.names,
		}, nil
	}
}

func method1(method1chan chan method1Req, db *gorm.DB) {
	for {
		req := <-method1chan
		var userNames []string
		result := db.Table("users").Select("name").Find(&userNames)

		time.Sleep(time.Duration(req.waitTime) * time.Second)
		if result.Error != nil {
			req.responseChan <- response{err: result.Error}
		} else {
			req.responseChan <- response{names: userNames}
		}
	}
}

func method2(waitTime int32, db *gorm.DB) response {
		//get all user names
		var userNames []string
		result := db.Table("users").Select("name").Find(&userNames)

		time.Sleep(time.Duration(waitTime) * time.Second)
		if result.Error != nil {
			return response{err: result.Error}
		} else {
			return response{names: userNames}
		}
}

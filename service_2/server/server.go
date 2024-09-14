package server

import (
	"boarding-week2/service_2/handler"

	pb "boarding-week2/pb"
)

func InitializeServer() pb.Svc2Server {
	return handler.NewHandler()
}

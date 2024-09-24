package test

import (
	"context"
	"github.com/msprojectlb/project-common/grpc/test/proto"
	"math/rand"
)

type AppServer struct {
	proto.UnimplementedTestServiceServer
}

func (a *AppServer) Hello(ctx context.Context, req *proto.HelloReq) (*proto.HelloRes, error) {
	name := req.Name
	return &proto.HelloRes{
		Id:       int32(rand.Intn(2000)),
		UserName: name,
	}, nil
}

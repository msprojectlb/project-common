package test

import (
	"context"
	"github.com/msprojectlb/project-common/mygrpc/test/gen"
	"math/rand"
)

type AppServer struct {
	gen.UnimplementedAppServiceServer
}

func (a *AppServer) Hello(ctx context.Context, req *gen.HelloReq) (*gen.HelloRes, error) {
	name := req.Name
	return &gen.HelloRes{
		Id:       int32(rand.Intn(2000)),
		Username: name,
	}, nil
}

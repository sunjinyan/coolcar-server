package server

import (
	"coolcar/shared/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GRPCConfig struct {
	Logger *zap.Logger
	Addr string
	Name string
	AuthPublicKeyFile string
	RegisterFunc func(server *grpc.Server)
}

func RunGRPCServer(c *GRPCConfig) (err error) {
	nameField := zap.String("name",c.Name)

	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		c.Logger.Fatal("can not listen",nameField,zap.Error(err))
		return
	}

	//添加拦截器
	var in grpc.UnaryServerInterceptor
	var opts []grpc.ServerOption
	if c.AuthPublicKeyFile != ""{
		in,err = auth.Interceptor(c.AuthPublicKeyFile)
		if err != nil {
			c.Logger.Fatal("can not Interceptor",nameField ,zap.Error(err))
			return
		}
		opts = append(opts,grpc.UnaryInterceptor(in))
	}

	//var s *grpc.Server
	//if in != nil {
	//	s = grpc.NewServer()
	//}else{
	//	s = grpc.NewServer(grpc.UnaryInterceptor(in))
	//}
	s := grpc.NewServer(opts...)
	c.RegisterFunc(s)

	//rentalpb.RegisterTripServiceServer(s,&trip.Service{
	//	Logger:                         logger,
	//})
	//err = s.Serve(lis)
	c.Logger.Sugar().Info("server started",nameField,zap.String("addr",c.Addr))
	return s.Serve(lis)
	//c.Logger.Fatal("can not server ",zap.Error(err))
}
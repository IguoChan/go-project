package grpcx

type GrpcServer interface {
	Serve() error
	Stop()
}

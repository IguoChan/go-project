package rpcx

type GrpcServer interface {
	Serve() error
	Stop()
}

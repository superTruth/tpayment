package grpc_pool

import "google.golang.org/grpc"

func GetConn(url string) (*grpc.ClientConn, error) {
	return grpc.Dial(url, grpc.WithInsecure())
}

func PutConn(url string, conn *grpc.ClientConn) {
	_ = conn.Close()
}

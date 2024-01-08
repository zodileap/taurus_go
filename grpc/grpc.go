package grpc

import (
	"context"
	"errors"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Server interface {
	Register(gRPC *grpc.Server) error
}

type Client interface {
}

var (
	Mu      sync.RWMutex
	servers = make(map[string]map[string]Server)
	clients = make(map[string]*grpc.ClientConn)
)

// 注册一个GRPC下的服务
//
// gRPC作为主服务, 一个gRPC下可以有多个不同功能服务
//
// 参数:
//   - grpcName: GRPC名
//   - serverName: GRPC下的服务名字
//   - server:  GRPC下的服务
func RegisterServer(grpcName string, serverName string, server Server) {
	Mu.Lock()
	defer Mu.Unlock()
	if server == nil {
		panic(text_panic_server_nil)
	}
	if _, ok := servers[grpcName]; !ok {
		servers[grpcName] = make(map[string]Server)
	}
	if _, dup := servers[grpcName][serverName]; dup {
		panic(text_panic_server_register_twice + serverName)
	}

	servers[grpcName][serverName] = server
}

// 初始化GRPC服务
//
// 参数:
//   - grpcName: GRPC服务名
//   - keyFile: TLS私钥,如果为空则不使用TLS
//   - certFile: TLS证书，如果为空则不使用TLS
func InitServer(grpcName string, keyFile string, certFile string) *grpc.Server {
	Mu.Lock()
	defer Mu.Unlock()

	if _, ok := servers[grpcName]; !ok {
		panic(text_panic_grpc_not_register)
	}
	var grpcSrv *grpc.Server

	if keyFile != "" && certFile != "" {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			panic(text_panic_grpc_tls_fail + err.Error())
		}
		grpcSrv = grpc.NewServer(grpc.Creds(creds))
	} else {
		grpcSrv = grpc.NewServer()
	}
	for _, server := range servers[grpcName] {
		server.Register(grpcSrv)
	}

	return grpcSrv
}

// 注册GRPC客户端
//
// 参数:
//   - clientName: 客户端名
//   - target: 客户端地址
func RegisterClient(clientName string, target string, certFile string) {
	Mu.Lock()
	defer Mu.Unlock()
	if _, ok := clients[clientName]; ok {
		panic(text_panic_client_connect_register_twice + clientName)
	}
	var creds credentials.TransportCredentials

	if certFile != "" {
		var err error
		creds, err = credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			panic(text_panic_grpc_tls_fail + err.Error())
		}
	} else {
		creds = insecure.NewCredentials()
	}

	// 设置连接选项
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		// 添加重连策略
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1.0 * time.Second,
				Multiplier: 1.5,
				Jitter:     0.2,
				MaxDelay:   time.Minute,
			},
			MinConnectTimeout: 5 * time.Second,
		}),
	}

	// 创建一个不会超时的上下文
	ctx := context.Background()

	// // 创建一个带有超时的上下文
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// 使用带有超时的上下文来拨号
	conn, err := grpc.DialContext(ctx, target, opts...)

	// conn, err := grpc.Dial(target, grpc.WithTransportCredentials(creds))

	if err != nil {
		panic(text_panic_client_connect_fail + clientName)
	}

	clients[clientName] = conn

}

// 获取GRPC客户端
//
// 参数:
//   - clientName: 客户端名
func GetClient(clientName string) (*grpc.ClientConn, error) {
	Mu.Lock()
	defer Mu.Unlock()
	if _, ok := clients[clientName]; !ok {
		return nil, errors.New(text_err_get_client_fail + clientName)
	}

	return clients[clientName], nil
}

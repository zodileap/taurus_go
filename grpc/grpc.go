package grpc

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Server 定义 gRPC 服务注册器。
type Server interface {
	Register(gRPC *ggrpc.Server) error
}

// Manager 管理 gRPC server/client 的注册与连接。
type Manager struct {
	mu      sync.RWMutex
	servers map[string]Server
	clients map[string]*ggrpc.ClientConn
}

// NewManager 创建一个新的管理器实例。
func NewManager() *Manager {
	return &Manager{
		servers: make(map[string]Server),
		clients: make(map[string]*ggrpc.ClientConn),
	}
}

// RegisterServer 注册一个逻辑 gRPC 服务。
func (m *Manager) RegisterServer(name string, server Server) error {
	if server == nil {
		return errors.New("grpc server cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.servers[name]; exists {
		return fmt.Errorf("grpc server already registered: %s", name)
	}

	m.servers[name] = server
	return nil
}

// InitServer 根据已注册服务创建 gRPC Server。
func (m *Manager) InitServer(name string, keyFile string, certFile string) (*ggrpc.Server, error) {
	m.mu.RLock()
	server, exists := m.servers[name]
	m.mu.RUnlock()
	if !exists {
		return nil, fmt.Errorf("grpc server not registered: %s", name)
	}

	var grpcServer *ggrpc.Server
	if keyFile != "" && certFile != "" {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("create grpc tls credentials for %s: %w", name, err)
		}
		grpcServer = ggrpc.NewServer(ggrpc.Creds(creds))
	} else {
		grpcServer = ggrpc.NewServer()
	}

	if err := server.Register(grpcServer); err != nil {
		return nil, fmt.Errorf("register grpc service %s: %w", name, err)
	}

	return grpcServer, nil
}

// RegisterClient 建立并缓存一个客户端连接。
func (m *Manager) RegisterClient(ctx context.Context, name string, target string, certFile string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.clients[name]; exists {
		return fmt.Errorf("grpc client already registered: %s", name)
	}

	var creds credentials.TransportCredentials
	if certFile != "" {
		var err error
		creds, err = credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			return fmt.Errorf("create grpc client tls credentials for %s: %w", name, err)
		}
	} else {
		creds = insecure.NewCredentials()
	}

	conn, err := ggrpc.DialContext(
		ctx,
		target,
		ggrpc.WithTransportCredentials(creds),
		ggrpc.WithBlock(),
		ggrpc.WithConnectParams(ggrpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1 * time.Second,
				Multiplier: 1.5,
				Jitter:     0.2,
				MaxDelay:   time.Minute,
			},
			MinConnectTimeout: 5 * time.Second,
		}),
	)
	if err != nil {
		return fmt.Errorf("dial grpc client %s (%s): %w", name, target, err)
	}

	m.clients[name] = conn
	return nil
}

// GetClient 返回已注册的客户端连接。
func (m *Manager) GetClient(name string) (*ggrpc.ClientConn, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, exists := m.clients[name]
	if !exists {
		return nil, fmt.Errorf("grpc client not registered: %s", name)
	}

	return conn, nil
}

// Close 关闭所有已注册客户端连接。
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for name, conn := range m.clients {
		if err := conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close grpc client %s: %w", name, err))
		}
		delete(m.clients, name)
	}

	return errors.Join(errs...)
}

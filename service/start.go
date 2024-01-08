// Description: 关于启动服务的模块，包括http和grpc服务
package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

var (
	mu        sync.Mutex
	services  = make(map[string]bool)
	srvCancel = make(map[string]context.CancelFunc)
)

type ServiceConfig struct {
	HTTP []ServiceHTTPConfig
	GRPC []ServiceGRPCConfig
}

type ServiceHTTPConfig struct {
	Name    string
	Address string
	Handler http.Handler
}

type ServiceGRPCConfig struct {
	Name    string
	Address string
	Server  *grpc.Server
}

var stopSignal = make(chan struct{})

// 启动服务
func StartService(config ServiceConfig) {
	var wg sync.WaitGroup
	duplicateDetected := false

	// Handle for each HTTP server
	for _, httpConfig := range config.HTTP {
		fmt.Printf("启动http服务: %s\n", httpConfig.Name)
		if duplicateDetected {
			break
		}
		wg.Add(1)

		if serviceExist(httpConfig.Name) {
			fmt.Println(text_panic_http_server_exist)
			duplicateDetected = true
			continue
		}

		go func(config ServiceHTTPConfig) {
			defer wg.Done()
			srv := &http.Server{
				Addr:    config.Address,
				Handler: config.Handler,
			}

			go func() {
				if err := srv.ListenAndServe(); err != http.ErrServerClosed {
					fmt.Println("ListenAndServe Error:", err.Error())
				}
			}()

			ctx, cancel := context.WithCancel(context.Background())
			srvCancel[config.Name] = cancel

			// 通过select来同时监听ctx和stopSignal
			select {
			case <-ctx.Done():
				shutdownCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
				if err := srv.Shutdown(shutdownCtx); err != nil {
					fmt.Println("Shutdown Error:", err.Error())
				}
			case <-stopSignal:
				shutdownCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
				if err := srv.Shutdown(shutdownCtx); err != nil {
					fmt.Println("Shutdown Error:", err.Error())
				}
			}
		}(httpConfig)

	}
	if duplicateDetected {
		close(stopSignal)
	}

	// Handle for each gRPC server
	for _, grpcConfig := range config.GRPC {
		fmt.Printf("启动gRPC服务: %s\n", grpcConfig.Name)
		if duplicateDetected {
			break
		}

		wg.Add(1)
		if serviceExist(grpcConfig.Name) {
			fmt.Println(text_panic_rpc_server_exist)
			duplicateDetected = true
			continue
		}
		go func(config ServiceGRPCConfig) {
			defer wg.Done()

			lis, err := net.Listen("tcp", config.Address)
			if err != nil {
				fmt.Println("failed to listen:", err)
				return
			}

			if err := config.Server.Serve(lis); err != nil {
				fmt.Println("failed to serve:", err)
			}

			ctx, cancel := context.WithCancel(context.Background())
			srvCancel[config.Name] = cancel

			// 启动一个goroutine来处理服务器关闭
			go func() {
				select {
				case <-ctx.Done():
					// Here, you may want to use grpc's GracefulStop() or Stop() instead of Shutdown()
					// since gRPC servers don't have a Shutdown() method.
					config.Server.GracefulStop()
				case <-stopSignal:
					config.Server.GracefulStop()
				}
			}()

		}(grpcConfig)
	}

	if duplicateDetected {
		close(stopSignal)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	mu.Lock()
	for name, cancel := range srvCancel {
		fmt.Printf("Stopping service %s\n", name)
		cancel()
	}
	for _, grpcConfig := range config.GRPC {
		grpcConfig.Server.GracefulStop()
	}
	mu.Unlock()

	wg.Wait()
}

func serviceExist(name string) bool {
	mu.Lock()
	defer mu.Unlock()
	if services[name] {
		return true
	}
	services[name] = true
	return false
}

package tcp

import (
	"context"
	"go-redis/interface/tcp"
	"go-redis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {

	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal)
	// 监听如下退出信号, 并发送给sigChan
	// 终端断开链接, ctrl+/, 结束程序, ctrl+c
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	// 监听退出信号, 并发送结束信号
	go func() {
		sig := <-sigChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info("start listen")
	ListenAndServe(listener, handler, closeChan)
	return nil
}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	// 保证接收到退出信号时会关闭listener和handler
	go func() {
		<-closeChan
		logger.Info("shutting down")
		// 关闭socket监听
		_ = listener.Close()
		// 关闭handler处理
		_ = handler.Close()
	}()

	// 保证退出时会关闭listener和handler
	defer func() {
		// 关闭socket监听
		_ = listener.Close()
		// 关闭handler处理
		_ = handler.Close()
	}()

	ctx := context.Background()
	var waitDone sync.WaitGroup
	for true {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accepted link")
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}

	// 等待所有的socket处理完毕
	waitDone.Wait()
}

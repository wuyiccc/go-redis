package tcp

import (
	"bufio"
	"context"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}

func (e *EchoClient) Close() error {
	// time wait 10s
	e.Waiting.WaitWithTimeout(10 * time.Second)
	_ = e.Conn.Close()
	return nil
}

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		_ = conn.Close()
	}
	client := &EchoClient{
		Conn: conn,
	}
	h.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("connection close")
				h.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		// waitGroup+1, 防止当发送信息到一半时, Server关闭连接
		client.Waiting.Add(1)
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}

}

// Close stops echo handler
func (h *EchoHandler) Close() error {
	logger.Info("handler shutting down...")
	h.closing.Set(true)
	// 循环关闭client
	h.activeConn.Range(func(key, value interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})
	return nil
}

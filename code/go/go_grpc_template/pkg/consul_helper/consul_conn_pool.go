package consul_helper

import (
	"context"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type ClientPool struct {
	consulResolver ResolverNode
	connPool       sync.Map
	timeout        time.Duration
}

type ConnWithTs struct {
	UpdateTime int64
	Conn       *grpc.ClientConn
}

type ResolverNode interface {
	NodeAddress() (string, error)
}

func NewClientPool(resolver ResolverNode, timeout time.Duration) *ClientPool {
	clientPool := &ClientPool{}
	clientPool.consulResolver = resolver
	clientPool.timeout = timeout
	clientPool.InitPool()
	return clientPool
}

func (pool *ClientPool) InitPool() {

	// with initial pool size of 10
	for i := 0; i < 10; i++ {
		conn, addr, err := pool.NewConnect()
		if err != nil {
			continue
		}
		pool.connPool.Store(addr, &ConnWithTs{time.Now().Unix(), conn})
	}
	go pool.watch()
}

func (pool *ClientPool) watch() {
	// every 30 second
	// delete unused connection
	tick := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-tick.C:
			now := time.Now().Unix()
			pool.connPool.Range(func(key, val interface{}) bool {
				if connWithTs, ok := val.(*ConnWithTs); ok {
					if now-connWithTs.UpdateTime > 30 {
						connWithTs.Conn.Close()
						pool.connPool.Delete(key)
					}
				}
				return true
			})
		}
	}
	return
}

func (pool *ClientPool) Get() (*grpc.ClientConn, error) {
	var err error
	addr, err := pool.consulResolver.NodeAddress()
	if err != nil {
		return nil, err
	}

	// 1st: use the picked adrress from pool connection
	if val, ok := pool.connPool.Load(addr); ok {
		connWithTs := val.(*ConnWithTs)
		connWithTs.UpdateTime = time.Now().Unix()
		return connWithTs.Conn, nil
	}

	var conn *grpc.ClientConn

	// 2nd: use the picked address for new connection
	conn, err = pool.NewConnectWithAddr(addr)
	if err != nil {
		// 3rd: create a new connection
		conn, addr, err = pool.NewConnect()
	}

	// use load or store for repeated write
	val, loaded := pool.connPool.LoadOrStore(addr, &ConnWithTs{time.Now().Unix(), conn})

	if loaded {
		conn.Close()
	}

	connWithTs := val.(*ConnWithTs)
	return connWithTs.Conn, nil
}

func (pool *ClientPool) NewConnectWithAddr(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), pool.timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithBlock(), grpc.WithInsecure())
	return conn, err

}

func (pool *ClientPool) NewConnect() (*grpc.ClientConn, string, error) {
	retry := 0
	var err error
	// new with retry
	for {
		retry++
		if retry > 3 {
			break
		}

		addr, err := pool.consulResolver.NodeAddress()
		if addr == "" {
			continue
		}

		conn, err := pool.NewConnectWithAddr(addr)
		if err == nil {
			return conn, addr, err
		}

	}
	return nil, "", err
}

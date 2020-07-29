# Go 如何使得 Web 工作

下图是Go实现Web服务的工作模式的流程图:

![http包执行流程.png](https://upload-images.jianshu.io/upload_images/14738618-02403d4816dce673.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

golang 的标准库 net/http 提供了 http 编程有关的接口，封装了内部TCP连接和报文解析的复杂琐碎的细节，使用者只需要和 http.request 和 http.ResponseWriter 两个对象交互就行。
当每次客户端有请求的时候，把请求封装成 http.Request ，调用对应的 handler 的 ServeHTTP 方法，然后把操作后的 http.ResponseWriter 解析，返回到客户端。
也就是说，我们只要写一个 handler，请求会通过参数传递进来，而它要做的就是根据请求的数据做处理，把结果写到 Response 中。

```go
// 这就是用Go实现的一个最简短的hello world服务器.
package main

import (
	"io"
	"net/http"
)

func main() {
	// 注册函数HandleFunc，用户连接， 自动调用指定处理函数
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		// 给客户端回复数据
		io.WriteString(w, "hello, world!\n")
		w.Write([]byte("lisa"))
	})
	http.ListenAndServe(":3000", nil)
}
```

## 1.路由注册函数HandleFunc()

http中的 HandleFunc方法,主要用来注册路由。http.HandleFunc 接受两个参数：第一个参数是字符串表示的 url 路径，第二个参数是该 url 实际的处理对象。

```$go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
    DefaultServeMux.HandleFunc(pattern, handler)
}
```
ServeMux(DefaultServeMux)维护了一个路由拓扑表。Mux是 multiplexor 的缩写，就是多路传输的意思，请求传过来，根据某种判断，分流到后端多个不同的地方。ServeMux主要通过`map[string]muxEntry`，来存储了具体的url模式和handler（此handler是实现Handler接口的类型）。

```$go
type ServeMux struct {
    mu    sync.RWMutex
    m     map[string]muxEntry	//存放路由信息的字典！
    hosts bool 
}

type muxEntry struct {
    h        Handler
    pattern  string
}
```

### 1.1.Handler接口，HTTP字节流处理函数

Handler接口可以算是HTTP Server的一个枢纽（核心处理逻辑）。Handler 接口都要实现 ServeHTTP 这个方法。

```$go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

看到这里，我们知道用来注册路由的 HandleFunc方法其实是ServeMux的方法HandleFunc。

```$go
func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	mux.Handle(pattern, HandlerFunc(handler))
}

func (mux *ServeMux) Handle(pattern string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	//边界情况处理。pattern == "" || handler == nil
	// mux.m == nil, make(map[string]muxEntry)
	// ...
	
	//把 handler 和 pattern 模式绑定到 map[string]muxEntry 上
	mux.m[pattern] = muxEntry{h: handler, pattern: pattern}

	if pattern[0] != '/' {
		mux.hosts = true
	}

	//这是一个很有用的小技巧，如果注册了 `/tree/`， `serveMux` 会自动添加一个 `/tree` 的路径并重定向到 `/tree/`。当然这个 `/tree` 路径会被用户显示的路由信息覆盖。。
	n := len(pattern)
	if n > 0 && pattern[n-1] == '/' &&
		_, exist := mux.m[pattern[0:n-1]];!exist{

		path := pattern
		if pattern[0] != '/'{
		    path = pattern[strings.Index(pattern, "/"):]
		}
		url := &url.URL{Path: path}
		mux.m[pattern[0:n-1]] = muxEntry{h: RedirectHandler(url.String(), StatusMovedPermanently), pattern: pattern}
	}
}
```

## 2.服务监听函数ListenAndServe()

http中的 ListenAndServe方法,主要用来监听服务。http.ListenAndServe接受两个参数：第一个参数是字符串表示的监听地址（包含端口号），第二个参数是Handler 类型的处理对象。

```$go
func ListenAndServe(addr string, handler Handler) error {
	//创建一个Server结构体，调用该结构体的ListenAndServer方法然后返回。Server 保存了运行 HTTP 服务需要的参数。
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
```

接下来继续看看Server.ListenAndServe里面都做了些什么，ListenAndServe调用 net.Listen 监听在对应的 tcp 端口，tcpKeepAliveListener 设置了 TCP 的 KeepAlive 功能，最后调用 srv.Serve()方法开始真正的循环逻辑。

```$go
//初始化监听地址Addr，同时调用Listen方法设置监听。
//最后将监听的TCP对象传入Serve方法：
func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		//如果不指定服务器地址信息，默认以":http"作为地址信息(等价于":80")
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

//Serve接受 Listener l 传递过来的请求，为每个请求创建 goroutine 进行后台处理，goroutine 会读取请求，调用 srv.Handler。
func (srv *Server) Serve(l net.Listener) error {
	if fn := testHookServerServe; fn != nil {
		fn(srv, l)
	}

	l = &onceCloseListener{Listener: l}
	defer l.Close()

	if err := srv.setupHTTP2_Serve(); err != nil {
		return err
	}

	if !srv.trackListener(&l, true) {
		return ErrServerClosed
	}
	defer srv.trackListener(&l, false)

	var tempDelay time.Duration     // how long to sleep on accept failure

	baseCtx := context.Background() // base is always background, per Issue 16220
	ctx := context.WithValue(baseCtx, ServerContextKey, srv)

	//开启循环进行监听，通过传进来的listener接收来自客户端的请求并建立连接，然后为每一个连接创建routine执行具体的服务处理c.serve()
	for {
		//通过Listener的Accept方法用来获取连接数据
		rw, e := l.Accept()
		if e != nil {
			select {
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				srv.logf("http: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		tempDelay = 0
		//等待客户端请求。通过获得的连接数据，创建newConn连接对象
		c := srv.newConn(rw)
		c.setState(c.rwc, StateNew) // before Serve can return
		//开启goroutine发送连接请求
		go c.serve(ctx)
	}
}

func (c *conn) serve(ctx context.Context) {
	c.remoteAddr = c.rwc.RemoteAddr().String()
	//连接关闭相关的处理
	defer func() {
		if err := recover(); err != nil && err != ErrAbortHandler {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			c.server.logf("http: panic serving %v: %v\n%s", c.remoteAddr, err, buf)
		}
		if !c.hijacked() {
			c.close()
			c.setState(c.rwc, StateClosed)
		}
	}()

	.....
	ctx, cancelCtx := context.WithCancel(ctx)
	c.cancelCtx = cancelCtx
	defer cancelCtx()
	c.r = &connReader{conn: c}
	//handler实际操作的读写对象
	c.bufr = newBufioReader(c.r)
	c.bufw = newBufioWriterSize(checkConnErrorWriter{c}, 4<<10)

	for {
		//读取客户端的请求
		w, err := c.readRequest(ctx)
		if c.r.remain != c.server.initialReadLimitSize() {
			// If we read any bytes off the wire, we're active.
			c.setState(c.rwc, StateActive)
		}
		.................
		//处理网络数据的状态
		// Expect 100 Continue support
		req := w.req
		if req.expectsContinue() {
			if req.ProtoAtLeast(1, 1) && req.ContentLength != 0 {
				// Wrap the Body reader with one that replies on the connection
				req.Body = &expectContinueReader{readCloser: req.Body, resp: w}
			}
		} else if req.Header.get("Expect") != "" {
			w.sendExpectationFailed()
			return
		}


		c.curReq.Store(w)

		if requestBodyRemains(req.Body) {
			registerOnHitEOF(req.Body, w.conn.r.startBackgroundRead)
		} else {
			if w.conn.bufr.Buffered() > 0 {
				w.conn.r.closeNotifyFromPipelinedRequest()
			}
			w.conn.r.startBackgroundRead()
		}

		//请求处理方法，调用最早传递给 Server 的 Handler 函数
		serverHandler{c.server}.ServeHTTP(w, w.req)
		
		w.cancelCtx()
		if c.hijacked() {
			return
		}
		w.finishRequest()
		if !w.shouldReuseConnection() {
			if w.requestBodyLimitHit || w.closedRequestBodyEarly() {
				c.closeWriteAndWait()
			}
			return
		}
		c.setState(c.rwc, StateIdle)
		c.curReq.Store((*response)(nil))

		if !w.conn.server.doKeepAlives() {
			return
		}
		if d := c.server.idleTimeout(); d != 0 {
			c.rwc.SetReadDeadline(time.Now().Add(d))
			if _, err := c.bufr.Peek(4); err != nil {
				return
			}
		}
		c.rwc.SetReadDeadline(time.Time{})
	}
}
```

## 几个HTTP服务器例子

### 1.实现简单的tomcat 服务器

```go
package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	//http请求处理
	http.HandleFunc("/", handler1)
	//绑定监听地址和端口
	http.ListenAndServe("localhost:8080", nil)
}

//请求处理函数
func handler1(w http.ResponseWriter, r *http.Request) {
	//获取请求资源
	path := r.URL.Path
	if strings.Contains(path[1:], "") {
		//返回请求资源
		fmt.Fprintf(w, getHtmlFile("index.html"))
	} else {
		if strings.Contains(path[1:], ".html") {
			w.Header().Set("content-type", "text/html")
			fmt.Fprintf(w, getHtmlFile(path[1:]))
		}
		if strings.Contains(path[1:], ".css") {
			w.Header().Set("content-type", "text/css")
			fmt.Fprintf(w, getHtmlFile(path[1:]))
		}
		if strings.Contains(path[1:], ".js") {
			w.Header().Set("content-type", "text/javascript")
			fmt.Fprintf(w, getHtmlFile(path[1:]))
		}
		if strings.Contains(path[1:], "") {
			fmt.Print(strings.Contains(path[1:], ""))
		}
	}

}

func getHtmlFile(path string) (fileHtml string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n')

		if err != nil || io.EOF == err {
			break
		}
		fileHtml += line
	}
	return fileHtml
}
```


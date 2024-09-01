package main

import (
	"flag"
	"fmt"
	"gomicro/Services"
	"gomicro/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	mymux "github.com/gorilla/mux"
	"golang.org/x/time/rate"
)

func main() {
	go func() {}()
	go (func() {})()
	name := flag.String("name", "", "服务名称")
	port := flag.Int("port", 0, "服务端口")
	flag.Parse()
	if *name == "" {
		log.Fatal("请指定服务名")
	}
	if *port == 0 {
		log.Fatal("请指定端口")
	}
	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stdout)
		logger = kitlog.WithPrefix(logger, "mykit", "1.0")
		logger = kitlog.WithPrefix(logger, "time", kitlog.DefaultTimestampUTC) //加上前缀时间
		logger = kitlog.WithPrefix(logger, "caller", kitlog.DefaultCaller)     //加上前缀，日志输出时的文件和第几行代码

	}
	utils.SetServiceNameAndPort(*name, *port) //设置服务名和端口

	user := Services.UserService{}
	limit := rate.NewLimiter(1, 5)
	endp := Services.RateLimit(limit)(Services.UserServiceLogMiddleware(logger)(Services.GenUserEnPoint(user)))
	/*我们分析一下上面这段代码Services.RateLimit(limit)返回一个Middware，type Middleware func(Endpoint) Endpoint
	  也就是说这段代码的返回值必须是Endpoint类型
	  type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
	  才可以传入Middware (Services.UserServiceLogMiddleware(logger)(Services.GenUserEnPoint(user)))

	  再拆分Services.UserServiceLogMiddleware(logger)也返回一个Middware，同理Services.GenUserEnPoint(user)必然是返回一个EndPoint，这里GenUserEnPoint(user)返回值是\
	  func(ctx context.Context, request interface{}) (response interface{}, err error)所以是EndPoint类型
	  那么(Services.UserServiceLogMiddleware(logger)是middleware,Services.GenUserEnPoint(user))值作为参数是Endpoint，返回值依然是一个Endpoint，这个返回值作为参数传递给了Services.RateLimit(limit)，Services.RateLimit(limit)也是一个Middware，所以这样写是成立的
	*/

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(Services.MyErrorEncoder),
	}

	serverHandler := httptransport.NewServer(endp, Services.DecodeUserRequest, Services.EncodeUserResponse, options...) //使用go kit创建server传入我们之前定义的两个解析函数

	r := mymux.NewRouter()
	//r.Handle(`/user/{uid:\d+}`, serverHandler) //这种写法支持多种请求方式
	r.Methods("GET", "DELETE").Path(`/user/{uid:\d+}`).Handler(serverHandler) //这种写法仅支持Get，限定只能Get请求
	r.Methods("GET").Path("/health").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-type", "application/json")
		writer.Write([]byte(`{"status":"ok"}`))
	})
	errChan := make(chan error)
	go func() {
		utils.RegService()                                                 //调用注册服务程序
		err := http.ListenAndServe(":"+strconv.Itoa(utils.ServicePort), r) //启动http服务
		if err != nil {
			log.Println(err)
			errChan <- err
		}
	}()
	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sigChan)
	}()
	getErr := <-errChan
	utils.UnRegService()
	log.Println(getErr)
}

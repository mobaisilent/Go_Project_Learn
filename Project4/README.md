## 		项目基本信息

> 项目融入了不少其他包：虽然核心都是围绕着go-kit的

### 关于go-kit

>https://github.com/go-kit/kit

优秀的微服务工具包合集，利用它提供的API和规范可以创建健壮，可维护性高的微服务体系。



文档：

https://godoc.org/github.com/go-kit/kit



安装：

```go
go get github.com/go-kit/kit
```



类似框架：

- https://github.com/micro/go-micro
- https://github.com/koding/kite

### 微服务体系的基本需求（非全部

1. HTTPREST、RPC
2. 日志功能
3. 限流
4. API监控
5. 服务注册与发现
6. API网关服务链路追踪
7. 服务熔断

### Go-kit的三层架构

1. Transport

   主要负责与HTTP、gRPC、thrift等相关的逻辑

2. Endpoint
   定义Request和Response格式，并可以使用装饰器包装函数，
   以此来实现各种中间件嵌套。

3. Service
   这里就是我们的业务类、接口等

> 定义顺序是：Service -> Endpoint -> Transport

### 三层架构定义

#### UserService.go

```go
package Services

// 定义接口
type IUserService interface {
	GetName(userid int) string
}

// 定义结构体
type UserService struct{}

// 实现接口
func (userService UserService) GetName(userid int) string {
	if userid == 101 {
		return "jerry"
	}
	return "guest"
}

```

#### UserEndPoint.go

```go
package Services

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

// 封装User请求结构体
type UserRequest struct {
	Uid int `json:"uid"`
}

// 封装User响应结构体
type UserResponse struct {
	Result string `json:"result"`
}

// 生成User端点：通过实现endpoint.Endpoint接口来生成User端点
// 相当于以及获取了request信息 然后根据类型调用函数
func GenUserEndPoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest) //通过类型断言获取请求结构体
		result := userService.GetName(r.Uid)
		return UserResponse{Result: result}, nil
	}
}

```

#### UserTransport.go

```go
package Services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// 例如封装浏览器的请求url->request
func DecodeUserRequest(c context.Context, r *http.Request) (interface{}, error) { //这个函数决定了使用哪个request结构体来请求
	if r.URL.Query().Get("uid") != "" {
		uid, _ := strconv.Atoi(r.URL.Query().Get("uid"))
		return UserRequest{Uid: uid}, nil
	}
	return nil, errors.New("参数错误")
}

// 封装服务器根据request生成的结果response：使得结果展现为json格式在屏幕中
func EncodeUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json") //设置响应格式为json，这样客户端接收到的值就是json，就是把我们设置的UserResponse给json化了
	return json.NewEncoder(w).Encode(response)         //判断响应格式是否正确
}

```

> 一般是三层架构三个文件夹，我这里演示方便把三层的东西都放到了一个中
>
> 多理一下这三者间的关系就没问题

## HTTP服务

使用kit跑起HTTP服务

### 不带路由信息

#### main.go

```go
package main

import (
	"net/http"    // Import the missing package
	"p4/Services" //引入我们的服务包  这里不使用 . 导入容易发生命名冲突

	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	user := Services.UserService{}
	endp := Services.GenUserEndPoint(user)

	serverHandler := httptransport.NewServer(endp, Services.DecodeUserRequest, Services.EncodeUserResponse)
	http.ListenAndServe(":12345", serverHandler) // ListenAndServe是go内置的包
}

```

http://127.0.0.1:8080/

返回结果：

```
参数错误
```

http://127.0.0.1:8080/?uid=101

返回结果：

![image-20240810222302804](./images/image-20240810222302804.png)

http://127.0.0.1:8080/?uid=10190

返回结果：

![image-20240810222911002](./images/image-20240810222911002.png)

### 带路由信息

![img](./images/20191221111117.png)

#### main.go

```go
package main

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"p4/Services"
)

func main() {
	user := Services.UserService{}
	endp := Services.GenUserEndPoint(user)

	serverHandler := httptransport.NewServer(endp, Services.DecodeUserRequest, Services.EncodeUserResponse) //使用go kit创建server传入我们之前定义的两个解析函数

	r := mux.NewRouter() //使用mux来使服务支持路由
	//r.Handle(`/user/{uid:\d+}`, serverHandler) //这种写法支持多种请求方法，访问Examp: http://localhost:8080/user/121便可以访问
	r.Methods("GET").Path(`/user/{uid:\d+}`).Handler(serverHandler) //这种写法限定了请求只支持GET方法
	http.ListenAndServe(":8080", r)

}
```

修改`UserTransport`的`DecodeUserRequest`函数

```go
func DecodeUserRequest(c context.Context, r *http.Request) (interface{}, error) { //这个函数决定了使用哪个request来请求
	vars := mux.Vars(r)             //通过这个返回一个map，map中存放的是参数key和值，因为我们路由地址是这样的/user/{uid:\d+}，索引参数是uid,访问Examp: http://localhost:8080/user/121，所以值为121
	if uid, ok := vars["uid"]; ok { //
		uid, _ := strconv.Atoi(uid)
		return UserRequest{Uid: uid}, nil
	}
	return nil, errors.New("参数错误")
}
```

### HttpMethods实现不同方法

#### UserService.go

```go
package Services

import "errors"

// 定义接口
type IUserService interface {
	GetName(userid int) string
	DelUser(userid int) error // 注意这里的修改：改半天才改出来
}

// 定义结构体
type UserService struct{}

// 实现接口
func (userService UserService) GetName(userid int) string {
	if userid == 101 {
		return "jerry"
	}
	return "guest"
}

func (userService UserService) DelUser(userid int) error {
	if userid == 101 {
		return errors.New("权限不够")
	}
	return nil
}

```

#### UserEndPoint.go

```go
package Services

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
)

type UserRequest struct { //封装User请求结构体
	Uid    int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
}

func GenUserEndPoint(userService IUserService) endpoint.Endpoint { //当EncodeUserRequest和DecodeUserResponse都不报错的时候才会走这个函数
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest) //通过类型断言获取请求结构体
		result := "nothings"
		if r.Method == "GET" {
			result = userService.GetName(r.Uid)
		} else if r.Method == "DELETE" {
			err := userService.DelUser(r.Uid)
			if err != nil {
				result = err.Error()
			} else {
				result = fmt.Sprintf("userid为%d的用户已删除", r.Uid)
			}
		}
		return UserResponse{Result: result}, nil
	}
}

```

#### UserTransport.go

```go
package Services

import (
	"context"
	"encoding/json"
	"errors"
	mymux "github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func DecodeUserRequest(c context.Context, r *http.Request) (interface{}, error) { //这个函数决定了使用哪个request来请求
	vars := mymux.Vars(r)
	if uid, ok := vars["uid"]; ok {
		uid, _ := strconv.Atoi(uid)
		return UserRequest{Uid: uid, Method: r.Method}, nil //组装请求参数和方法
	}
	return nil, errors.New("参数错误") //如果发生错误返回给客户端这个错误，如果没有返回endpoint的执行结果
}

func EncodeUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

```

#### main.go

```go
package main

import (
	"net/http"
	"p4/Services"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func main() {
	user := Services.UserService{}
	endp := Services.GenUserEndPoint(user)

	serverHandler := httptransport.NewServer(endp, Services.DecodeUserRequest, Services.EncodeUserResponse) //使用go kit创建server传入我们之前定义的两个解析函数

	r := mux.NewRouter() //使用mux来使服务支持路由
	//r.Handle(`/user/{uid:\d+}`, serverHandler) //这种写法支持多种请求方法，访问Examp: http://localhost:8080/user/121便可以访问
	r.Methods("GET","DELETE").Path(`/user/{uid:\d+}`).Handler(serverHandler) //这种写法限定了请求只支持GET方法
	http.ListenAndServe(":8080", r)

}

```

![image-20240811154835496](./images/image-20240811154835496.png)![image-20240811154844911](./images/image-20240811154844911.png)

![image-20240811154847219](./images/image-20240811154847219.png)

## 服务注册

![img](./images/20191221120207.png)

### 关于consul

#### 环境变量安装

下载exe文件然后配置好环境变量就行

![test1](./images/20191221120053.png)

#### docker安装

![img](./images/20191221120400-1723364136813-7.png)

![img](./images/20191221120400.png)

#### 服务查看与注册

![img](./images/20191221122849.png)

![img](./images/20191221122917.png)

对应main.go

```go
package main

import (
	"net/http"
	"p4/Services"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func main() {
	user := Services.UserService{}
	endp := Services.GenUserEndPoint(user)

	serverHandler := httptransport.NewServer(endp, Services.DecodeUserRequest, Services.EncodeUserResponse) //使用go kit创建server传入我们之前定义的两个解析函数

	r := mux.NewRouter() //使用mux来使服务支持路由
	{
		//r.Handle(`/user/{uid:\d+}`, serverHandler) //这种写法支持多种请求方法，访问Examp: http://localhost:8080/user/121便可以访问
		r.Methods("GET", "DELETE").Path(`/user/{uid:\d+}`).Handler(serverHandler) //这种写法限定了请求只支持GET方法
		r.Methods("GET").Path("/health").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-type", "application/json")
			writer.Write([]byte(`{"status":"ok"}`)) // 字符串转字节切片
		})
	}
	http.ListenAndServe(":8080", r)

}

```

对应json文件：

```json
{
    "ID":"userservice",
    "Name":"userservice",
    "Tags""{
    	"primary"
	},
	"Address":"127.0.0.1",
    "Port":8080,
	"Check":{
        "HTTP":"127.0.0.1:8080/health",
        "Interval":"5s"
    }
}
```

提交服务：(linux)

```bash
curl\
  --request PUT\
  --data @p.json\
  localhost:8500/v1/agent/service/register
```

提交服务：（win）

```bash
consul agent -server -bootstrap -ui -client 127.0.0.1 -bind 127.0.0.1 -advertise 127.0.0.1 -data-dir ./consul_data
// 注意内网IP：或者直接就本机IP了
// ./consul_data 日志存储文件
```

或者直接用PostMan提交

反注册：
```bash
curl -Method PUT -Uri http://127.0.0.1:8500/v1/agent/service/deregister/userservice
```

![img](./images/20191222104743.png)

#### 使用consul

```bash
consul agent -server -bootstrap -ui -client 127.0.0.1 -bind 127.0.0.1 -advertise 127.0.0.1 -data-dir ./consul_data
```

![image-20240811162859790](./images/image-20240811162859790.png)

​	注册结果如下：
```bash
statusCode        : 200
StatusDescription : OK
Content           : {}
RawContent        : HTTP/1.1 200 OK
                    Vary: Accept-Encoding
                    X-Consul-Default-Acl-Policy: allow
                    Content-Length: 0
                    Date: Sun, 11 Aug 2024 08:35:55 GMT


Headers           : {[Vary, Accept-Encoding], [X-Consul-Default-Acl-Policy, allow], [Content-Length, 0], [Date, Sun, 11 Aug 
                     2024 08:35:55 GMT]}
RawContentLength  : 0
```

![image-20240811163807168](./images/image-20240811163807168.png)

这就注册好了服务监听。

### 使用go注册consul服务

需要先将consul服务启动，这里不要搞错了

#### consul.go

```go
package utils

import (
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

func RegService() {
	config := consulapi.DefaultConfig() //创建consul配置
	config.Address = "127.0.0.1:8500"
    client, err := consulapi.NewClient(config) //创建客户端
	if err != nil {
		log.Fatal(err)
	}
    
	reg := consulapi.AgentServiceRegistration{}
	reg.Name = "userservice"      //注册service的名字
	reg.Address = "127.0.0.1" //注册service的ip
	reg.Port = 8080               //注册service的端口
	reg.Tags = []string{"primary"}

	check := consulapi.AgentServiceCheck{}          //创建consul的检查器
	check.Interval = "5s"                           //设置consul心跳检查时间间隔
	check.HTTP = "http://127.0.0.1:8080/health" //设置检查使用的url
	reg.Check = &check // 绑定检查器

	
	err = client.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatal(err)
	}
}

```

#### main.go

```go
package utils

import (
	"log"

	consulapi "github.com/hashicorp/consul/api"
)

func RegService() {
	config := consulapi.DefaultConfig() //创建consul配置
	config.Address = "127.0.0.1:8500"
	reg := consulapi.AgentServiceRegistration{}
	reg.Name = "userservice"  //注册service的名字
	reg.Address = "127.0.0.1" //注册service的ip
	reg.Port = 8080           //注册service的端口
	reg.Tags = []string{"primary"}

	check := consulapi.AgentServiceCheck{}      //创建consul的检查器
	check.Interval = "5s"                       //设置consul心跳检查时间间隔
	check.HTTP = "http://127.0.0.1:8080/health" //设置检查使用的url

	reg.Check = &check // 绑定检查器

	client, err := consulapi.NewClient(config) //创建客户端
	if err != nil {
		log.Fatal(err)
	}
	err = client.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatal(err)
	}
}

```

运行main.go就可以看见服务已经添加好监听了

![image-20240811172339194](./images/image-20240811172339194.png)

### Go退出时反注册/优雅关闭

先看一下直接关闭ctrl+c exit的结果如下：

![image-20240811173306549](./images/image-20240811173340501.png)

也就是相当于服务失效。而不是反注册。



修改consul.go

#### consul.go

```go
package utils

import (
	consulapi "github.com/hashicorp/consul/api"
	"log"
)

// 创建反注册Client
var ConsulClient *consulapi.Client

func init() {
	config := consulapi.DefaultConfig() //创建consul配置
	config.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(config) //创建客户端
	if err != nil {
		log.Fatal(err)
	}
	ConsulClient = client
}

func RegService() {
	reg := consulapi.AgentServiceRegistration{}
	reg.Name = "userservice"  //注册service的名字
	reg.Address = "127.0.0.1" //注册service的ip
	reg.Port = 8080           //注册service的端口
	reg.Tags = []string{"primary"}

	check := consulapi.AgentServiceCheck{}      //创建consul的检查器
	check.Interval = "5s"                       //设置consul心跳检查时间间隔
	check.HTTP = "http://127.0.0.1:8080/health" //设置检查使用的url
	reg.Check = &check

	err := ConsulClient.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatal(err)
	}

}

// 多了这个反注册函数
func UnRegService() {
	ConsulClient.Agent().ServiceDeregister("userservice") // 添加反注册监听
}

```

修改main.go

#### main.go

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"p4/Services"
	"p4/utils"
	"syscall"

	httptransport "github.com/go-kit/kit/transport/http"
	mymux "github.com/gorilla/mux"
)

func main() {
	user := Services.UserService{}
	endp := Services.GenUserEndPoint(user)

	serverHandler := httptransport.NewServer(endp, Services.DecodeUserRequest, Services.EncodeUserResponse) //使用go kit创建server传入我们之前定义的两个解析函数

	r := mymux.NewRouter()
	{
		//r.Handle(`/user/{uid:\d+}`, serverHandler) //这种写法支持多种请求方式
		r.Methods("GET", "DELETE").Path(`/user/{uid:\d+}`).Handler(serverHandler)
		r.Methods("GET").Path("/health").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-type", "application/json")
			writer.Write([]byte(`{"status":"ok"}`))
		}) //这种写法仅支持Get限定只能Get,DELETE请求
	}

	errChan := make(chan error)
	go func() {
		utils.RegService() //调用注册服务程序
		err := http.ListenAndServe(":8080", r)
		if err != nil {
			log.Println(err)
			errChan <- err // 这里是注册失败的error
		}
		// http.ListenAndServe函数返回一个错误，那么这个错误会被发送到errChan通道
	}()
	go func() {
		sigChan := make(chan os.Signal, 1) //1是信号通道的缓冲大小。这意味着信号通道可以存储一个信号，即使没有Goroutine立即接收它，防止丢失信号
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sigChan) // 这里是检测结束的error信息
	}()
	getErr := <-errChan //只要报错 否则service关闭阻塞在这里的会进行下去
	utils.UnRegService()
	log.Println(getErr)
}

```

#### 测试

![image-20240811174333987](./images/image-20240811174333987.png)

![image-20240811174343151](./images/image-20240811174343151.png)

可见服务关闭了。



## 服务发现

![img](./images/20191222180841.png)

> 客户端直接调用api

![img](./images/20191222183616.png)



### 开个客户端项目（新项目

![image-20240813220607101](./images/image-20240813220607101.png)

#### UserEndpoint.go

```go
package Services

// 定义信息传送结构体
type UserRequest struct {
	Uid    int    `json:"uid"`
	Method string `json:"method"`
}

type UserResponse struct {
	Result string `json:"result"`
}

```

#### UserTransport.go

```go
package Services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// 处理http请求和响应的两个函数：：也就相当于客户端直接向浏览器一样封装相关信息直接调用API了
func GetUserInfo_Request(_ context.Context, request *http.Request, r interface{}) error {
	user_request := r.(UserRequest)
	request.URL.Path += "/user/" + strconv.Itoa(user_request.Uid)
	return nil
}

func GetUserInfo_Response(_ context.Context, res *http.Response) (response interface{}, err error) {
	if res.StatusCode > 400 {
		return nil, errors.New("no data")
	}
	var user_response UserResponse
	err = json.NewDecoder(res.Body).Decode(&user_response)
	if err != nil {
		return nil, err
	}
	return user_response, err
}

```

#### main.go

```go
package main

import (
	"context"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"gomicro2/Services"
	"net/url"
	"os"
)

func main() {
	tgt, _ := url.Parse("http://127.0.0.1:8080")
	//创建一个直连client，这里我们必须写两个func,一个是如何请求,一个是响应我们怎么处理
	client := httptransport.NewClient("GET", tgt, Services.GetUserInfo_Request, Services.GetUserInfo_Response)
	getUserInfo := client.Endpoint() //通过这个拿到了定义在服务端的endpoint也就是上面这段代码return出来的函数，直接在本地就可以调用服务端的代码

	ctx := context.Background() //创建一个上下文

	//执行：：传入ctx和结构体：已知是GET方法和路径和结构体，对应使用服务端的函数
	res, err := getUserInfo(ctx, Services.UserRequest{Uid: 101}) //使用go-kit插件来直接调用服务
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	userinfo := res.(Services.UserResponse)
	fmt.Println(userinfo.Result)

}

```

结果如下：
![image-20240813221415342](./images/image-20240813221415342.png)

### consul获取服务实例，调用测试

> 在上个项目2的基础上直接修改product.go就行，服务端代码不需要修改

![image-20240814145234807](./images/image-20240814145234807.png)

![image-20240814145246283](./images/image-20240814145246283.png)

![image-20240814145306646](./images/image-20240814145306646.png)

#### product.go

```go
package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/hashicorp/consul/api"
	"gomicro2/Services"
	"io"
	"net/url"
	"os"
)

func main() {
	//第一步创建client
	{
		config := api.DefaultConfig()
		config.Address = "127.0.0.1:8500"  // 配置信息
		api_client, _ := api.NewClient(config)
		client := consul.NewClient(api_client)  // 链式创建客户端

		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stdout)
			var Tag = []string{"primary"}
            
			//第二部创建一个consul的实例
			instancer := consul.NewInstancer(client, logger, "userservice", Tag, true) //最后的true表示只有通过健康检查的服务才能被得到
			{
				factory := func(service_url string) (endpoint.Endpoint, io.Closer, error) { //factory定义了如何获得服务端的endpoint,这里的service_url是从consul中读取到的service的address我这里是127.0.0.1:8000
					tart, _ := url.Parse("http://" + service_url)                                                                                 //server ip +8080真实服务的地址
					return httptransport.NewClient("GET", tart, Services.GetUserInfo_Request, Services.GetUserInfo_Response).Endpoint(), nil, nil //我再GetUserInfo_Request里面定义了访问哪一个api把url拼接成了http://127.0.0.1:8000/v1/user/{uid}的形式
				}
				endpointer := sd.NewEndpointer(instancer, factory, logger)
				endpoints, _ := endpointer.Endpoints()
				fmt.Println("服务有", len(endpoints), "条")
				getUserInfo := endpoints[0] //写死获取第一个
				ctx := context.Background() //第三步：创建一个context上下文对象

				//第四步：执行
				res, err := getUserInfo(ctx, Services.UserRequest{Uid: 101})
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				//第五步：断言，得到响应值
				userinfo := res.(Services.UserResponse)
				fmt.Println(userinfo.Result)
			}
		}
	}

}
```

结果如下：
![image-20240814150327620](./images/image-20240814150327620.png)



### 根据命令行参数注册多个服务

![image-20240814165144180](./images/image-20240814165144180.png)

![image-20240814165305991](./images/image-20240814165305991.png)

​	修改服务端p4里main的代码和consul代码

#### consul.go

```go
package utils

import (
	"fmt"
	"github.com/google/uuid"
	consulapi "github.com/hashicorp/consul/api"
	"log"
	"strconv"
)

var ConsulClient *consulapi.Client
var ServiceID string
var ServiceName string
var ServicePort int

func init() {
	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	client, err := consulapi.NewClient(config) //创建客户端
	if err != nil {
		log.Fatal(err)
	}
	ConsulClient = client
	ServiceID = "userservice" + uuid.New().String() //因为最终这段代码是在不同的机器上跑的，是分布式的，有好几台机器提供相同的server，所以这里存到consul中的id必须是唯一的，否则只有一台服务器可以注册进去，这里使用uuid保证唯一性
}

func SetServiceNameAndPort(name string, port int) {
	ServiceName = name
	ServicePort = port
}

func RegService() {
	reg := consulapi.AgentServiceRegistration{}
	reg.ID = ServiceID        //设置不同的Id，即使是相同的service name也得有不同的id
	reg.Name = ServiceName    //注册service的名字
	reg.Address = "127.0.0.1" //注册service的ip
	reg.Port = ServicePort    //注册service的端口
	reg.Tags = []string{"primary"}

	check := consulapi.AgentServiceCheck{}                                   //创建consul的检查器
	check.Interval = "5s"                                                    //设置consul心跳检查时间间隔
	check.HTTP = "http://127.0.0.1:" + strconv.Itoa(ServicePort) + "/health" //设置检查使用的url
	reg.Check = &check

	err := ConsulClient.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatal(err)
	}
}

func UnRegService() {
	ConsulClient.Agent().ServiceDeregister("userservice")
}

```
> 重新启动main.go注册服务（consul是由consul.go单独启动的consul服务）

#### main.go

```go
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"p4/Services"
	"p4/utils"
	"strconv"
	"syscall"

	httptransport "github.com/go-kit/kit/transport/http"
	mymux "github.com/gorilla/mux"
)

func main() {

	name := flag.String("name", "", "服务名称")
	port := flag.Int("port", 0, "服务端口")
	flag.Parse()

	if *name == "" {
		log.Fatal("请指定服务名")
	}

	if *port == 0 {
		log.Fatal("请指定端口")
	}

	utils.SetServiceNameAndPort(*name, *port) //设置服务名和端口

	user := Services.UserService{}
	endp := Services.GenUserEndPoint(user)

	serverHandler := httptransport.NewServer(endp, Services.DecodeUserRequest, Services.EncodeUserResponse) //使用go kit创建server传入我们之前定义的两个解析函数

	r := mymux.NewRouter()
	{
		//r.Handle(`/user/{uid:\d+}`, serverHandler) //这种写法支持多种请求方式
		r.Methods("GET", "DELETE").Path(`/user/{uid:\d+}`).Handler(serverHandler)
		r.Methods("GET").Path("/health").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-type", "application/json")
			writer.Write([]byte(`{"status":"ok"}`))
		}) //这种写法仅支持Get限定只能Get,DELETE请求
	}

	errChan := make(chan error)
	go func() {
		utils.RegService()                                     //调用注册服务程序
		err := http.ListenAndServe(":"+strconv.Itoa(*port), r) //启动http服务
		if err != nil {
			log.Println(err)
			errChan <- err // 这里是注册失败的error
		}
		// http.ListenAndServe函数返回一个错误，那么这个错误会被发送到errChan通道
	}()
	go func() {
		sigChan := make(chan os.Signal, 1) //1是信号通道的缓冲大小。这意味着信号通道可以存储一个信号，即使没有Goroutine立即接收它，防止丢失信号
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sigChan) // 这里是检测结束的error信息
	}()
	getErr := <-errChan //只要报错 否则service关闭阻塞在这里的会进行下去
	utils.UnRegService()
	log.Println(getErr)
}

```

> 也就是通过启动的时候设置好参数进行使用了
>
> go run main.go -name userservice -port 8080	
>
> go run main.go -name userservice -port 8081	

结果如下：
![image-20240814172851354](./images/image-20240814172851354.png)

### 使用多个服务

> 修改p4既服务端UserEndPoint.go获取端口信息

![image-20240814222233633](./images/image-20240814222233633.png)

```go
package Services

import (
	"context"
	"fmt"
	"p4/utils"
	"strconv"

	"github.com/go-kit/kit/endpoint"
)

type UserRequest struct { //封装User请求结构体
	Uid    int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
}

func GenUserEndPoint(userService IUserService) endpoint.Endpoint { //当EncodeUserRequest和DecodeUserResponse都不报错的时候才会走这个函数
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest) //通过类型断言获取请求结构体
		result := "nothings"
		if r.Method == "GET" {
			result = userService.GetName(r.Uid) + strconv.Itoa(utils.ServicePort)
		} else if r.Method == "DELETE" {
			err := userService.DelUser(r.Uid)
			if err != nil {
				result = err.Error()
			} else {
				result = fmt.Sprintf("userid为%d的用户已删除", r.Uid)
			}
		}
		return UserResponse{Result: result}, nil
	}
}

```

#### product.go

循环测试：修改客户端product.go代码

> 就是增加了个死循环和Sleep

```go
package main

import (
	"context"
	"fmt"
	"gomicro2/Services"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/hashicorp/consul/api"
)

func main() {
	//第一步创建client
	{
		config := api.DefaultConfig()
		config.Address = "localhost:8500"
		api_client, _ := api.NewClient(config)
		client := consul.NewClient(api_client)

		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stdout)
			var Tag = []string{"primary"}
			//第二部创建一个consul的实例
			instancer := consul.NewInstancer(client, logger, "userservice", Tag, true) //最后的true表示只有通过健康检查的服务才能被得到
			{
				factory := func(service_url string) (endpoint.Endpoint, io.Closer, error) { //factory定义了如何获得服务端的endpoint,这里的service_url是从consul中读取到的service的address我这里是192.168.3.14:8000
					tart, _ := url.Parse("http://" + service_url)                                                                                 //server ip +8080真实服务的地址
					return httptransport.NewClient("GET", tart, Services.GetUserInfo_Request, Services.GetUserInfo_Response).Endpoint(), nil, nil //我再GetUserInfo_Request里面定义了访问哪一个api把url拼接成了http://192.168.3.14:8000/v1/user/{uid}的形式
				}
				endpointer := sd.NewEndpointer(instancer, factory, logger)
				endpoints, _ := endpointer.Endpoints()
				fmt.Println("服务有", len(endpoints), "条")

				for {
					getUserInfo := endpoints[0] //写死获取第一个
					ctx := context.Background() //第三步：创建一个context上下文对象

					//第四步：执行
					res, err := getUserInfo(ctx, Services.UserRequest{Uid: 101})
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					//第五步：断言，得到响应值
					userinfo := res.(Services.UserResponse)
					fmt.Println(userinfo.Result)
					time.Sleep(3 * time.Second)
				}
			}
		}
	}

}

```

结果如下：

![image-20240814222938347](./images/image-20240814222938347.png)

> 因为是写死的下标第0条没用负载均衡的效果

## 负载均衡

三个服务如下；
![image-20240814223814809](./images/image-20240814223814809.png)

### 轮询

![image-20240814223140079](./images/image-20240814223140079.png)

> 这个函数的原理是“简单的轮询算法”

将product的部分代码改成这样：
```go
        fmt.Println("服务有", len(endpoints), "条")

        mylb := lb.NewRoundRobin(endpointer)
        for {
            //getUserInfo := endpoints[0] //写死获取第一个
            getUserInfo, _ := mylb.Endpoint() //通过lb获取
            ctx := context.Background()       //第三步：创建一个context上下文对象

            //第四步：执行
            res, err := getUserInfo(ctx, Services.UserRequest{Uid: 101})
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
            //第五步：断言，得到响应值
            userinfo := res.(Services.UserResponse)
            fmt.Println(userinfo.Result)
            time.Sleep(3 * time.Second)
        }
```

结果如下：

![image-20240814223704329](./images/image-20240814223704329.png)

### 随机

修改一行代码即可

```go
		mylb := lb.NewRandom(endpointer, time.Now().UnixNano())
```

结果如下：
![image-20240814224419290](./images/image-20240814224419290.png)

## API限流

### 桶算法

![image-20240814225024733](./images/image-20240814225024733.png)

![image-20240814225059493](./images/image-20240814225059493.png)

### Wait

#### test.go

```go
package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

func main() {

	r := rate.NewLimiter(1, 5)  // 桶的初始容量为5  每次消耗一个
	ctx := context.Background()
	for {
		err := r.Wait(ctx)  // err := r.WaitN(ctx,2)   //这里就是写死每次消耗两个	
		if err != nil {
			fmt.Println("error")
			return
		}
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	}
}

```

>go build test.go
>
>./test.exe

### Allow

```go
package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	r := rate.NewLimiter(1, 5) //1表示每次放进筒内的数量，桶内的令牌数是5，最大令牌数也是5，这个筒子是自动补充的，你只要取了令牌不管你取多少个，这里都会在每次取完后自动加1个进来，因为我们设置的是1

	for {

		if r.AllowN(time.Now(), 2) { //AllowN表示取当前的时间，这里是一次取2个，如果当前不够取两个了，本次就不取，再放一个进去，然后返回false
			fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
		} else {
			fmt.Println("too many request")
		}
		time.Sleep(time.Second)
	}

}
```

![image-20240814231435137](./images/image-20240814231435137.png)

##### 实践

**test.go**

```go
package main

import (
	"golang.org/x/time/rate"
	"net/http"
)

var r = rate.NewLimiter(1, 5)

func MyLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !r.Allow() {
			http.Error(writer, "too many requests", http.StatusTooManyRequests)
		} else {
			next.ServeHTTP(writer, request)
		}	
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("OK!!!"))
	})
	http.ListenAndServe(":8080", MyLimit(mux))
}

```

> 封装到url的api中服务中

![image-20240814232624636](./images/image-20240814232624636.png)

### 限流应用到项目中

在p4文件中

#### UserEndpoint.go

```go
package Services

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
	"p4/utils"
	"strconv"
)

type UserRequest struct { //封装User请求结构体
	Uid    int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
}

// 加入限流功能中间件
func RateLimit(limit *rate.Limiter) endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
	return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, errors.New("too many request")
			}
			return next(ctx, request)
		}
	}
}

func GenUserEndPoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest) //通过类型断言获取请求结构体
		result := "nothings"
		if r.Method == "GET" {
			result = userService.GetName(r.Uid) + strconv.Itoa(utils.ServicePort)
			fmt.Println(result)
		} else if r.Method == "DELETE" {
			err := userService.DelUser(r.Uid)
			if err != nil {
				result = err.Error()
			} else {
				result = fmt.Sprintf("userid为%d的用户已删除", r.Uid)
			}
		}
		return UserResponse{Result: result}, nil
	}
}

```

> 添加了一个限流中间件的使用

#### main.go

```go
package main

import (
	"flag"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	mymux "github.com/gorilla/mux"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"os/signal"
	"p4/Services"
	"p4/utils"
	"strconv"
	"syscall"
)

func main() {
	name := flag.String("name", "", "服务名称")
	port := flag.Int("port", 0, "服务端口")
	flag.Parse()
	if *name == "" {
		log.Fatal("请指定服务名")
	}
	if *port == 0 {
		log.Fatal("请指定端口")
	}
	utils.SetServiceNameAndPort(*name, *port) //设置服务名和端口

	user := Services.UserService{}
	limit := rate.NewLimiter(1, 5)
	endp := Services.RateLimit(limit)(Services.GenUserEndPoint(user)) //调用限流代码生成的中间件

	serverHandler := httptransport.NewServer(endp, Services.DecodeUserRequest, Services.EncodeUserResponse) //使用go kit创建server传入我们之前定义的两个解析函数

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
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-sigChan)
	}()
	getErr := <-errChan
	utils.UnRegService()
	log.Println(getErr)
}

```

## 重写自定义Error

### 自定义ErrorEncoder

```go
func MyErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
    contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
    w.Header().Set("Content-type", contentType) //设置请求头
    w.WriteHeader(429) //写入返回码
    w.Write(body)
}
```

### 生成ServerOption继而生成相应的Handler

```go
options := []httptransport.ServerOption{
        httptransport.ServerErrorEncoder(Services.MyErrorEncoder),
        //ServerErrorEncoder支持ErrorEncoder类型的参数 type ErrorEncoder func(ctx context.Context, err error, w http.ResponseWriter)
           //我们自定义的MyErrorEncoder只要符合ErrorEncoder类型就可以传入
} //创建ServerOption切片

serverHandler := httptransport.NewServer(endp, Services.DecodeUserRequest, Services.EncodeUserResponse, options...)//在创建handler的同事把切片展开传入

```

> p4项目中

#### main.go

```go
package main

import (
	"flag"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	mymux "github.com/gorilla/mux"
	"golang.org/x/time/rate"
	"gomicro/Services"
	"gomicro/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	name := flag.String("name", "", "服务名称")
	port := flag.Int("port", 0, "服务端口")
	flag.Parse()
	if *name == "" {
		log.Fatal("请指定服务名")
	}
	if *port == 0 {
		log.Fatal("请指定端口")
	}
	utils.SetServiceNameAndPort(*name, *port) //设置服务名和端口

	user := Services.UserService{}
	limit := rate.NewLimiter(1, 5)
	endp := Services.RateLimit(limit)(Services.GenUserEnPoint(user))

	options := []httptransport.ServerOption{ //生成ServerOtion切片，传入我们自定义的错误处理函数
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

//因为设置了限流，访问过快就会报错，此时的报错Code就是我们自定义的429

```

#### UserEndpoint.go

```go
package Services

import (
	"context"
	"errors"
	"fmt"
	"p4/utils"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/time/rate"
)

type UserRequest struct { //封装User请求结构体
	Uid    int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
}

// 加入限流功能中间件
func RateLimit(limit *rate.Limiter) endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
	return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, errors.New("too many requests")
			}
			return next(ctx, request)
		}
	}
}

func GenUserEndPoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest) //通过类型断言获取请求结构体
		result := "nothings"
		if r.Method == "GET" {
			result = userService.GetName(r.Uid) + strconv.Itoa(utils.ServicePort)
			fmt.Println(result)
		} else if r.Method == "DELETE" {
			err := userService.DelUser(r.Uid)
			if err != nil {
				result = err.Error()
			} else {
				result = fmt.Sprintf("userid为%d的用户已删除", r.Uid)
			}
		}
		return UserResponse{Result: result}, nil
	}
}

```

### 自定义Error

#### 自定义错误结构体

```go
package utils

type MyError struct {
    Code    int
    Message string
}


func NewMyError(code int, msg string) error {
    return &MyError{Code: code, Message: msg}
}

func (this *MyError) Error() string {
    return this.Message
}

```

#### 如何处理我们自定义的Error

```
func MyErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
    contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
    w.Header().Set("Content-type", contentType) //设置请求头
    if myerr, ok := err.(*utils.MyError); ok { //通过类型断言判断当前error的类型，走相应的处理
        w.WriteHeader(myerr.Code)
        w.Write(body)
    } else {
        w.WriteHeader(500)
        w.Write(body)
    }

}
```

#### 调用自定义结构体

```go
func RateLimit(limit *rate.Limiter) endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
    return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
        return func(ctx context.Context, request interface{}) (response interface{}, err error) {
            if !limit.Allow() {
                return nil, utils.NewMyError(429, "toot many request") //使用我们自定的错误结构体
            }
            return next(ctx, request)
        }
    }
}

```

## 熔断器

> 核心使用的包：hystrix
>
> 官网：https://github.com/Netflix/Hystrix
>
> 获取：go get github.com/afex/hystrix-go

本质：防止级联故障

![image-20240901143359825](/Users/mobai/Desktop/Golang/Project/Project4/images/image-20240901143359825.png)

### 简单例子

模拟延时的例子：

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Product struct {
	ID    int
	Title string
	Price int
}

// 创建商品结构体

func getProduct() (Product, error) {
	r := rand.Intn(10)
	if r < 6 { //模拟api卡顿和超时效果  // 模拟3秒
		time.Sleep(time.Second * 3)
	}
	return Product{
		ID:    101,
		Title: "Golang从入门到精通",
		Price: 12,
	}, nil
}

// 创建获取商品信息的函数 根据随机获取的数值来模拟api卡顿和超时效果

func main() {
	rand.Seed(time.Now().UnixNano()) //设置种子
	for {
		p, _ := getProduct()
		fmt.Println(p)
		time.Sleep(time.Second)
	}
}

```

### 调用hystrix

> timeout为2s那么延时3s就会报错

```go
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

type Product struct {
	ID    int
	Title string
	Price int
}

func getProduct() (Product, error) {
	r := rand.Intn(10)
	if r < 6 { //模拟api卡顿和超时效果
		time.Sleep(time.Second * 4)
	}
	return Product{
		ID:    101,
		Title: "Golang从入门到精通",
		Price: 12,
	}, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	configA := hystrix.CommandConfig{ //创建一个hystrix的config
		Timeout: 3000, //command运行超过3秒就会报超时错误
	}
	// configA2 := hystrix.CommandConfig{ //创建一个hystrix的config
	// 	Timeout: 4000, //定义第二个config 对于4秒超时生效
	// }

	// 也就是实际的业务需要修改报错提醒的时候只需要修改使用的config参数就行，不需要修改实现的代码

	hystrix.ConfigureCommand("get_prod", configA) //hystrix绑定command
	for {
		err := hystrix.Do("get_prod", func() error { //使用hystrix来讲我们的操作封装成command
			p, _ := getProduct() //这里会随机延迟0-4秒
			fmt.Println(p)
			return nil
		}, nil)
		if err != nil {
			fmt.Println(err)
		}
	}
}

```

下面是三秒的超时报错：
![image-20240901150230380](images/image-20240901150230380.png)

### 降级

> 能够一定程度的解决降级问题

![image-20240901150424650](images/image-20240901150424650.png)

```go
package main

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"math/rand"
	"time"
)

type Product struct {
	ID    int
	Title string
	Price int
}

func getProduct() (Product, error) {
	r := rand.Intn(10)
	if r < 6 { //模拟api卡顿和超时效果
		time.Sleep(time.Second * 4)
	}
	return Product{
		ID:    101,
		Title: "Golang从入门到精通",
		Price: 12,
	}, nil
}

// 返回普通商品的函数

func RecProduct() (Product, error) {
	return Product{
		ID:    999,
		Title: "推荐商品",
		Price: 120,
	}, nil
}

// 返回推荐商品的函数

func main() {
	rand.Seed(time.Now().UnixNano())
	configA := hystrix.CommandConfig{
		Timeout: 3000,
	}
	hystrix.ConfigureCommand("get_prod", configA) //绑定command

	for {
		err := hystrix.Do("get_prod", func() error { //使用hystrix来讲我们的操作封装成command
			p, _ := getProduct() //这里会随机延迟0-4秒
			fmt.Println(p)
			return nil
		}, func(e error) error {
			fmt.Println(RecProduct()) //超时后调用回调函数返回推荐商品
			return errors.New("my timeout")
		})
		
		if err != nil {
			//如果降级也失败了，在这里定义业务逻辑
		}
	}
}

```

结果如下：

![image-20240901150908764](images/image-20240901150908764.png)

### 异步执行

异步执行和服务降级，使用hystrix.Go()函数的返回值是chan err

```go
package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"math/rand"
	"time"
)

type Product struct {
	ID    int
	Title string
	Price int
}

func getProduct() (Product, error) {
	r := rand.Intn(10)
	if r < 6 { // 模拟 API 卡顿和超时效果
		time.Sleep(time.Second * 4)
	}
	return Product{
		ID:    101,
		Title: "Golang从入门到精通",
		Price: 12,
	}, nil
}

// 生成普通商品

func RecProduct() (Product, error) {
	return Product{
		ID:    999,
		Title: "推荐商品",
		Price: 120,
	}, nil
}

// 生成推荐商品

func main() {
	rand.Seed(time.Now().UnixNano())
	configA := hystrix.CommandConfig{
		Timeout: 5000,
	}
	hystrix.ConfigureCommand("get_prod", configA)
	resultChan := make(chan Product, 1)
	defer close(resultChan)
	// 创建管程

	for {
		errs := hystrix.Go("get_prod", func() error {
			p, err := getProduct()
			if err != nil {
				return err
			}
			resultChan <- p
			return nil
		}, func(e error) error {
			rcp, err := RecProduct()
			if err != nil {
				return err
			}
			resultChan <- rcp
			return nil
		})

		select {
		case getProd := <-resultChan:
			fmt.Println(getProd)
		case err := <-errs:
			fmt.Println("Error:", err)
		}
		// 监管管程 即监管了 getProd 也监管了 err
	}
}

```

### 熔断器控制最大并发数

```go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

type Product struct {
	ID    int
	Title string
	Price int
}

func getProduct() (Product, error) {
	r := rand.Intn(10)
	if r < 6 { //模拟api卡顿和超时效果
		//time.Sleep(time.Second * 4)
	}
	return Product{
		ID:    101,
		Title: "Golang从入门到精通",
		Price: 12,
	}, nil
}

func RecProduct() (Product, error) {
	return Product{
		ID:    999,
		Title: "推荐商品",
		Price: 120,
	}, nil

}

func main() {
	rand.Seed(time.Now().UnixNano())
	configA := hystrix.CommandConfig{ //创建一个hystrix的config
		Timeout:               3000, //command运行超过3秒就会报超时错误
		MaxConcurrentRequests: 5,    //控制最大并发数为5，如果超过5会调用我们传入的回调函数降级
	}
	hystrix.ConfigureCommand("get_prod", configA) //hystrix绑定command
	resultChan := make(chan Product, 1)

	wg := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		go (func() {
			wg.Add(1)
			defer wg.Done()

			errs := hystrix.Go("get_prod", func() error { //使用hystrix来讲我们的操作封装成command,hystrix返回值是一个chan error
				p, _ := getProduct() //这里会随机延迟0-4秒
				resultChan <- p
				return nil //这里返回的error在回调中可以获取到，也就是下面的e变量
			}, func(e error) error {
				fmt.Println(e)           // 查看error信息
				rcp, err := RecProduct() //推荐商品,如果这里的err不是nil,那么就会忘errs中写入这个err，下面的select就可以监控到
				resultChan <- rcp
				return err
			})

			select {
			case getProd := <-resultChan:
				fmt.Println(getProd)
			case err := <-errs: //使用hystrix.Go时返回值是chan error各个协程的错误都放到errs中
				fmt.Println(err, 1)
			}
		})()
	}
	wg.Wait()
}

```

结果如下：
![image-20240901154822014](images/image-20240901154822014.png)

### 熔断器三种状态 打开与关闭和半开

![test](images/20191223165830.png)

![image-20240901153150762](images/image-20240901153150762.png)

```go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

type Product struct {
	ID    int
	Title string
	Price int
}

func getProduct() (Product, error) {
	r := rand.Intn(10)
	if r < 6 { //模拟api卡顿和超时效果
		time.Sleep(time.Second * 10)
	}
	return Product{
		ID:    101,
		Title: "Golang从入门到精通",
		Price: 12,
	}, nil
}

func RecProduct() (Product, error) {
	return Product{
		ID:    999,
		Title: "推荐商品",
		Price: 120,
	}, nil

}

func main() {
	rand.Seed(time.Now().UnixNano())

	// 下面是熔断器 config 的详细配置信息
	configA := hystrix.CommandConfig{ //创建一个hystrix的config
		Timeout:                3000,                  //command运行超过3秒就会报超时错误，并且在一个统计窗口内处理的请求数量达到阈值会调用我们传入的降级回调函数
		MaxConcurrentRequests:  5,                     //控制最大并发数为5，并且在一个统计窗口内处理的请求数量达到阈值会调用我们传入的降级回调函数
		RequestVolumeThreshold: 5,                     //判断熔断的最少请求数，默认是5；只有在一个统计窗口内处理的请求数量达到这个阈值，才会进行熔断与否的判断
		ErrorPercentThreshold:  5,                     //判断熔断的阈值，默认值5，表示在一个统计窗口内有50%的请求处理失败，比如有20个请求有10个以上失败了会触发熔断器短路直接熔断服务
		SleepWindow:            int(time.Second * 10), //熔断器短路多久以后开始尝试是否恢复，这里设置的是10
	}
	hystrix.ConfigureCommand("get_prod", configA) //hystrix绑定command
	c, _, _ := hystrix.GetCircuit("get_prod")     //返回值有三个，第一个是熔断器指针,第二个是bool表示是否能够取到，第三个是error

	resultChan := make(chan Product, 1)
	wg := &sync.WaitGroup{}

	for i := 0; i < 20; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()

			errs := hystrix.Do("get_prod", func() error { //使用hystrix来讲我们的操作封装成command,hystrix返回值是一个chan error
				p, _ := getProduct() //这里会随机延迟0-4秒
				fmt.Println("hello")
				resultChan <- p
				return nil //这里返回的error在回调中可以获取到，也就是下面的e变量
			}, func(e error) error {
				fmt.Println("hello")
				rcp, err := RecProduct() //推荐商品,如果这里的err不是nil,那么就会忘errs中写入这个err，下面的select就可以监控到
				resultChan <- rcp
				return err
			})
			// 这里和前面的笔记别无二致

			if errs != nil { //这里errs是error接口，但是使用hystrix.Go异步执行时返回值是chan error各个协程的错误都放到errs中
				fmt.Println(errs)
			} else {
				select {
				case prod := <-resultChan:
					fmt.Println(prod)
				}
			}

			fmt.Println(c.IsOpen())       //查看熔断器是否打开，一旦打开所有的请求都会走fallback
			fmt.Println(c.AllowRequest()) //查看是否允许请求服务
		}()
	}
	wg.Wait()
}

```

下面是输出，可以看到hello执行的次数明显没有200次，所有的输出加起来也只有89行，但是我们开了200个协程，也就是说熔断器直接熔断服务了，一部分请求直接被拒绝了，只有等待我们设置的10s后再把熔断器设置成半打开状态，再次执行根据结果判断熔断器应该设置成什么状态

```bash
hello
hello
hello
hello
hello
hello
hello
hello
hello
hello
{999 推荐商品 120}
true
false
hello
hello
hello
hello
{101 Golang从入门到精通 12}
true
false
{999 推荐商品 120}
true
false
{999 推荐商品 120}
true
false
{999 推荐商品 120}
true
false
{999 推荐商品 120}
true
false
{999 推荐商品 120}
true
false
{999 推荐商品 120}
true
false
hello
hello
hello
{999 推荐商品 120}
true
false
{999 推荐商品 120}
true
false
{999 推荐商品 120}
true
false
{999 推荐商品 120}
true
false
{999 推荐商品 120}
true
false
{101 Golang从入门到精通 12}
true
false
{999 推荐商品 120}
true
false
{999 推荐商品 120}
true
false
{999 推荐商品 120}
true
false
hello
hello
hello
{999 推荐商品 120}
true
false
{999 推荐商品 120}
true
false
{999 推荐商品 120}
true
false

```

> 在true的情况下会直接执行降级方法：那么一下就可以输出下一行

### 熔断器应用

熔断器放哪？
![image-20240901155952930](images/image-20240901155952930.png)

不打开服务中心尝试链接 查看熔断器

#### 封装到服务端 (package.go)

```go
package main

import (
    "fmt"
    "github.com/afex/hystrix-go/hystrix"
    "gomicro2/util"
    "log"
    "time"
)



func main() {
    configA := hystrix.CommandConfig{
        Timeout:                2000,
        MaxConcurrentRequests:  5,
        RequestVolumeThreshold: 3,
        SleepWindow:            int(time.Second * 10),
        ErrorPercentThreshold:  20,
    }

    hystrix.ConfigureCommand("getuser", configA)
    err := hystrix.Do("getuser", func() error {
        res, err := util.GetUser() //调用方法
        fmt.Println(res)
        return err
    }, func(e error) error {
        fmt.Println("降级用户")
        return e
    })
    if err != nil {
        log.Fatal(err)
    }
}

```

#### Util.GetUser

```go
package util

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	consulapi "github.com/hashicorp/consul/api"
	"gomicro2/Services"
	"io"
	"net/url"
	"os"
	"time"
)

func GetUser() (string, error) {
	//第一步创建client
	{
		config := consulapi.DefaultConfig()            //初始化consul的配置
		config.Address = "localhost:8500"              //consul的地址
		api_client, err := consulapi.NewClient(config) //根据consul的配置初始化client
		if err != nil {
			return "", err
		}
		client := consul.NewClient(api_client) //根据client创建实例

		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stdout)
			var Tag = []string{"primary"}
			instancer := consul.NewInstancer(client, logger, "userservice", Tag, true) //最后的true表示只有通过健康检查的服务才能被得到
			{
				factory := func(service_url string) (endpoint.Endpoint, io.Closer, error) { //factory定义了如何获得服务端的endpoint,这里的service_url是从consul中读取到的service的address我这里是192.168.3.14:8000
					tart, _ := url.Parse("http://" + service_url)                                                                                 //server ip +8080真实服务的地址
					return httptransport.NewClient("GET", tart, Services.GetUserInfo_Request, Services.GetUserInfo_Response).Endpoint(), nil, nil //我再GetUserInfo_Request里面定义了访问哪一个api把url拼接成了http://192.168.3.14:8000/v1/user/{uid}的形式
				}
				endpointer := sd.NewEndpointer(instancer, factory, logger)
				endpoints, err := endpointer.Endpoints() //获取所有的服务端当前server的所有endpoint函数
				if err != nil {
					return "", err
				}
				fmt.Println("服务有", len(endpoints), "条")

				mylb := lb.NewRandom(endpointer, time.Now().UnixNano()) //使用go-kit自带的轮询

				for {
					getUserInfo, err := mylb.Endpoint() //写死获取第一个
					ctx := context.Background()         //第三步：创建一个context上下文对象

					//第四步：执行
					res, err := getUserInfo(ctx, Services.UserRequest{Uid: 101})
					if err != nil {
						return "", err
					}
					//第五步：断言，得到响应值
					userinfo := res.(Services.UserResponse)
					return userinfo.Result, nil
				}
			}
		}
	}
}

```

> 因为文件结果混乱 具体代码不过多展示了

## 日志

基本使用

![image-20240901161536344](images/image-20240901161536344.png)

### 创建日志

```go
package Services

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/time/rate"
	"p4/utils"
	"os"
	"strconv"
)

type UserRequest struct { //封装User请求结构体
	Uid    int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
}

// 加入限流功能中间件
func RateLimit(limit *rate.Limiter) endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
	return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, utils.NewMyError(429, "toot many request")
			}
			return next(ctx, request)
		}
	}
}

func GenUserEnPoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stdout)
			logger = log.WithPrefix(logger, "mykit", "1.0")
			logger = log.WithPrefix(logger, "time", log.DefaultTimestampUTC) //加上前缀时间
			logger = log.WithPrefix(logger, "caller", log.DefaultCaller)     //加上前缀，日志输出时的文件和第几行代码

		}
		r := request.(UserRequest) //通过类型断言获取请求结构体
		result := "nothings"
		if r.Method == "GET" {
			result = userService.GetName(r.Uid) + strconv.Itoa(utils.ServicePort)
			logger.Log("method", r.Method, "event", "get user", "userid", r.Uid)

		} else if r.Method == "DELETE" {
			err := userService.DelUser(r.Uid)
			if err != nil {
				result = err.Error()
			} else {
				result = fmt.Sprintf("userid为%d的用户已删除", r.Uid)
			}
		}
		return UserResponse{Result: result}, nil
	}
}

```

> 修改 UserEndPoint 实现

### 由中间件包装Log

```go
package Services

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"golang.org/x/time/rate"
	"gomicro/utils"
	"strconv"
)

type UserRequest struct { //封装User请求结构体
	Uid    int `json:"uid"`
	Method string
}

type UserResponse struct {
	Result string `json:"result"`
}

// 日志中间件,每一个service都应该有自己的日志中间件
func UserServiceLogMiddleware(logger log.Logger) endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
	return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest) //通过类型断言获取请求结构体
			logger.Log("method", r.Method, "event", "get user", "userid", r.Uid)
			return next(ctx, request)
		}
	}
}

// 加入限流功能中间件
func RateLimit(limit *rate.Limiter) endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
	return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, utils.NewMyError(429, "toot many request")
			}
			return next(ctx, request)
		}
	}
}

func GenUserEnPoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest) //通过类型断言获取请求结构体
		result := "nothings"
		if r.Method == "GET" {
			result = userService.GetName(r.Uid) + strconv.Itoa(utils.ServicePort)

		} else if r.Method == "DELETE" {
			err := userService.DelUser(r.Uid)
			if err != nil {
				result = err.Error()
			} else {
				result = fmt.Sprintf("userid为%d的用户已删除", r.Uid)
			}
		}
		return UserResponse{Result: result}, nil
	}
}

```

调用日志中间件

```go
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

```

## JWT

### JWT的基本使用

```go
package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type UserClaim struct {
	Uname              string `json:"username"`
	jwt.StandardClaims        //嵌套了这个结构体就实现了Claim接口
}

// 创建claims结构体：并且使用标准clainms接口

func main() {
	sec := []byte("123abc")                                                               //秘钥
	token_obj := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{Uname: "xiahualou"}) //使用HS256加密算法加密，创建JWT对象
	token, _ := token_obj.SignedString(sec)                                               //把秘钥传进去，生成签名token，（也就是最终的JWT字符串）
	fmt.Println(token)

	uc := UserClaim{} //验证jwt
	getToken, _ := jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		return sec, nil //这里是对称加密，所以只要有人拿到了这个sec就可以进行访问不安全
	})
	//创建Parse函数，返回sec
	//用下面这种解析方式可以把解析后的结果保存到结构体中去
	getToken, _ = jwt.ParseWithClaims(token, &uc, func(token *jwt.Token) (i interface{}, e error) {
		return sec, nil
	})

	//验证jwt是否有效
	if getToken.Valid {
		fmt.Println(getToken.Claims.(*UserClaim).Uname) //使用断言判断具体的claim直接取值 断言为UserClaim打印Uname
	}

	// 从最后的输出信息可见 解密成功
}

// claim 是声明信息

```

### 对称加密/非对称加密-公钥/私钥

生成公钥/私钥初始化代码

```go
package utils

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "io/ioutil"
)

func GenRSAPubAndPri(bits int,filepath string ) error {
    // 生成私钥文件
    privateKey, err := rsa.GenerateKey(rand.Reader, bits)
    if err != nil {
        return err
    }
    derStream := x509.MarshalPKCS1PrivateKey(privateKey)
    priBlock := &pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: derStream,
    }

    err= ioutil.WriteFile(filepath+"/private.pem",pem.EncodeToMemory(priBlock), 0644)
    if err!=nil{
        return err
    }
    fmt.Println("=======私钥文件创建成功========")
    // 生成公钥文件
    publicKey := &privateKey.PublicKey
    derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
    if err != nil {
        return err
    }
    publicBlock := &pem.Block{
        Type:  "PUBLIC KEY",
        Bytes: derPkix,
    }

    err= ioutil.WriteFile(filepath+"/public.pem",pem.EncodeToMemory(publicBlock), 0644)
    if err!=nil{
        return err
    }
    fmt.Println("=======公钥文件创建成功=========")

    return nil
}


```

调用上面代码生成公钥/私钥

```go
package main

import (
    "gomicro/utils"
    "log"
)

func main() {
    err := utils.GenRSAPubAndPri(1024, "./pem") //1024是长度，长度越长安全性越高，但是性能也就越差
    if err != nil {
        log.Fatal(err)
    }
    //执行完生成公钥和私钥，公钥给别人私钥给自己
}

```

私钥加密

```go
package main

import (
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "io/ioutil"
    "log"
)

type UserClaim struct { //这个结构体主要是用来宣示当前公钥的使用者是谁，只有使用者和公钥的签名者是同一个人才可以用来正确的解密，还可以设置其他的属性，可以去百度一下
    Uname              string `json:"username"`
    jwt.StandardClaims        //嵌套了这个结构体就实现了Claim接口
}

func main() {
    priBytes, err := ioutil.ReadFile("./pem/private.pem")
    if err != nil {
        log.Fatal("私钥文件读取失败")
    }

    pubBytes, err := ioutil.ReadFile("./pem/public.pem")
    if err != nil {
        log.Fatal("公钥文件读取失败")
    }
    pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
    if err != nil {
        log.Fatal("公钥文件不正确")
    }

    priKey, err := jwt.ParseRSAPrivateKeyFromPEM(priBytes)
    if err != nil {
        log.Fatal("私钥文件不正确")
    }

    token_obj := jwt.NewWithClaims(jwt.SigningMethodRS256, UserClaim{Uname: "xiahualou"}) //所有人给xiahualou发送公钥加密的数据，但是只有xiahualou本人可以使用私钥解密
    token, _ := token_obj.SignedString(priKey)

    uc := &UserClaim{}
    getToken, _ := jwt.ParseWithClaims(token, uc, func(token *jwt.Token) (i interface{}, e error) { //使用私钥解密
        return pubKey, nil //这里的返回值必须是公钥，不然解密肯定是失败
    })
    if getToken.Valid { //服务端验证token是否有效
        fmt.Println(getToken.Claims.(*UserClaim).Uname)
    }

}

```

### Token设置过期时间

```go
package main

import (
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "io/ioutil"
    "log"
    "time"
)

type UserClaim struct {
    Uname              string `json:"username"`
    jwt.StandardClaims        //嵌套了这个结构体就实现了Claim接口
}

func main() {
    priBytes, err := ioutil.ReadFile("./pem/private.pem")
    if err != nil {
        log.Fatal("私钥文件读取失败")
    }

    pubBytes, err := ioutil.ReadFile("./pem/public.pem")
    if err != nil {
        log.Fatal("公钥文件读取失败")
    }
    pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
    if err != nil {
        log.Fatal("公钥文件不正确")
    }

    priKey, err := jwt.ParseRSAPrivateKeyFromPEM(priBytes)
    if err != nil {
        log.Fatal("私钥文件不正确")
    }
    user := UserClaim{Uname: "xiahualou"}
    user.ExpiresAt = time.Now().Add(time.Second * 5).Unix()      //UserClaim嵌套了jwt.StandardClaims，使用它的Add方法添加过期时间是5秒后，这里要使用unix()
    token_obj := jwt.NewWithClaims(jwt.SigningMethodRS256, user) //所有人给xiahualou发送公钥加密的数据，但是只有xiahualou本人可以使用私钥解密
    token, _ := token_obj.SignedString(priKey)
     //通过一秒一次for循环来验证过期生效
    for {
        uc := UserClaim{}
        getToken, err := jwt.ParseWithClaims(token, &uc, func(token *jwt.Token) (i interface{}, e error) { //使用私钥解密
            return pubKey, nil //这里的返回值必须是公钥，不然解密肯定是失败
        })
        if getToken.Valid { //服务端验证token是否有效
            fmt.Println(getToken.Claims.(*UserClaim).Uname)
        } else if ve, ok := err.(*jwt.ValidationError); ok { //官方写法招抄就行
            if ve.Errors&jwt.ValidationErrorMalformed != 0 {
                fmt.Println("错误的token")
            } else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
                fmt.Println("token过期或未启用")
            } else {
                fmt.Println("无法处理这个token", err)
            }

        }
        time.Sleep(time.Second)
    }

}

```















### 集成使用JWT处理Token

> 还是三步骤创建EndPoint，创建Transport，调用请求

#### 第一步创建transport

```go
package Services

import (
    "context"
    "encoding/json"
    "errors"
    "github.com/tidwall/gjson"
    "io/ioutil"
    "net/http"
)



func DecodeAccessRequest(c context.Context, r *http.Request) (interface{}, error){
    body,_:=ioutil.ReadAll(r.Body)
    result:=gjson.Parse(string(body)) //第三方库解析json
    if result.IsObject() { //如果是json就返回true
        username:=result.Get("username")
        userpass:=result.Get("userpass")
        return AccessRequest{Username:username.String(),Userpass:userpass.String(),Method:r.Method},nil
    }
    return nil,errors.New("参数错误")

}
func EncodeAccessResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
    w.Header().Set("Content-type","application/json")
    return json.NewEncoder(w).Encode(response) //返回一个bool值判断response是否可以正确的转化为json，不能则抛出异常，返回给调用方
}
```

#### 在UserTransport中修改一下代码

```go
package Services

import (
    "context"
    "encoding/json"
    "errors"
    mymux "github.com/gorilla/mux"
    "gomicro/utils"
    "net/http"
    "strconv"
)

func DecodeUserRequest(c context.Context, r *http.Request) (interface{}, error) { //这个函数决定了使用哪个request来请求
    vars := mymux.Vars(r)
    if uid, ok := vars["uid"]; ok {
        uid, _ := strconv.Atoi(uid)
        return UserRequest{Uid: uid, Method: r.Method, Token: r.URL.Query().Get("token")}, nil //请求必须携带token过来，如果找不到这里返回空字符串，因为request访问的先后顺序是先DecodeUserRequest，再EncodeUserResponse再到我们的EndPoint，所以这里就已经给我们的request结构体存入了Token，那么我们EndPoint里面的request类型断言成UserRequest结构体实例后里面就有Token了
    }
    return nil, errors.New("参数错误")
}

func EncodeUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
    w.Header().Set("Content-type", "application/json")
    return json.NewEncoder(w).Encode(response)
}

func MyErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
    contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
    w.Header().Set("Content-type", contentType) //设置请求头
    if myerr, ok := err.(*utils.MyError); ok {
        w.WriteHeader(myerr.Code)
        w.Write(body)
    } else {
        w.WriteHeader(500)
        w.Write(body)
    }

}
```

#### 创建Endpoint

这段是生成Token的代码，我这里先调用一下这段代码，拿到token后使用localhost:8080/user?token=token来访问我们的接口

```go
package Services

import (
    "context"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/go-kit/kit/endpoint"
    "time"
)
const secKey="123abc"//秘钥
type UserClaim struct {
    Uname string `json:"username"`
    jwt.StandardClaims
}
type IAccessService interface {
    GetToken(uname string,upass string) (string,error)
}
type AccessService struct {}

func(this * AccessService) GetToken (uname string,upass string ) (string,error)  {
    if uname=="jerry" && upass=="123"{
        userinfo:=&UserClaim{Uname:uname}
        userinfo.ExpiresAt=time.Now().Add(time.Second*60).Unix() //设置60秒的过期时间
        token_obj:=jwt.NewWithClaims(jwt.SigningMethodHS256,userinfo)
        token,err:=token_obj.SignedString([]byte(secKey))
        return token,err
    }
    return "",fmt.Errorf("error uname and password")
}

type AccessRequest struct {
    Username string
    Userpass string
    Method string
}
type AccessResponse struct {
    Status string
    Token string
}
func  AccessEndpoint(accessService IAccessService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (response interface{}, err error) {
        r:=request.(AccessRequest)
        result:=AccessResponse{Status:"OK"}
        if r.Method=="POST"{
            token,err:=accessService.GetToken(r.Username,r.Userpass)
            if err!=nil{
                result.Status="error:"+err.Error()
            }else{
                result.Token=token
            }
        }
        return result,nil
    }
}
```

#### 在UserEndPoint中增加CheckToken的中间件

```go
package Services

import (
    "context"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/go-kit/kit/endpoint"
    "github.com/go-kit/kit/log"
    "golang.org/x/time/rate"
    "gomicro/utils"
    "strconv"
)

type UserRequest struct { //封装User请求结构体
    Uid    int `json:"uid"`
    Method string
    Token  string
}

type UserResponse struct {
    Result string `json:"result"`
}

//token验证中间件
func CheckTokenMiddleware() endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
    return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
        return func(ctx context.Context, request interface{}) (response interface{}, err error) {
            r := request.(UserRequest) //通过类型断言获取请求结构体
            uc := UserClaim{}
            //下面的r.Token是在代码DecodeUserRequest那里封装进去的
            getToken, err := jwt.ParseWithClaims(r.Token, &uc, func(token *jwt.Token) (i interface{}, e error) {
                return []byte(secKey), err
            })
            fmt.Println(err, 123)
            if getToken != nil && getToken.Valid { //验证通过
                newCtx := context.WithValue(ctx, "LoginUser", getToken.Claims.(*UserClaim).Uname)
                return next(newCtx, request)
            } else {
                return nil, utils.NewMyError(403, "error token")
            }

            //logger.Log("method", r.Method, "event", "get user", "userid", r.Uid)

        }
    }
}

//日志中间件,每一个service都应该有自己的日志中间件
func UserServiceLogMiddleware(logger log.Logger) endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
    return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
        return func(ctx context.Context, request interface{}) (response interface{}, err error) {
            r := request.(UserRequest) //通过类型断言获取请求结构体
            logger.Log("method", r.Method, "event", "get user", "userid", r.Uid)
            return next(ctx, request)
        }
    }
}

//加入限流功能中间件
func RateLimit(limit *rate.Limiter) endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
    return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
        return func(ctx context.Context, request interface{}) (response interface{}, err error) {
            if !limit.Allow() {
                return nil, utils.NewMyError(429, "toot many request")
            }
            return next(ctx, request) //执行endpoint
        }
    }
}

func GenUserEnPoint(userService IUserService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (response interface{}, err error) {
        r := request.(UserRequest) //通过类型断言获取请求结构体
        fmt.Println("当前登录用户为", ctx.Value("LoginUser"))
        result := "nothings"
        if r.Method == "GET" {
            result = userService.GetName(r.Uid) + strconv.Itoa(utils.ServicePort)

        } else if r.Method == "DELETE" {
            err := userService.DelUser(r.Uid)
            if err != nil {
                result = err.Error()
            } else {
                result = fmt.Sprintf("userid为%d的用户已删除", r.Uid)
            }
        }
        return UserResponse{Result: result}, nil
    }
}
```

#### 调用checkToken中间件的代码

```go
package main

import (
    "flag"
    "fmt"
    kitlog "github.com/go-kit/kit/log"
    httptransport "github.com/go-kit/kit/transport/http"
    mymux "github.com/gorilla/mux"
    "golang.org/x/time/rate"
    "gomicro/Services"
    "gomicro/utils"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "syscall"
)

func main() {
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

    //用户服务
    user := Services.UserService{}
    limit := rate.NewLimiter(1, 5)
    endp := Services.RateLimit(limit)(Services.UserServiceLogMiddleware(logger)(Services.CheckTokenMiddleware()(Services.GenUserEnPoint(user))))

    //增加handler用于获取token
    accessService := &Services.AccessService{}
    accessServiceEndpoint := Services.AccessEndpoint(accessService)
    accessHandler := httptransport.NewServer(accessServiceEndpoint, Services.DecodeAccessRequest, Services.EncodeAccessResponse)

    options := []httptransport.ServerOption{
        httptransport.ServerErrorEncoder(Services.MyErrorEncoder), //使用我们的自定义错误
    }

    serverHandler := httptransport.NewServer(endp, Services.DecodeUserRequest, Services.EncodeUserResponse, options...) //使用go kit创建server传入我们之前定义的两个解析函数

    r := mymux.NewRouter()
    //r.Handle(`/user/{uid:\d+}`, serverHandler) //这种写法支持多种请求方式
    r.Methods("POST").Path("/access-token").Handler(accessHandler)            //注册token获取的handler
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
```

#### 使用postman测试接口

![img](images/20191224142350.png)

> gjson适用于解析json的三方库，速度很快，可以取Github了解下
>
> 仓库地址如下：
> https://github.com/tidwall/gjson
>
> `gjson` 一个强大且高效的 Go 语言库，用于解析和提取 JSON 数据。它通过路径表达式提供了一种简洁的方式来访问 JSON 数据，而无需将 JSON 解析为结构体或映射

> 现在我们拿到了token，那么现在我们在访问的时候需要带上token去访问接口在UserRequest中新加一个Token字段用于请求认证

```go
package Services

import (
    "context"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/go-kit/kit/endpoint"
    "github.com/go-kit/kit/log"
    "golang.org/x/time/rate"
    "gomicro/utils"
    "strconv"
)

type UserRequest struct { //封装User请求结构体
    Uid    int `json:"uid"`
    Method string
    Token  string //新加的token字段，用于读取url中的token封装进来再传递给下一层的请求处理
}

type UserResponse struct {
    Result string `json:"result"`
}

//token验证中间件
func CheckTokenMiddleware() endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
    return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
        return func(ctx context.Context, request interface{}) (response interface{}, err error) {
            r := request.(UserRequest) //通过类型断言获取请求结构体
            uc := UserClaim{}
            getToken, err := jwt.ParseWithClaims(r.Token, &uc, func(token *jwt.Token) (i interface{}, e error) { //验证token是否正确
                return []byte(secKey), err
            })
            if getToken != nil && getToken.Valid { //验证通过
                //这里很关键，验证通过后我们把用户名通过ctx传入到下一层的请求，标识当前用户已经通过验证，即GenUserEndPoint返回的endpoint方法
                newCtx := context.WithValue(ctx, "LoginUser", getToken.Claims.(*UserClaim).Uname)
                return next(newCtx, request)
            } else {
                return nil, utils.NewMyError(403, "error token")
            }

            //logger.Log("method", r.Method, "event", "get user", "userid", r.Uid)

        }
    }
}

//日志中间件,每一个service都应该有自己的日志中间件
func UserServiceLogMiddleware(logger log.Logger) endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
    return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
        return func(ctx context.Context, request interface{}) (response interface{}, err error) {
            r := request.(UserRequest) //通过类型断言获取请求结构体
            logger.Log("method", r.Method, "event", "get user", "userid", r.Uid)
            return next(ctx, request)
        }
    }
}

//加入限流功能中间件
func RateLimit(limit *rate.Limiter) endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
    return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
        return func(ctx context.Context, request interface{}) (response interface{}, err error) {
            if !limit.Allow() {
                return nil, utils.NewMyError(429, "toot many request")
            }
            return next(ctx, request) //执行endpoint
        }
    }
}

func GenUserEnPoint(userService IUserService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (response interface{}, err error) {
        r := request.(UserRequest) //通过类型断言获取请求结构体
        fmt.Println("当前登录用户为", ctx.Value("LoginUser")) //读取上面newCtx设置的用户name，判断能否处理请求，我这里是简写了，如果读不到应该拒绝处理
        result := "nothings"
        if r.Method == "GET" {
            result = userService.GetName(r.Uid) + strconv.Itoa(utils.ServicePort)

        } else if r.Method == "DELETE" {
            err := userService.DelUser(r.Uid)
            if err != nil {
                result = err.Error()
            } else {
                result = fmt.Sprintf("userid为%d的用户已删除", r.Uid)
            }
        }
        return UserResponse{Result: result}, nil
    }
}
```

> 启动server后使用postman使用localhost:8080/user/101访问结果如下,因为url没有携带token所以报错了

![img](images/20191224145308.png)

> 接下来我们使用正确的方式再次访问

![img](images/20191224151917.png)

> 一般我們把token存入到Redis中防止用戶一直請求，也起到了缓存的作用

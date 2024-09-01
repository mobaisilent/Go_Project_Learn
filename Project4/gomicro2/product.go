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
	"github.com/go-kit/kit/sd/lb"
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

				mylb := lb.NewRandom(endpointer, time.Now().UnixNano())
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
			}
		}
	}

}

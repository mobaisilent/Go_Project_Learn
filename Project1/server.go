package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux" // 注意这个mux包
)

// 导入基础包

func main() {
	router := mux.NewRouter()      // 创建路由  根据需求创建路由  ：：强大的HTTP路由库
	go h.run()                     // 在一个新的goroutine中启动实体
	router.HandleFunc("/ws", myws) // 经典的gorilla/mux路由：根据ws路径启动myws方法
	// 创建HTTP服务器并监听在127.0.0.1的8080端口
	// 使用router作为请求的处理器
	err := http.ListenAndServe("127.0.0.1:8080", router) // 监听和服务端口如果没问题就启动服务方法
	// 检查服务器启动过程中是否遇到错误
	if err != nil {
		// 如果有错误，打印错误信息
		fmt.Println("err:", err)
	}
}

//经典的gorilla/mux路由：根据ws路径启动myws方法

/*
导入必要的包：fmt 用于格式化输出，net/http 用于处理 HTTP 请求，github.com/gorilla/mux 是一个强大的 HTTP 路由库，用于创建和管理路由。
定义 main 函数，程序的入口点。在这个函数中，首先创建一个路由器 router，然后启动 hub 的运行逻辑（h.run()）在一个新的 goroutine 中，以便它可以异步处理 WebSocket 连接。
使用 router.HandleFunc("/ws", myws) 设置一个路由规则，当 HTTP 请求的路径为 /ws 时，调用 myws 函数处理该请求。这里 myws 函数的具体实现没有在代码片段中给出，但基于上下文，我们可以推断 myws 函数负责将 HTTP 连接升级到 WebSocket 连接，并处理后续的通信。
最后，使用 http.ListenAndServe("127.0.0.1:8080", router) 启动 HTTP 服务器，监听本地主机的 8080 端口。这个服务器使用之前创建的 router 作为请求的处理器，根据不同的路径分发请求到相应的处理函数。
如果服务器启动过程中遇到错误，将通过 fmt.Println("err:", err) 打印错误信息。
*/

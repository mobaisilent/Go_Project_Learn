package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket" // 这里这里依旧导入了wesocket块
)

type connection struct {
	ws   *websocket.Conn // websocket连接器
	sc   chan []byte     // 发送消息的管道
	data *Data           // 数据结构体：注意上面创建的data
}

// 创建connection对象：用于存放升级为ws的对象

var wu = &websocket.Upgrader{
	ReadBufferSize:  512,
	WriteBufferSize: 512,
	CheckOrigin:     func(r *http.Request) bool { return true }}

// 定义websocket的变量 设置三个配置信息读写buf和checkorigin  作用于http请求

func myws(w http.ResponseWriter, r *http.Request) { // 注意这里传入参数
	ws, err := wu.Upgrade(w, r, nil) // 尝试将http升级到ws连接
	if err != nil {
		return
	} //失败处理
	c := &connection{
		ws:   ws,
		sc:   make(chan []byte, 256),
		data: &Data{}}
	// 创建connection对象
	h.r <- c
	go c.writer()
	c.reader()

	defer func() {
		c.data.Type = "logout"
		user_list = del(user_list, c.data.User)
		c.data.UserList = user_list
		c.data.Content = c.data.User
		data_b, _ := json.Marshal(c.data)
		h.b <- data_b
		h.r <- c
	}()
	// 处理登出的逻辑
}

// myws的实现：看server.go那里引用了这个方法

func (c *connection) writer() {
	for message := range c.sc {
		c.ws.WriteMessage(websocket.TextMessage, message)
	}
	c.ws.Close()
}

var user_list = []string{}

func (c *connection) reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			h.r <- c
			break
		}
		json.Unmarshal(message, &c.data)
		switch c.data.Type {
		case "login":
			c.data.User = c.data.Content
			c.data.From = c.data.User
			user_list = append(user_list, c.data.User)
			c.data.UserList = user_list
			data_b, _ := json.Marshal(c.data)
			h.b <- data_b
		case "user":
			c.data.Type = "user"
			data_b, _ := json.Marshal(c.data)
			h.b <- data_b
		case "logout":
			c.data.Type = "logout"
			user_list = del(user_list, c.data.User)
			data_b, _ := json.Marshal(c.data)
			h.b <- data_b
			h.r <- c
		default:
			fmt.Print("========default================")
		}
	}
}

func del(slice []string, user string) []string {
	count := len(slice)
	if count == 0 {
		return slice
	}
	if count == 1 && slice[0] == user {
		return []string{}
	}
	var n_slice = []string{}
	for i := range slice {
		if slice[i] == user && i == count {
			return slice[:count]
		} else if slice[i] == user {
			n_slice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	fmt.Println(n_slice)
	return n_slice
}

// del数组处理

/*
connection 结构体包含三个字段：ws 是指向 websocket.Conn 的指针，用于WebSocket连接；sc 是一个字节切片的通道，用于发送消息；data 是指向 Data 结构体的指针，存储与连接相关的数据。
wu 是 websocket.Upgrader 的实例，配置了读写缓冲区大小和跨域检查函数，用于将 HTTP 连接升级到 WebSocket 连接。
myws 函数首先尝试使用 wu.Upgrade 方法将 HTTP 请求升级到 WebSocket 连接。如果升级成功，它会创建一个新的 connection 实例，并将其注册到 hub 中。然后，启动独立的 goroutine 来处理消息的发送（writer 方法）和接收（reader 方法）。最后，定义了一个 defer 语句来处理用户登出的逻辑，包括从用户列表中删除用户并通知其他用户。
writer 方法循环读取 sc 通道中的消息，并通过 WebSocket 连接发送它们。如果通道关闭，它也会关闭 WebSocket 连接。
reader 方法循环读取 WebSocket 连接上的消息。对于每条接收到的消息，它会根据消息类型（登录、普通用户消息、登出）进行不同的处理，如更新用户列表、广播消息等。
del 函数用于从字符串切片中删除指定的元素。它遍历切片，找到匹配的元素并返回一个不包含该元素的新切片。
*/

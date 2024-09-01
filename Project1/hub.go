package main

import "encoding/json"

var h = hub{
	c: make(map[*connection]bool), // 用户池
	b: make(chan []byte),          // 广播信道
	r: make(chan *connection),     // 处理注册
	u: make(chan *connection),     // 处理注销
}

// 这是一个全局对象：所以可见server里面直接调用了h.run

// 四个切片：c、u、b、r，分别用来存储连接、用户、消息和注册连接

type hub struct {
	c map[*connection]bool
	b chan []byte
	r chan *connection
	u chan *connection
}

// 这是一个结构体

func (h *hub) run() {
	for {
		select {
		case c := <-h.r:
			h.c[c] = true
			c.data.Ip = c.ws.RemoteAddr().String()
			c.data.Type = "handshake"
			c.data.UserList = user_list
			data_b, _ := json.Marshal(c.data)
			c.sc <- data_b
		case c := <-h.u:
			if _, ok := h.c[c]; ok {
				delete(h.c, c)
				close(c.sc)
			}
		case data := <-h.b:
			for c := range h.c {
				select {
				case c.sc <- data:
				default:
					delete(h.c, c)
					close(c.sc)
				}
			}
		}
	}
}

/*
方法使用一个无限循环 for，内部通过 select 语句监听三个通道：h.r、h.u 和 h.b。
当从 h.r 通道接收到一个 connection 实例时（表示有新的连接需要注册）：
将该连接添加到 h.c 映射中，标记为活跃状态。
设置连接的 IP 地址、类型为 "handshake"，并分配用户列表。
将连接数据序列化为 JSON 格式，并通过连接的发送通道 sc 发送给客户端。
当从 h.u 通道接收到一个 connection 实例时（表示有连接需要注销）：
检查该连接是否存在于 h.c 映射中，如果存在，则从映射中删除，并关闭连接的发送通道 sc。
当从 h.b 通道接收到数据时（表示有消息需要广播给所有客户端）：
遍历 h.c 映射中的所有连接，尝试通过每个连接的发送通道 sc 发送数据。
如果发送失败（通道已满或其他原因），则从映射中删除该连接，并关闭连接的发送通道 sc。
*/

### server.go

```bash
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
```

### hub.go

```go
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

```

### connection.go

```go
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

```

### data.go

```
package main

type Data struct {
	Ip       string   `json:"ip"`
	User     string   `json:"user"`
	From     string   `json:"from"`
	Type     string   `json:"type"`
	Content  string   `json:"content"`
	UserList []string `json:"user_list"`
}

// Data 结构体，包含了用户的 IP 地址、用户名、消息来源、消息类型、消息内容和用户列表。
// 后面的json其实就是标签罢了：方便使用的

```

### local.html

```html
<!DOCTYPE html>
<html lang="en">

<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>mobai演示聊天室</title>
  <style>
    body {
      font-family: Arial, sans-serif;
    }

    .chat-container {
      width: 800px;
      height: 600px;
      margin: 30px auto;
      text-align: center;
    }

    .chat-header {
      margin-bottom: 20px;
    }

    .chat-main {
      display: flex;
      border: 1px solid gray;
      height: 300px;
    }

    .user-list {
      width: 200px;
      height: 300px;
      float: left;
      text-align: left;
      overflow: auto;
    }

    .message-list {
      flex-grow: 1;
      border-left: 1px solid gray;
      overflow: auto;
    }

    .message-input {
      width: calc(100% - 20px);
      margin-top: 20px;
    }

    .send-button {
      margin-top: 10px;
    }
  </style>
</head>

<body>
<div class="chat-container">
  <h1 class="chat-header">mobai演示聊天室</h1>
  <div class="chat-main">
    <div class="user-list">
      <p>当前在线:<span id="user_num">0</span></p>
      <div id="user_list"></div>
    </div>
    <div id="msg_list" class="message-list"></div>
  </div>
  <textarea id="msg_box" class="message-input" rows="6" cols="50" onkeydown="confirm(event)"></textarea><br>
  <input type="button" class="send-button" value="发送" onclick="send()">
</div>
</body>

</html>
<script type="text/javascript">
  var uname = prompt('请输入用户名', 'user' + uuid(8, 16));
  var ws = new WebSocket("ws://127.0.0.1:8080/ws");
  ws.onopen = function () {
    var data = "系统消息：建立连接成功";
    listMsg(data);
  };
  ws.onmessage = function (e) {
    var msg = JSON.parse(e.data);
    var sender, user_name, name_list, change_type;
    switch (msg.type) {
      case 'system':
        sender = '系统消息: ';
        break;
      case 'user':
        sender = msg.from + ': ';
        break;
      case 'handshake':
        var user_info = { 'type': 'login', 'content': uname };
        sendMsg(user_info);
        return;
      case 'login':
      case 'logout':
        user_name = msg.content;
        name_list = msg.user_list;
        change_type = msg.type;
        dealUser(user_name, change_type, name_list);
        return;
    }
    var data = sender + msg.content;
    listMsg(data);
  };
  ws.onerror = function () {
    var data = "系统消息 : 出错了,请退出重试.";
    listMsg(data);
  };
  function confirm(event) {
    var key_num = event.keyCode;
    if (13 == key_num) {
      send();
    } else {
      return false;
    }
  }
  function send() {
    var msg_box = document.getElementById("msg_box");
    var content = msg_box.value;
    var reg = new RegExp("\r\n", "g");
    content = content.replace(reg, "");
    var msg = { 'content': content.trim(), 'type': 'user' };
    sendMsg(msg);
    msg_box.value = '';
  }
  function listMsg(data) {
    var msg_list = document.getElementById("msg_list");
    var msg = document.createElement("p");
    msg.innerHTML = data;
    msg_list.appendChild(msg);
    msg_list.scrollTop = msg_list.scrollHeight;
  }
  function dealUser(user_name, type, name_list) {
    var user_list = document.getElementById("user_list");
    var user_num = document.getElementById("user_num");
    while (user_list.hasChildNodes()) {
      user_list.removeChild(user_list.firstChild);
    }
    for (var index in name_list) {
      var user = document.createElement("p");
      user.innerHTML = name_list[index];
      user_list.appendChild(user);
    }
    user_num.innerHTML = name_list.length;
    user_list.scrollTop = user_list.scrollHeight;
    var change = type == 'login' ? '上线' : '下线';
    var data = '系统消息: ' + user_name + ' 已' + change;
    listMsg(data);
  }
  function sendMsg(msg) {
    var data = JSON.stringify(msg);
    ws.send(data);
  }
  function uuid(len, radix) {
    var chars = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'.split('');
    var uuid = [], i;
    radix = radix || chars.length;
    if (len) {
      for (i = 0; i < len; i++) uuid[i] = chars[0 | Math.random() * radix];
    } else {
      var r;
      uuid[8] = uuid[13] = uuid[18] = uuid[23] = '-';
      uuid[14] = '4';
      for (i = 0; i < 36; i++) {
        if (!uuid[i]) {
          r = 0 | Math.random() * 16;
          uuid[i] = chars[(i == 19) ? (r & 0x3) | 0x8 : r];
        }
      }
    }
    return uuid.join('');
  }
</script>
```
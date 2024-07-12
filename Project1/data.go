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

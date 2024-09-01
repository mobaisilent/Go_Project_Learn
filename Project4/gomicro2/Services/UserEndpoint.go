package Services

// 定义信息传送结构体
type UserRequest struct {
	Uid    int    `json:"uid"`
	Method string `json:"method"`
}

type UserResponse struct {
	Result string `json:"result"`
}

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

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

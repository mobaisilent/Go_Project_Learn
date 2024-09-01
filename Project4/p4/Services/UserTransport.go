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

func MyErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	w.Header().Set("Content-type", contentType) //设置请求头
	w.WriteHeader(429)                          //写入返回码
	w.Write(body)
}

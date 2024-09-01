package Services

import (
    "context"
    "fmt"
    "github.com/go-kit/kit/endpoint"
    "github.com/go-kit/kit/log"
    "golang.org/x/time/rate"
    "p4/utils"
    "os"
    "strconv"
)

type UserRequest struct { //封装User请求结构体
    Uid    int `json:"uid"`
    Method string
}

type UserResponse struct {
    Result string `json:"result"`
}

//加入限流功能中间件
func RateLimit(limit *rate.Limiter) endpoint.Middleware { //Middleware type Middleware func(Endpoint) Endpoint
    return func(next endpoint.Endpoint) endpoint.Endpoint { //Endpoint type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
        return func(ctx context.Context, request interface{}) (response interface{}, err error) {
            if !limit.Allow() {
                return nil, utils.NewMyError(429, "toot many request")
            }
            return next(ctx, request)
        }
    }
}

func GenUserEnPoint(userService IUserService) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (response interface{}, err error) {
        var logger log.Logger
        {
            logger = log.NewLogfmtLogger(os.Stdout)
            logger = log.WithPrefix(logger, "mykit", "1.0")
            logger = log.WithPrefix(logger, "time", log.DefaultTimestampUTC) //加上前缀时间
            logger = log.WithPrefix(logger, "caller", log.DefaultCaller) //加上前缀，日志输出时的文件和第几行代码

        }
        r := request.(UserRequest) //通过类型断言获取请求结构体
        result := "nothings"
        if r.Method == "GET" {
            result = userService.GetName(r.Uid) + strconv.Itoa(utils.ServicePort)
            logger.Log("method", r.Method, "event", "get user", "userid", r.Uid)

        } else if r.Method == "DELETE" {
            err := userService.DelUser(r.Uid)
            if err != nil {
                result = err.Error()
            } else {
                result = fmt.Sprintf("userid为%d的用户已删除", r.Uid)
            }
        }
        return UserResponse{Result: result}, nil
    }
}

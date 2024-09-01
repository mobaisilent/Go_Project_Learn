# 新手Go语言项目练手

> go项目推荐来源：
>
> https://learnku.com/articles/58970
>
> 
>
> 备注：
>
> 鉴于我需要整理这几个项目，不方便fork，出于尊重原作者的意愿，将所以项目的链接都放到project名称下面。
>
> 
>
> 愿景：
>
> Go语言的学习不是一蹴而就的，希望/也只能坚持。

## Project1

> websocket编程：
>
> mobai在线聊天。
>
> 
>
> 参考来源：
>
> https://www.topgoer.com/%E7%BD%91%E7%BB%9C%E7%BC%96%E7%A8%8B/WebSocket%E7%BC%96%E7%A8%8B.html
>
> 
>
> 评价：
>
> 直接上手websocket还是难了点：但是简单的tcp编程又比较简单。

## Project2

> gin-vue-admin
>
> 
>
> 前情提要：
>
> 项目综合性比较高；难度比较大。
>
> 
>
> 项目github地址：
> https://github.com/flipped-aurora/gin-vue-admin

## Project3

>grpc-todolist
>
>
>
>提要：
>
>使用etcd作为键值对数据储存和处理，可要为高效处理分布式高一致性场景提高经验。
>
>
>
>项目地址：
>
>https://github.com/CocaineCong/gRPC-todoList

## Project4

> 项目融入了不少其他包：虽然核心都是围绕着go-kit的

### 关于go-kit

>https://github.com/go-kit/kit

优秀的微服务工具包合集，利用它提供的API和规范可以创建健壮，可维护性高的微服务体系。



文档：

https://godoc.org/github.com/go-kit/kit



安装：

```go
go get github.com/go-kit/kit
```



类似框架：

- https://github.com/micro/go-micro
- https://github.com/koding/kite

### 微服务体系的基本需求（非全部

1. HTTPREST、RPC
2. 日志功能
3. 限流
4. API监控
5. 服务注册与发现
6. API网关服务链路追踪
7. 服务熔断

### Go-kit的三层架构

1. Transport

   主要负责与HTTP、gRPC、thrift等相关的逻辑

2. Endpoint
   定义Request和Response格式，并可以使用装饰器包装函数，
   以此来实现各种中间件嵌套。

3. Service
   这里就是我们的业务类、接口等

> 定义顺序是：Service -> Endpoint -> Transport

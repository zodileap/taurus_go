package grpc

const text_panic_server_nil = "server参数不能为nil"

const text_panic_server_register_twice = "不能在同一个gRPC服务下注册两个相同名字的服务,服务名: "

const text_panic_grpc_not_register = "该gRPC服务没有注册过grpc服务,服务名: "

const text_panic_client_connect_fail = "连接gRPC服务失败,客户端: "

const text_panic_client_connect_register_twice = "不能同一个gRPC客户端注册两次,客户端: "

const text_err_get_client_fail = "获取gRPC客户端失败,客户端: "

const text_panic_grpc_tls_fail = "gRPC在进行TLS初始化失败,错误: "

# goutils
golang 公共组件

## 组件

### log

通用日志库，go-logging封装

### context

日志上下文组件，方便做tracing及日志统一输出

### http

http请求连接池，支持配置http请求超时时间，连接池大小，能做到快速拒绝，防止服务由于超时或者大流量下造成的雪崩

### redis

redis client封装，支持集群和单机模式

# Go 压测工具

这是宇润使用 Go 语言开发的一个 Http 压测工具。

Swoole 版本：<https://github.com/Yurunsoft/net-test-swoole>

压测命令：`net-test http -u http://127.0.0.1:8080 -c 100 -n 100000`

参数说明：

```shell
--url value, -u value     压测地址
--co value, -c value      并发数（协程数量） (default: 100)
--number value, -n value  总请求次数 (default: 100)
```

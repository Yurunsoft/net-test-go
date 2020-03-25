package main

import (
	"net-test/http"
	"os"

	"github.com/urfave/cli"
)

func main() {
	//实例化cli
	app := cli.NewApp()
	//Name可以设定应用的名字
	app.Name = "压测工具"
	app.Usage = "一个用 go 语言开发的压测工具"
	// Version可以设定应用的版本号
	app.Version = "1.0.0"
	// Commands用于创建命令
	app.Commands = []cli.Command{
		{
			// 命令的名字
			Name: "http",
			// 命令的缩写，就是不输入language只输入lang也可以调用命令
			Aliases: []string{"http"},
			// 命令的用法注释，这里会在输入 程序名 -help的时候显示命令的使用方法
			Usage: "Http 压测",
			// 命令的处理函数
			Action: http.Test,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "url,u",
					Usage:    "压测地址",
					Required: true,
				},
				cli.Int64Flag{
					Name:  "co,c",
					Usage: "并发数（协程数量）",
					Value: 100,
				},
				cli.Int64Flag{
					Name:  "number,n",
					Usage: "总请求次数",
					Value: 100,
				},
			},
		},
	}
	// 接受os.Args启动程序
	app.Run(os.Args)
}

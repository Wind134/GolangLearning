package main

import (
	"flag"
	"fmt"

	rabbitmq "github.com/Wind134/GolangLearning/rabbitMQ"
)

func main() {
	// 定义命令行标志
	f := flag.String("mode", "", "Program running flag")

	// 解析命令行参数
	flag.Parse()

	flagInfo := *f

	switch flagInfo {
	case "producer":
		rabbitmq.Producer()
	case "consumer":
		rabbitmq.Consumer()
	default:
		fmt.Println("Hello World")
	}
}

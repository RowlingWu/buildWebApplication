package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"service"
)

func main() {
	port := "8080"

	flag.Parse() // 解析端口参数
	if len(flag.Args()) != 0 {
		port = flag.Args()[0]
	}

	mx := service.InitRoutes()

	defer os.Stdout.Close()

	fmt.Println("GoCron listening on", port)
	http.ListenAndServe(":"+port, mx)
}

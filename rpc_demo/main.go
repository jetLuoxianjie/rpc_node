package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"rpc_demo/client"
	"rpc_demo/server"
	"time"

)

type User struct {
	Name string
	Age  int
}

var userDB = map[int]User{
	1: User{"Ankur", 85},
	9: User{"Anand", 25},
	8: User{"Ankur Anand", 27},
}

func QueryUser(id int) (User, error) {
	if u, ok := userDB[id]; ok {
		return u, nil
	}
	return User{}, fmt.Errorf("id %d not in user db", id)
}

func main() {
	// 注册编码解码的
	gob.Register(User{})
	addr := "localhost:3212"
	srv := server.NewServer(addr)

	// 注册服务
	srv.Register("QueryUser", QueryUser)
	//运行
	go srv.Run()

	// 等待一秒免得服务器还没有启动完毕就开始链接
	time.Sleep(1 * time.Second)

	// 客户端开始
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	cli := client.NewClient(conn)

	var Query func(int) (User, error)
	//调用
	cli.CallRPC("QueryUser", &Query)

	u, err := Query(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(u)

	u2, err := Query(8)
	if err != nil {
		panic(err)
	}
	fmt.Println(u2)
}

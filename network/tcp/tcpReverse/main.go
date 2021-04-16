package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	server()
}

func server() {
	lis, err := net.Listen("tcp", ":7879")
	if err != nil {
		panic(err)
		return
	}
	defer lis.Close()
	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println("建立连接错误:%v\n", err)
			continue
		}
		fmt.Println(conn.RemoteAddr(), conn.LocalAddr())
		go handle(conn)
	}
}

func handle(sconn net.Conn) {
	defer sconn.Close()
	ip := "www.baidu.com:443"
	dconn, err := net.DialTimeout("tcp", ip, time.Second*2)
	if err != nil {
		fmt.Printf("连接%v失败:%v\n", ip, err)
		return
	}
	ExitChan := make(chan bool, 2)
	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		_, err := io.Copy(dconn, sconn)
		fmt.Printf("往%v发送数据失败:%v\n", ip, err)
		ExitChan <- true

	}(sconn, dconn, ExitChan)
	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		_, err := io.Copy(sconn, dconn)
		sconn.Close()
		fmt.Printf("从%v接收数据失败:%v\n", ip, err)
		ExitChan <- true
	}(sconn, dconn, ExitChan)
	<-ExitChan
	<-ExitChan
	dconn.Close()
}

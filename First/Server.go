package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"time"
	"unsafe"
)

func main() {
	StartListen()
}

const (
	TCP_HOST = "127.0.0.1"
	TCP_PORT = "5005"
	TCP_TYPE = "tcp"
)


func StartListen() {
	listen, err := net.Listen(TCP_TYPE, TCP_HOST+":"+TCP_PORT)
	if err != nil {
		fmt.Println("Error Listening", err.Error())
		os.Exit(1)
	}
	defer listen.Close()
	fmt.Println("Listening success ")
	for {
		conn, err := listen.Accept()
		fmt.Println(" accept connect "+conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("Error Accepting:", err.Error())
			os.Exit(1)
		}
		go ReceiveClientData(conn)
	}
}

func ReceiveClientData(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			time.Sleep(5* time.Second)
			//fmt.Println("Error reading:", err.Error())
		}else{
			receiveMsg := Bytes2String(buf)
			fmt.Println(receiveMsg)
			conn.Write([]byte("server Response :"+receiveMsg))
		}
	}
	//conn.Close()
}

func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

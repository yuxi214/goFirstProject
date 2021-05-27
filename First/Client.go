package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"unsafe"
)

var tcpAddr *net.TCPAddr
var conn *net.TCPConn = nil
var isClose bool = false
//var err error
func main()  {
	go scanner()
	tcpAddr, _ = net.ResolveTCPAddr("tcp4", "127.0.0.1:5005")
	conn,_ = net.DialTCP("tcp",nil, tcpAddr)
	go ReceiveData(conn)
	for !isClose {
		time.Sleep(1)
	}
	//conn.Close()
}

func scanner() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		sendMsg := scanner.Text()
		if sendMsg == "close" {
			fmt.Println("execute cmd:" + sendMsg)
			conn.Close()
			isClose = true
		} else {
			temp := []byte(sendMsg)
			conn.Write(temp)
		}
	}
}


func ReceiveData(conn *net.TCPConn)  {
	////read
	//dataPackage  := &tcpEcho.EchoProtocol{}
	//p,err :=dataPackage.ReadPacket(conn)
	//if err == nil {
	//	echoPacket := p.(*tcpEcho.EchoPacket)
	//	fmt.Printf("Server reply:[%v] [%v]\n", echoPacket.GetLength(), string(echoPacket.GetBody()))
	//}
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			time.Sleep(5)
		}else{
			receiveMsg :=Bytes2String(buf)
			fmt.Println(receiveMsg)
		}
	}
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func checkError(err error)  {
	if err != nil{
		log.Fatal(err)
	}
}

package main

import (
	"./Utils"
	"bufio"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"net"
)



type TcpServer struct {
	listener   *net.TCPListener
	hawkServer *net.TCPAddr
}

var (
	listenAddress = "127.0.0.1:3001"
)

func main() {
	//init server ip and port
	hawkServer, err := net.ResolveTCPAddr("tcp", listenAddress)
	Utils.CheckServerErr(err)
	// listen
	listen, err := net.ListenTCP("tcp", hawkServer)
	Utils.CheckServerErr(err)
	//close listen
	defer listen.Close()
	tcpServer := &TcpServer{
		listener:   listen,
		hawkServer: hawkServer,
	}
	fmt.Println(" start server success")
	for {
		conn, err := tcpServer.listener.Accept()
		Utils.CheckServerErr(err)
		fmt.Println("accept tcp client %s", conn.RemoteAddr().String())
		go Handle(conn)
	}
}

func Handle(conn net.Conn) {
	defer conn.Close()
	state := 0x00
	length := uint16(0)
	crc16 := uint16(0)
	var recvBuffer []byte
	cursor := uint16(0)
	bufferReader := bufio.NewReader(conn)
	for {
		recvByte,err := bufferReader.ReadByte()
		if err != nil {
			//这里因为做了心跳，所以就没有加deadline时间，如果客户端断开连接
			//这里ReadByte方法返回一个io.EOF的错误，具体可考虑文档
			if err == io.EOF {
				fmt.Printf("client %s is close!\n",conn.RemoteAddr().String())
			}
			//在这里直接退出goroutine，关闭由defer操作完成
			return
		}
		//进入状态机，根据不同的状态来处理
		switch state {
		case 0x00:
			if recvByte == 0xFF {
				state = 0x01
				//初始化状态机
				recvBuffer = nil
				length = 0
				crc16 = 0
			}else{
				state = 0x00
			}
			break
		case 0x01:
			if recvByte == 0xFF {
				state = 0x02
			}else{
				state = 0x00
			}
			break
		case 0x02:
			length += uint16(recvByte) * 256
			state = 0x03
			break
		case 0x03:
			length += uint16(recvByte)
			// 一次申请缓存，初始化游标，准备读数据
			recvBuffer = make([]byte,length)
			cursor = 0
			state = 0x04
			break
		case 0x04:
			//不断地在这个状态下读数据，直到满足长度为止
			recvBuffer[cursor] = recvByte
			cursor++
			if(cursor == length){
				state = 0x05
			}
			break
		case 0x05:
			crc16 += uint16(recvByte) * 256
			state = 0x06
			break
		case 0x06:
			crc16 += uint16(recvByte)
			state = 0x07
			break
		case 0x07:
			if recvByte == 0xFF {
				state = 0x08
			}else{
				state = 0x00
			}
		case 0x08:
			if recvByte == 0xFE {
				//执行数据包校验
				if (crc32.ChecksumIEEE(recvBuffer) >> 16) & 0xFFFF == uint32(crc16) {
					var packet Utils.Packet
					//把拿到的数据反序列化出来
					json.Unmarshal(recvBuffer,&packet)
					//新开协程处理数据
					go processRecvData(&packet,conn)
				}else{
					fmt.Println("丢弃数据!")
				}
			}
			//状态机归位,接收下一个包
			state = 0x00
		}
	}
}

func processRecvData(packet *Utils.Packet,conn net.Conn)  {
	switch packet.PacketType {
	case Utils.HEART_BEAT_PACKET:
		var beatPacket Utils.HeartPacket
		json.Unmarshal(packet.PacketContent,&beatPacket)
		fmt.Printf("recieve heat beat from [%s] ,data is [%v]\n",conn.RemoteAddr().String(),beatPacket)
		conn.Write([]byte("heartBeat\n"))
		return
	case Utils.REPORT_PACKET:
		var reportPacket Utils.ReportPacket
		json.Unmarshal(packet.PacketContent,&reportPacket)
		fmt.Printf("recieve report data from [%s] ,data is [%v]\n",conn.RemoteAddr().String(),reportPacket)
		conn.Write([]byte("Report data has recive\n"))
		return
	}
}

package main

import (
	"./Utils"
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"
)

//server ip and port
var (
	server = "127.0.0.1:3001"
)

//Client Objects
type TcpClient struct {
	connection *net.TCPConn
	hawkServer *net.TCPAddr
	stopChan   chan struct{}
}

func main() {
	hawkServer, err := net.ResolveTCPAddr("tcp", server)
	Utils.CheckClientErr(err)
	connection, err := net.DialTCP("tcp", nil, hawkServer)
	Utils.CheckClientErr(err)
	client := &TcpClient{
		connection: connection,
		hawkServer: hawkServer,
		stopChan:   make(chan struct{}),
	}
	//start receive data
	go client.receivePackets()

	//send heart packet
	go func() {
		heartBeatTick := time.Tick(2 * time.Second)
		for {
			select {
			case <-heartBeatTick:
				client.sendHeartPacket()
			case <-client.stopChan:
				return
			}
		}
	}()

	for i := 0; i < 300; i++ {
		go func() {
			sendTimer := time.After(1 * time.Second)
			for {
				select {
				case <-sendTimer:
					client.sendReportPacket()
					sendTimer = time.After(1 * time.Second)
				case <-client.stopChan:
					return
				}
			}
		}()
	}
	// await exit
	<-client.stopChan
}

func (client *TcpClient) receivePackets() {
	reader:=bufio.NewReader(client.connection)
	for{
		msg,err:=reader.ReadString('\n')
		if err!=nil{
			close(client.stopChan)
			break
		}else{
			fmt.Print(msg)
		}
	}
}

func (client *TcpClient) sendHeartPacket() {
	heartPacket:= Utils.HeartPacket{
		Version: "1.0",
		Timestamp: time.Now().Unix(),
	}
	packetBytes,err:=json.Marshal(heartPacket)
	if err!=nil{
		fmt.Println(err.Error())
	}
	packet:= Utils.Packet{
		PacketType:    Utils.HEART_BEAT_PACKET,
		PacketContent: packetBytes,
	}
	sendBytes,err:=json.Marshal(packet)
	if err!=nil{
		fmt.Println(err.Error())
	}
	client.connection.Write(Utils.EnPackSendData(sendBytes))
	fmt.Println("Send heartPacket success!")
}

func (client *TcpClient) sendReportPacket() {
	reportPacket := Utils.ReportPacket{
		Content:   Utils.GetRandString(),
		Timestamp: time.Now().Unix(),
		Rand:      rand.Int(),
	}
	packetBytes,err := json.Marshal(reportPacket)
	if err!=nil{
		fmt.Println(err.Error())
	}
	//这一次其实可以不需要，在封包的地方把类型和数据传进去即可
	packet := Utils.Packet{
		PacketType:    Utils.REPORT_PACKET,
		PacketContent: packetBytes,
	}
	sendBytes,err := json.Marshal(packet)
	if err!=nil{
		fmt.Println(err.Error())
	}
	//发送
	client.connection.Write(Utils.EnPackSendData(sendBytes))
	fmt.Println("Send metric data success!")
}


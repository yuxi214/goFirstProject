package Utils

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"os"
)

//数据包类型
const (
	HEART_BEAT_PACKET = 0x00
	REPORT_PACKET = 0x01
)
//data packet
type Packet struct {
	PacketType    byte
	PacketContent []byte
}

//heart packet
type HeartPacket struct {
	Version   string `json:"version"`
	Timestamp int64  `json:"timestamp"`
}

//report data packet
type ReportPacket struct {
	Content   string `json:"content"`
	Rand      int    `json:"rand"`
	Timestamp int64  `json:"timestamp"`
}

func CheckClientErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
func CheckServerErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func GetRandString()string  {
	length := rand.Intn(50)
	strBytes := make([]byte,length)
	for i:=0;i<length;i++ {
		strBytes[i] = byte(rand.Intn(26) + 97)
	}
	return string(strBytes)
}


//使用的协议与服务器端保持一致
func EnPackSendData(sendBytes []byte) []byte {
	packetLength := len(sendBytes) + 8
	result := make([]byte,packetLength)
	result[0] = 0xFF
	result[1] = 0xFF
	result[2] = byte(uint16(len(sendBytes)) >> 8)
	result[3] = byte(uint16(len(sendBytes)) & 0xFF)
	copy(result[4:],sendBytes)
	sendCrc := crc32.ChecksumIEEE(sendBytes)
	result[packetLength-4] = byte(sendCrc >> 24)
	result[packetLength-3] = byte(sendCrc >> 16 & 0xFF)
	result[packetLength-2] = 0xFF
	result[packetLength-1] = 0xFE
	fmt.Println(result)
	return result
}


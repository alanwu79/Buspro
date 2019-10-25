package main

import (
	conf "Buspro/config"
	_ "Buspro/crc"
	"C"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	_ "io"
	_"math/big"
	_"regexp"
	_"time"

	//"github.com/goburrow/serial"
	"github.com/tarm/serial"
	"log"
	"strconv"
	_"regexp"
)

func ByteToHex(data []byte) string {
	buffer := new(bytes.Buffer)
	for _, b := range data {

		s := strconv.FormatInt(int64(b&0xff), 16)
		if len(s) == 1 {
			buffer.WriteString("0")
		}
		buffer.WriteString(s)
	}

	return buffer.String()
}

func bytesTo32Int(b []byte) int {
	buf := bytes.NewBuffer(b)
	var tmp uint32
	binary.Read(buf, binary.BigEndian, &tmp)
	return int(tmp)
}


func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

//func Reconnect(port *serial.Port){
//	buf := make([]byte, 4960)
//	pos := 0
//	var content []byte
//	var message string
//	result := regexp.Match("aaaa")
//	log.Println(result)
//
//	for{
//		time.Sleep(1000 * time.Millisecond) //等待回傳所需的時間1000ms
//		for i := 0; i < 10; i++ {
//			bytesRead, err := port.Read(buf) //讀資料回來
//			if err != nil {
//				log.Println("Read fail")
//			}
//			if bytesRead > 0 {
//				pos += bytesRead
//				content = append(content, buf[:bytesRead]...)
//			}
//		}
//		log.Println("content=", ByteToHex(content))
//	}
//}


//func readSizedSerialPortLoop(bufSize int, port *serial.Port) error {
//	readBuf := make([]byte, 512)
//	var sumBuf = []byte{}
//	var renewBuf = []byte{}
//	sendBuf := make([]byte, 256)
//
//
//	for {
//		num, err := port.Read(readBuf)
//		if err == io.EOF {
//			fmt.Println("EOF:", num)
//			break
//		}
//		if err != nil {
//			return fmt.Errorf("cannnot open serial port: serialPort: %v, Error: %v", port, err)
//		}
//		if num > 0 {
//			log.Println("readBuf: %v, len: %v", readBuf, len(readBuf))
//			for index := range readBuf[:num] {
//				sumBuf = append(sumBuf, readBuf[index])
//			}
//			for len(sumBuf) >= bufSize {
//				sendBuf = sumBuf[:bufSize]
//
//				// Truncate sumBuf by Size
//				log.Println("sumBuf: %v, len: %v", sumBuf, len(sumBuf))
//				log.Println("sumBuf: %v, len: %v", sendBuf, len(sendBuf))
//
//				renewBuf = []byte{}
//				for index := bufSize; index < len(sumBuf); index++ {
//					renewBuf = append(renewBuf, sumBuf[index])
//				}
//				sumBuf = renewBuf
//				log.Println("renewed sumBuf: %v, len: %v / Size: %v", sumBuf, len(sumBuf), bufSize)
//			}
//		}
//		time.Sleep(1000 * time.Millisecond) //等待回傳所需的時間1000ms
//	}
//}

func main() {
	//開啟SerialPort
	port, err := serial.OpenPort(conf.Conf)
	if err != nil {
		log.Fatal("Comport open fail")
	}
	defer port.Close() //程式結束時關閉SerialPort



	wbuf := []byte{0xAA, 0xAA, 0x10, 0x01, 0xFA, 0x23, 0x9C, 0x00, 0x31, 0x01, 0x05, 0x01, 0x64, 0x00, 0x00, 0x01,
		0x26, 0xA7}

	_, err = port.Write(wbuf) //寫資料出去
	if err != nil {
		log.Fatal(err)
	}

//	go readSizedSerialPortLoop(4096,port)

	readBuf := make([]byte, 512)
	var sumBuf = []byte{}
	var renewBuf = []byte{}
	sendBuf := make([]byte, 256)
	bufSize := 20

	for {
		num, err := port.Read(readBuf)
		if err == io.EOF {
			fmt.Println("EOF:", num)
			break
		}
		if err != nil {
			 fmt.Errorf("cannnot open serial port: serialPort: %v, Error: %v", port, err)
		}
		if num > 0 {
			log.Println("enter")
			//log.Println("readBuf:", readBuf ,"len:" , len(readBuf))
			for index := range readBuf[:num] {
				sumBuf = append(sumBuf, readBuf[index])
				log.Println("index:", index,"readBuf[index]:", readBuf[index] ,"sumBuf:" , sumBuf)
			}
			for len(sumBuf) >= bufSize {
				sendBuf = sumBuf[:bufSize]

				// Truncate sumBuf by Size
				log.Println("sumBuf: ", sumBuf, "len(sumBuf): ", len(sumBuf))
				log.Println("sendBuf: ", sendBuf, "len(sendBuf): ", len(sendBuf))
				renewBuf = []byte{}
				for index := bufSize; index < len(sumBuf); index++ {
					renewBuf = append(renewBuf, sumBuf[index])
				}
				sumBuf = renewBuf
				log.Println("renewed sumBuf: %v, len: %v / Size: %v", sumBuf, len(sumBuf), bufSize)
			}
			log.Println("break")
		}
	}
	//buf := make([]byte, 1024)
	////pos := 0
	//var content []byte
	//var commandSize []byte
	//var packageSize int64
	//var headSize int64= 2
	//var commandSizeHex string
	//
	//time.Sleep(10 * time.Millisecond) //等待回傳所需的時間1000ms
	//for {
	//	bytesRead, err := port.Read(buf) //讀資料回來
	//	 if err != nil {
	//	 	log.Fatal(err)
	//	 }
	//	 if bytesRead > 0 {
	//	 	log.Println("buf = ",buf)
	//	 	log.Println("bytesRead = ",bytesRead)
	//
	//	 	commandSize = buf[2:3]
	//	 	commandSizeHex= ByteToHex(commandSize)
	//	 	commandSizeInt, err := strconv.ParseInt(commandSizeHex, 16, 32)
	//	 	if err != nil {
	//	 		panic(err)
	//	 	}
	//	 	log.Println("commandSizeInt = ",commandSizeInt)
	//	 	packageSize = commandSizeInt+headSize
	//	 	log.Println("packageSize = ",packageSize)
	//		 content = append(content, buf[:commandSizeInt]...)
	//	 	//log.Println("content1=", ByteToHex(content[:packageSize]))
	//	 }
	//
	//}
	//log.Println("content2=", ByteToHex(content[:packageSize]))
}


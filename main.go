package main

import (
	conf "Buspro/config"
	_ "Buspro/crc"

	//"C"
	"Buspro/phyDev"
	_ "Buspro/transfer"

	//"Buspro/sockServicce"
	"bytes"
	_ "encoding/asn1"
	"encoding/binary"
	"fmt"
	_ "io"
	_ "math/big"
	_ "regexp"
	_ "sync"
	_ "time"

	//"github.com/goburrow/serial"
	"log"
	_ "regexp"
	"strconv"

	"github.com/tarm/serial"
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

func bytesToIntU(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0}, b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp uint8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp uint16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp uint32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0, fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}

func main() {

	port, err := serial.OpenPort(conf.Conf) //開啟SerialPort
	if err != nil {
		log.Fatal("Comport open fail")
	}
	defer port.Close() //程式結束時關閉SerialPort

	//sockServicce.RunUnix()

	readPipe := make(chan []byte)
	writepipe := make(chan []byte, 115200)
	wbuf := []byte{0xAA, 0xAA, 0x10, 0x01, 0xFA, 0x23, 0x9C, 0x00, 0x31, 0x01, 0x05, 0x01, 0x64, 0x00, 0x00, 0x01,
		0x26, 0xA7}

	writepipe <- wbuf
	phyDev.Wg.Add(1)
	go phyDev.ReadSerialPort(port, readPipe)
	//phyDev.WriteSerial(port, writepipe)

	phyDev.Wg.Wait()

}

//readBuf := make([]byte, 512)
//var sumBuf = []byte{}
//var renewBuf = []byte{}
//sendBuf := make([]byte, 256)
//var commandSize = 0
//bufSize := 4096
//headSize := 2
//
//for {
//	num, err := port.Read(readBuf)
//	if err == io.EOF {
//		fmt.Println("EOF:", num)
//		break
//	}
//	if err != nil {
//		 fmt.Errorf("cannnot open serial port: serialPort: %v, Error: %v", port, err)
//	}
//	if num > 0 {
//		log.Println("num:",num)
//		//log.Println("readBuf:", readBuf ,"len:" , len(readBuf))
//		for index := range readBuf[:num] {
//
//			sumBuf = append(sumBuf, readBuf[index])
//
//			log.Println("index:", index,"readBuf[index]:", readBuf[index] ,"sumBuf:" , sumBuf)
//		}
//
//		if(ByteToHex(sumBuf[2:3])!="00") {
//			commandSize,err = bytesToIntU(sumBuf[2:3])
//			if err != nil {
//				log.Println("bytesToIntU Fail:")
//				break
//			}
//
//			bufSize = commandSize + headSize
//
//			log.Println("bufSize:" , bufSize)
//		}
//
//		for len(sumBuf) >= bufSize {
//			sendBuf = sumBuf[:bufSize]
//
//			// Truncate sumBuf by Size
//			log.Println("sumBuf: ", sumBuf, "len(sumBuf): ", len(sumBuf))
//			log.Println("sendBuf: ", sendBuf, "len(sendBuf): ", len(sendBuf))
//			renewBuf = []byte{}
//			for index := bufSize; index < len(sumBuf); index++ {
//				renewBuf = append(renewBuf, sumBuf[index])
//				log.Println("renewBuf[0]:",renewBuf[0])
//				if(renewBuf[0])!=byte(170) {
//					renewBuf = renewBuf[1:]
//				}
//			}
//			log.Println("renewBuf: ",renewBuf)
//			sumBuf = renewBuf
//			log.Println("renewed sumBuf: ",sumBuf, "len:",len(sumBuf) ,"Size:",bufSize )
//		}
//		log.Println("break")
//	}
//}
//

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

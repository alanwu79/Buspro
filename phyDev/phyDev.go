package phyDev

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"sync"

	//"github.com/goburrow/serial"
	"io"
	"log"

	"github.com/stellar/go/crc16"
	"github.com/tarm/serial"
)

var Wg sync.WaitGroup

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

func WriteSerial(port *serial.Port, writepipe <-chan []byte) {

	for {
		log.Println("write")

		command := <-writepipe
		log.Println(command)
		n, err := port.Write([]byte(command))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf(string(n))
	}
}

func CheckHeadAndLen(sumBuf []byte, bufSize int) ([]byte, int) {
	commandSize := 0
	headSize := 2

	for j := range sumBuf[:] {
		if j == 0 {
			if (sumBuf[j] != byte(170) && sumBuf[j+1] != byte(170)) || (sumBuf[j] == byte(170) && sumBuf[j+1] != byte(170)) {
				sumBuf = append(sumBuf[2:])
			} else if sumBuf[j] != byte(170) && sumBuf[j+1] == byte(170) {
				sumBuf = append(sumBuf[1:])
			}
		}

		if j == 2 {
			if sumBuf[j] > byte(78) || sumBuf[j] <= byte(11) {
				sumBuf = append(sumBuf[3:])
			} else {
				commandSize, _ = bytesToIntU(sumBuf[2:3])
				bufSize = commandSize + headSize
			}
		}
	}
	return sumBuf, bufSize
}

func ReadSerialPort(port *serial.Port, readpipe chan []byte) {

	readBuf := make([]byte, 256)
	var sumBuf = []byte{}
	sendBuf := make([]byte, 256)
	bufSize := 4096

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
			for index := range readBuf[:num] {
				sumBuf = append(sumBuf, readBuf[index])
				log.Println("readBuf[index]:", readBuf[index])
				log.Println("sumBuf:", sumBuf, "Len (sumBuf):", len(sumBuf))
			}

			if len(sumBuf) > 3 {
				sumBuf, bufSize = CheckHeadAndLen(sumBuf, bufSize)
				// for j := range sumBuf[:] {
				// if j == 0 {
				// 	if (sumBuf[j] != byte(170) && sumBuf[j+1] != byte(170)) || (sumBuf[j] == byte(170) && sumBuf[j+1] != byte(170)) {
				// 		sumBuf = append(sumBuf[2:])
				// 	} else if sumBuf[j] != byte(170) && sumBuf[j+1] == byte(170) {
				// 		sumBuf = append(sumBuf[1:])
				// 	}
				// }

				// if j == 2 {
				// 	if sumBuf[j] > byte(78) || sumBuf[j] <= byte(11) {
				// 		sumBuf = append(sumBuf[3:])
				// 	} else {
				// 		commandSize, _ = bytesToIntU(sumBuf[2:3])
				// 		bufSize = commandSize + headSize
				// 	}
				// }
				// }
			}

			if len(sumBuf) >= bufSize {
				//fmt.Print("CRC:%v", crc16.Checksum(sumBuf[2:bufSize-2]))
				su, _ := bytesToIntU(sumBuf[bufSize-2 : bufSize])
				cr, _ := bytesToIntU(crc16.Checksum(sumBuf[2 : bufSize-2]))
				if su != cr {
					sumBuf = append(sumBuf[bufSize:])
				}
			}

			for len(sumBuf) >= bufSize {
				sendBuf = sumBuf[:bufSize]
				readpipe <- sendBuf
				// Truncate sumBuf by Size
				log.Println("sumBuf: ", sumBuf, "len(sumBuf): ", len(sumBuf))
				log.Println("sendBuf: ", sendBuf, "len(sendBuf): ", len(sendBuf))
				sumBuf = append(sumBuf[bufSize:len(sumBuf)])
				log.Println("Renew sumBuf:", sumBuf, "Renew Len :", len(sumBuf))
			}
		}
	}
	Wg.Done()
}

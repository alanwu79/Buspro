package phyDev

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"sync"

	"Buspro/transfer"

	//"github.com/goburrow/serial"
	"io"
	"log"

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

func ReadSerialPort(port *serial.Port, readpipe chan []byte) {
	log.Println("read")

	readBuf := make([]byte, 256)
	var sumBuf = []byte{}
	var renewBuf = []byte{}
	sendBuf := make([]byte, 256)
	var commandSize = 0
	bufSize := 4096
	headSize := 2

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
			log.Println("num:", num)
			//log.Println("readBuf:", readBuf ,"len:" , len(readBuf))
			for index := range readBuf[:num] {

				sumBuf = append(sumBuf, readBuf[index])

				log.Println("index:", index, "readBuf[index]:", readBuf[index], "sumBuf:", sumBuf)
			}

			if ByteToHex(sumBuf[2:3]) != "00" {
				commandSize, err = bytesToIntU(sumBuf[2:3])
				if err != nil {
					log.Println("bytesToIntU Fail:")
					break
				}

				bufSize = commandSize + headSize

				log.Println("bufSize:", bufSize)
			}

			for len(sumBuf) >= bufSize {
				sendBuf = sumBuf[:bufSize]
				readpipe <- sendBuf
				go transfer.ReadPipeTransJSON(readpipe)

				// Truncate sumBuf by Size
				log.Println("sumBuf: ", sumBuf, "len(sumBuf): ", len(sumBuf))
				log.Println("sendBuf: ", sendBuf, "len(sendBuf): ", len(sendBuf))
				renewBuf = []byte{}
				for index := bufSize; index < len(sumBuf); index++ {
					renewBuf = append(renewBuf, sumBuf[index])
					log.Println("renewBuf[0]:", renewBuf[0])
					if (renewBuf[0]) != byte(170) {
						renewBuf = renewBuf[1:]
					}
				}
				log.Println("renewBuf: ", renewBuf)
				sumBuf = renewBuf
				log.Println("renewed sumBuf: ", sumBuf, "len:", len(sumBuf), "Size:", bufSize)
			}
			log.Println("break")
		}
	}
	Wg.Done()
}

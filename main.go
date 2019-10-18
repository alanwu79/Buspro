package main

import (
	conf "Buspro/config"
	_ "Buspro/crc"
	"C"
	"github.com/goburrow/serial"
	"log"
	"time"
)

func main() {
	//開啟SerialPort
	port, err := serial.Open(conf.Conf)
	if err != nil {
		log.Fatal("Comport open fail")
	}
	defer port.Close() //程式結束時關閉SerialPort



	wbuf := []byte{0xAA, 0xAA, 0x10, 0x01, 0xFA, 0x23, 0x9C, 0x00, 0x31, 0x01, 0x05, 0x01, 0x64, 0x00, 0x00, 0x01,
		0x26, 0xA7}

	// length := (bytes.Count([]byte(wbuf), nil) - 1) + 11
	// fmt.Printf("length:%X \n", length)
	// checksum := CRC16Sum(wbuf)

	// int16buf := new(bytes.Buffer)

	// binary.Write(int16buf, binary.LittleEndian, checksum)

	// fmt.Printf("output-before:%X \n", wbuf)

	// wbuf = append(wbuf, int16buf.Bytes()...)

	// fmt.Printf("output-after:%X \n", wbuf)
	_, err = port.Write(wbuf) //寫資料出去
	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, 4960)
	pos := 0
	var content []byte

	time.Sleep(100 * time.Millisecond) //等待回傳所需的時間1000ms
	for i := 0; i < 10; i++ {
		log.Printf("i=%d\n", i)
		bytesRead, err := port.Read(data) //讀資料回來
		//content = append(content, buffer[:bytesRead]...)
		if err != nil {
			log.Println("gg")
		}
		if bytesRead > 0 {
			pos += bytesRead
			content = append(content, data[:bytesRead]...)
		}
	}
	log.Println("content=", content)
}

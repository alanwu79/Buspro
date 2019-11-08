package main

import (
	conf "Buspro/config"
	"Buspro/transfer"

	//"C"
	"Buspro/phyDev"

	_ "Buspro/sockServicce"
	"bytes"
	_ "encoding/asn1"
	"encoding/binary"
	_ "encoding/json"

	"fmt"
	_ "io/ioutil"
	_ "math/big"
	_ "net"
	_ "regexp"
	_ "sync"
	_ "time"

	//"github.com/goburrow/serial"
	"log"
	_ "regexp"
	"strconv"

	_ "github.com/gogf/gf/encoding/gjson"
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

// func JsonUnmarshal() {
// 	testjson, err := ioutil.ReadFile("testjson.dat")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	var jsonObj map[string]interface{}
// 	err = json.Unmarshal(testjson, &jsonObj)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	if _, ok := jsonObj["BusProDevList"]; !ok {
// 		fmt.Println("BusProDevList err")
// 		return
// 	}
// 	BusProDevList := jsonObj["BusProDevList"].([]interface{})

// 	for _, b := range BusProDevList {

// 		BusProDev := b.(map[string]interface{})
// 		subnetID := int(BusProDev["subnetID"].(float64))
// 		fmt.Println(subnetID)

// 		// devices := jsonObj[BusProDev["devices"]].([]interface{})
// 		// for _, d := range devices {
// 		// 	device := d.(map[string]interface{})
// 		// 	deviceType := int(BusProDev["subnetID"].(float64))
// 		// 	fmt.Println(deviceType)
// 		// }
// 	}
// }

func main() {

	port, err := serial.OpenPort(conf.Conf) //開啟SerialPort
	if err != nil {
		log.Fatal("Comport open fail")
	}
	defer port.Close() //程式結束時關閉SerialPort

	// addr, err := net.ResolveUnixAddr("unix", unixSocketFile)
	// checkError(err)
	// unixCon, err := net.DialUnix("unix", nil, addr)

	//sockServicce.RunUnix()

	transfer.ReadFromSocketTransCmd()
	readPipe := make(chan []byte)
	writepipe := make(chan []byte, 115200)
	//	sendSocketpipe := make(chan []byte)

	wbuf := []byte{0xAA, 0xAA, 0x10, 0x01, 0xFA, 0x23, 0x9C, 0x00, 0x31, 0x01, 0x05, 0x01, 0x64, 0x00, 0x00, 0x01,
		0x26, 0xA7}

	writepipe <- wbuf
	phyDev.Wg.Add(1)
	go phyDev.ReadSerialPort(port, readPipe)
	transfer.ReadFromDevTransJSON(readPipe)

	//go transfer.ReadPipeTransJSON(readPipe, sendSocketpipe)
	//sockServicce.ReadFromTransAndSendUnix(unixCon, sendSocketpipe)
	//phyDev.WriteSerial(port, writepipe)

	phyDev.Wg.Wait()

}

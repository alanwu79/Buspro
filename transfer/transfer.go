package transfer

import (
	"bytes"
	"encoding/binary"
	_ "encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/gogf/gf/encoding/gjson"

	//"github.com/goburrow/serial"
	"log"
)

type CommandFormat struct {
	SourceId []byte `json:"source_address"`
	TargetId []byte `json:"target_address"`
	Opcode   []byte `json:"operate_code"`
	Payload  []byte `json:"payload"`
}

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

func ReadFromDevTransJSON(readpipe chan []byte) {
	log.Println("enter ReadFromDevTransJSON")

	var m CommandFormat
	command := make([]byte, 256)
	for {
		command = <-readpipe
		msglen := command[2]
		//channel <- msg

		m.SourceId = command[3:5]
		m.TargetId = command[9:11]
		m.Opcode = command[7:9]
		m.Payload = command[11:msglen]

		// fmt.Println("sourceId:", m.SourceId)
		// fmt.Println("targetId:", m.TargetId)
		// fmt.Println("Opcode:", m.Opcode)
		// fmt.Println("payload:", m.Payload)

		// jsonBytes, err := json.Marshal(m)
		// if err != nil {
		// 	log.Println("bytesToIntU Fail:")
		// }

		//fmt.Printf("json 结果:%s\n", jsonBytes)
	}
}

// wbuf := []byte{0xAA, 0xAA, 0x10, 0x01, 0xFA, 0x23, 0x9C, 0x00, 0x31, 0x01, 0x05, 0x01, 0x64, 0x00, 0x00, 0x01,0x26, 0xA7}

func ReadFromSocketTransCmd() {
	testjson, err := ioutil.ReadFile("testjson.dat")
	if err != nil {
		fmt.Println(err)
		return
	}

	if j, err := gjson.DecodeToJson(testjson); err != nil {
		fmt.Println(err)
	} else {

		fmt.Println("deviceType:", j.GetInt("BusProDevList.0.devices.1.deviceType"))
		deviceType := fmt.Sprintf("deviceType:%x\n", j.GetInt("BusProDevList.0.devices.1.deviceType"))
		fmt.Println(deviceType)
		subnetId := fmt.Sprintf("subnetId:%x\n", j.GetInt("BusProDevList.0.devices.1.commInfo.BusPro.subnetId"))
		fmt.Println(subnetId)
		deviceId := fmt.Sprintf("deviceId:%x\n", j.GetInt("BusProDevList.0.devices.1.commInfo.BusPro.deviceId"))
		fmt.Println(deviceId)
		opCode := fmt.Sprintf("opCode:%x\n", j.GetInt("BusProDevList.0.devices.1.attrs.1.opCode"))
		fmt.Println(opCode)
		chId := fmt.Sprintf("chId:%x\n", j.GetInt("BusProDevList.0.devices.1.attrs.0.chId"))
		fmt.Println(chId)
	}
	send_buf := []byte{}
	send_buf = append(send_buf, 0xAA)
	send_buf = append(send_buf, 0xAA)

	// length_of_data_package = 11 + len(payload)
	// send_buf = append(send_buf, length_of_data_package)
	// send_buf = append(send_buf, length_of_data_package)
	// send_buf = append(send_buf, length_of_data_package)

}

// func ReadPipeTransJSON(readpipe chan []byte, sendSocketpipe chan []byte) {
// 	log.Println("enter ReadPipeTransJSON")

// 	var m CommandFormat
// 	command := make([]byte, 256)
// 	for {
// 		command = <-readpipe
// 		msglen := command[2]
// 		//channel <- msg

// 		m.SourceId = command[3:5]
// 		m.TargetId = command[9:11]
// 		m.Opcode = command[7:9]
// 		m.Payload = command[11:msglen]

// 		fmt.Println("sourceId:", m.SourceId)
// 		fmt.Println("targetId:", m.TargetId)
// 		fmt.Println("Opcode:", m.Opcode)
// 		fmt.Println("payload:", m.Payload)

// 		jsonBytes, err := json.Marshal(m)
// 		if err != nil {
// 			log.Println("bytesToIntU Fail:")
// 		}

// 		fmt.Printf("json 结果:%s\n", jsonBytes)
// 		sendSocketpipe <- jsonBytes
// 	}
// }

//
//OperateCode := map[byte]string{
//"0031 ": "SingleChannelControl ",
//"0032": "SingleChannelControlResponse ",
//"0033": "ReadStatusOfChannels ",
//"0034": "ReadStatusOfChannelsResponse ",
//"0002": "SceneControl ",
//"0003": "SceneControlResponse ",
//"E01C": "UniversalSwitchControl",
//"E01D": "UniversalSwitchControlResponse",
//"E018": "ReadStatusOfUniversalSwitch",
//"E019": "ReadStatusOfUniversalSwitchResponse  ",
//"E017": "BroadcastStatusOfUniversalSwitch  ",
//"1644": "BroadcastSensorStatusResponse   ",
//"1645": "ReadSensorStatus   ",
//"1646": "ReadSensorStatusResponse   ",
//"1647": "BroadcastSensorStatusAutoResponse  ",
//"E3E5": "BroadcastTemperatureResponse  ",
//"1944": "ReadFloorHeatingStatus    ",
//"1945": "ReadFloorHeatingStatusResponse   ",
//"1946": "ControlFloorHeatingStatus   ",
//"1947": "ControlFloorHeatingStatusResponse   ",
//"15CE": "ReadDryContactStatus    ",
//"15CF": "ReadDryContactStatusResponse    ",
//}

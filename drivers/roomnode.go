package drivers

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"

	"github.com/jcw/jeebus"
)

//The struct reflects the input data
type RoomNode struct {
	Misc  uint8
	Light uint8
	Humi  uint8
	Temp  int16
}

type RoomNodeDecoder struct {
}

func (s *RoomNodeDecoder) Handle(m *jeebus.Message) {
	b, err := hex.DecodeString(m.Get("text"))
	check(err)
	r := bytes.NewReader(b)
	l := len(b)
	mp := make(map[string]interface{})
	if l == 2 {
		//motion detect message
		mp["detect"] = (b[1] & 1)
	} else {
		var v RoomNode
		if l == 6 {
			r.Seek(1, 0)
		}
		err := binary.Read(r, binary.LittleEndian, &v)
		check(err)
		if l == 6 {
			mp["seq"] = v.Misc
		}
		mp["light"] = v.Light
		mp["humi"] = v.Humi >> 1
		if (v.Humi & 1) == 1 {
			mp["detect"] = (v.Humi & 1)
		}
		if v.Temp &= 0x3FF; v.Temp >= 0x200 {
			v.Temp -= 0x400
		}
		mp["temp"] = v.Temp
	}
	publish("roomnode", &mp, m)
}

func init() {
	register("roomnode", &RoomNodeDecoder{})
}

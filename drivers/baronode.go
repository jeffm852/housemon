package drivers

import (
	"bytes"
	"encoding/hex"
	"encoding/binary"

	"github.com/jcw/jeebus"
)

type BaroNode struct {
	//match input data
	NID       uint8        `json:"-"`
	Seq       uint8        `json:"seq,omitempty"`
	MsgT      uint8        `json:"-"`
	Temp      uint16       `json:"temp"`
	Pres      uint32       `json:"pres"`
}

type BaroNodeDecoder struct {
}

func (s *BaroNodeDecoder) Handle(m *jeebus.Message) {
	b, err := hex.DecodeString(m.Get("text"))
	check(err)
	var v BaroNode
	err = binary.Read(bytes.NewReader(b), binary.LittleEndian, &v)
	check(err)
	publish("baronode", &v, m)
}

func init() {
	register("baronode", &BaroNodeDecoder{})
}

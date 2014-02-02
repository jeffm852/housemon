package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"regexp"

	"github.com/jcw/jeebus"
	"github.com/jcw/housemon/drivers"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalf("usage: housemon <cmd> ...")
	}

	switch os.Args[1] {

	case "decode":
		client := jeebus.NewClient()
		client.Register("rd/RF12demo/#", &RF12demoDecodeService{})

		//drivers.JNodeMap()
		msg := map[string]interface{}{"text": "v"}
		jeebus.Publish("if/RF12demo", msg)
		msg = map[string]interface{}{"text": "c"}
		jeebus.Publish("if/RF12demo", msg)
		//"text":" A i1 g178 @ 868 MHz "


		<-make(chan byte) // wait forever

	default:
		log.Fatal("unknown sub-command: housemon ", os.Args[1], " ...")
	}
}

var node, grp, band int
var astx bool
var (
	rf12msg struct {
		ID   [3]int `json:"id"`
		Dev  string `json:"dev"`
		Loc  string `json:"loc"`
		Text string `json:"text"`
		Time int64  `json:"time"`
	}
)
var confRegex = regexp.MustCompile(`^ [A-Z[\\\]\^_@] i(\d+)(\*)? g(\d+) @ (\d\d\d) MHz`)

type RF12demoDecodeService struct {
}

func (s *RF12demoDecodeService) Handle(m *jeebus.Message) {
	text := m.Get("text")

	if strings.HasPrefix(text, "OK ") {
		var buf bytes.Buffer
		var vals []string
		vals = strings.Split(text[3:], " ")
		rnode , err := strconv.Atoi(vals[0])
		rnode &= 0x1F
		check(err)
		// convert the line of decimal byte values to a byte buffer
		for _, v := range vals {
			n, err := strconv.Atoi(v)
			check(err)
			buf.WriteByte(byte(n))
		}
		now := m.GetInt64("time")
		dev := strings.SplitN(m.T, "/", 3)[2]
		hex := fmt.Sprintf("%X", buf.Bytes())
		fmt.Printf("%d %s %s\n", now, dev, hex)
		rf12msg.ID[0] = band
		rf12msg.ID[1] = grp
		rf12msg.ID[2] = rnode
		rf12msg.Dev = dev
		rf12msg.Text = hex
		rf12msg.Time = now
		if found, nT, nL := drivers.JNodeType(band, grp, rnode, now); found {
			rf12msg.Loc = nL
			jeebus.Publish("rf12/"+nT, rf12msg)
		} else {
			jeebus.Publish("rf12/unknown", rf12msg)
		}
	} else if conf := confRegex.FindStringSubmatch(text); conf != nil {
		node, _ = strconv.Atoi(conf[1])
		astx = conf[2] == "*"
		grp, _ = strconv.Atoi(conf[3])
		band, _ = strconv.Atoi(conf[4])
		fmt.Println("groupID: ", grp)
	} else if strings.HasPrefix(text, "[") && strings.Contains(text, "]") {
		fmt.Println(text)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

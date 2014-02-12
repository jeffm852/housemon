package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"regexp"
	_"reflect"
	ring "container/ring"

	"github.com/jcw/jeebus"
	"github.com/jcw/housemon/drivers"
	"time"
)

//TODO: This is a quick fix update/upgrade to get RF12demo decoder working with registry system
//There are lots of optimisations that can now be made

func main() {
	var mqttAddr string

	flag.StringVar(&mqttAddr, "mqtt", ":1883",
		"connect to MQTT server on <host><:port>")
	flag.Parse()

	if !strings.Contains(mqttAddr, "://") {
		mqttAddr = "tcp://" + mqttAddr
	}
	murl, err := url.Parse(mqttAddr)
	check(err)

	client := jeebus.NewClient(murl)

	decoder := &RF12demoDecodeService{ &RF12demoDb{ make(map[string]*rf12msg),make(map[string]bool) },client, ring.New(10) }
	client.Register("io/RF12demo/+/#", decoder)  // + represents the 'instance' of the Tag
	
	<- client.Done

}


//TODO: Quick reconfig - NOT to be confused with .ino HEADER positions.
const (
	RF12BAND = iota
	RF12GRP
	RF12NODEID
)


type rf12msg struct {
	ID   [3]int `json:"id"`
	Dev  string `json:"dev"`
	Loc  string `json:"loc"`
	Text string `json:"text"`
	Time int64  `json:"time"`
	Astx bool
	Node int
}

//this contains the 'instances' of RF12 we have seen key=topic
type RF12demoDb struct {
	Db map[string]*rf12msg
	ConfigFlag map[string]bool
}

//"text":" A i1 g178 @ 868 MHz "
var confRegex = regexp.MustCompile(`^ [A-Z[\\\]\^_@] i(\d+)(\*)? g(\d+) @ (\d\d\d) MHz`)

type RF12demoDecodeService struct {
	*RF12demoDb
	client *jeebus.Client
	MsgBuffer *ring.Ring  //TODO: Add in message buffer when init moved outside (to capture OK with no config)
}

func (s *RF12demoDecodeService) Handle(m *jeebus.Message) {

	//fmt.Println("Message:", string(m.P) )

	keys := strings.Split(m.T,"/")
	rwTopic := strings.Join( keys[:len(keys)-1], "/")
	text := m.Get("text")

	inst,ok := s.Db[rwTopic]
	if !ok  {
		inst  = &rf12msg{} //not seen this rf12, create its 'state container'
		s.Db[rwTopic] = inst
		s.ConfigFlag[rwTopic] = false  //only tracks if we have sent init
	}

	if !s.ConfigFlag[rwTopic] {
		//hack: if we never seen a config from this rf12, ask for one
		//TODO: put this on watchdog timer
		s.ConfigFlag[rwTopic] = true
		//TODO: direct clone of prev code - only directed to specific 'instance' of rf12demo.
		msg := map[string]interface{}{"text": "c"}
		s.client.Publish(rwTopic, msg)
		<-time.After(50 * time.Millisecond ) //Added to help sketch
		msg = map[string]interface{}{"text": "v"}
		s.client.Publish(rwTopic, msg)

	}

	//TODO:ring buffer would wrap OK processing for a 'onetime' empty

	//TODO: This mainly unchanged apart from using 'instance' data - no attempt to tidy - yet!
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
		inst.Dev = rwTopic   //we can simply use the endpoint minus the timestamp instead of strings.SplitN(m.T, "/", 3)[2]
		hex := fmt.Sprintf("%X", buf.Bytes())

		fmt.Printf("%d %s %s\n", now, inst.Dev, hex)

		inst.ID[RF12NODEID] = rnode  //this changes possibly every message rec'd

		inst.Text = hex
		inst.Time = now

		if found, nT, nL := drivers.JNodeType(inst.ID[RF12BAND], inst.ID[RF12GRP], inst.ID[RF12NODEID], inst.Time); found {
			inst.Loc = nL
			s.client.Publish("rf12/"+nT, inst)
		} else {
			inst.Loc = "unknown"
			s.client.Publish("rf12/unknown", inst)
		}
	} else if conf := confRegex.FindStringSubmatch(text); conf != nil {
		inst.Node, _ = strconv.Atoi(conf[1]) //nodeid
		inst.Astx = conf[2] == "*"
		inst.ID[RF12GRP], _ = strconv.Atoi(conf[3]) //grp
		inst.ID[RF12BAND], _ = strconv.Atoi(conf[4]) //band
		//fmt.Println("groupID: ", inst.ID[RF12GRP])

		fmt.Println("Processing Config:", inst)
		s.Db[rwTopic] = inst

	} else if strings.HasPrefix(text, "[") && strings.Contains(text, "]") {
		fmt.Println("Sketch/Tag located:",text)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

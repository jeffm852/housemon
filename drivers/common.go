package drivers

import (
	"flag"
	"fmt"
	"log"
	"encoding/json"
	"os"
	"net/url"
	"strings"

	"github.com/jcw/jeebus"
)

var client *jeebus.Client

var murl *url.URL //all drivers are auto created via their own init, they must get this from main


func init() {  //if drivers (within this pkg) must init() then the mqtt param must be obtained this way!!
	var mqttAddr string

	//cannot redefine flag. state - use flagset
	flagset := flag.NewFlagSet("runflags", flag.ContinueOnError)

	flagset.StringVar(&mqttAddr, "mqtt", ":1883",
		"connect to MQTT server on <host><:port>")
	flagset.Parse(os.Args[1:])

	if !strings.Contains(mqttAddr, "://") {
		mqttAddr = "tcp://" + mqttAddr
	}
	var err error
	murl, err = url.Parse(mqttAddr)
	check(err)

	fmt.Println("Drivers using mqtt on:", murl)
}



func register(nT string, decoder jeebus.Service) {
	if client == nil {
		client = jeebus.NewClient(murl)
	}
	client.Register("rf12/"+nT+"/#", decoder)
}

func publish(nT string, v interface{}, m *jeebus.Message) {
	//TODO: implement reflection to replace marshalling to json and back
	//var b []byte
	b, err := json.Marshal(v)
	check(err)
	var im map[string]interface{}
	err = json.Unmarshal(b, &im)
	check(err)
	var vm = map[string]interface{} {"value": ""}

	for property, v := range im {
		vm["value"] = v
		val, err := json.Marshal(vm)
		check(err)
		topic := "/hm/" + m.Get("loc") + "/" + nT + "/" + property
		//topic += "/" + strconv.FormatInt(m.GetInt64("time"), 10)
		client.Publish(topic, val)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

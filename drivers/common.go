package drivers

import (
	"log"
	"encoding/json"

	"github.com/jcw/jeebus"
)

var client *jeebus.Client

func register(nT string, decoder jeebus.Service) {
	if client == nil {
		client = jeebus.NewClient(nil)
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

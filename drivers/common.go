package drivers

import (
	"log"
	"encoding/json"

	"github.com/jcw/jeebus"
)

func register(nT string, decoder jeebus.Service) {
	rf12Client := jeebus.NewClient("rf12")
	rf12Client.Register(nT+"/#", decoder)
}

func publish(nT string, v interface{}, m *jeebus.Message) {
	//TODO: implement reflection to replace marshalling to json and back
	//var b []byte
	b, err := json.Marshal(v)
	check(err)
	var f interface{}
	err = json.Unmarshal(b, &f)
	check(err)
	im := f.(map[string]interface{})
	var vm = map[string]interface{} {"value": ""}

	for property, v := range im {
		vm["value"] = v
		val, err := json.Marshal(vm)
		check(err)
		topic := "/hm/" + m.Get("loc") + "/" + nT + "/" + property
		//topic += "/" + strconv.FormatInt(m.GetInt64("time"), 10)
		jeebus.Publish(topic, val)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

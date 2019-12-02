package sorashell

import (
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	"strings"
	"time"
)

// naive cache which holds subscribers data for imsiFilterSuggestions
var devicesCache []prompt.Suggest

func (s *SoracomCompleter) deviceIdFilterSuggestions(word string) []prompt.Suggest {
	c := make(chan []prompt.Suggest, 1024)
	if len(devicesCache) == 0 {
		go getDevices(c, s.worker)
		select {
		case res := <-c:
			devicesCache = res
		case <-time.After(10 * time.Second):
			return []prompt.Suggest{{
				Text:        "Downloading device information in background",
				Description: "Hit space to see latest",
			}}
		}
	}
	return filterFunc(devicesCache, word, filterTextOrDescriptionFuzzy)
}

var getDevices = func(c chan<- []prompt.Suggest, worker *SoracomWorker) {
	var r []prompt.Suggest

	result := worker.Execute("devices list --fetch-all")
	if err := json.NewDecoder(strings.NewReader(result)).Decode(&devices); err != nil {
		c <- []prompt.Suggest{{
			Text:        "Error while running 'devices list --fetch-all'",
			Description: err.Error(),
		}}
	}
	for _, device := range devices {
		online := "offline"
		if device.Online {
			online = "online"
		}
		r = append(r, prompt.Suggest{
			Text: device.DeviceId,
			Description: fmt.Sprintf("%-14s | %-8s | %-8s | %-15s | %15s | %15s | %s",
				trunc(device.Endpoint, 14),
				device.Status,
				online,
				device.Imei,
				device.Imsi,
				device.IPAddress,
				device.Manufacturer,
			),
		})
	}
	c <- r
}

package sorashell

import (
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	"strings"
	"time"
)

// naive cache which holds sigfox devices data for sigfoxDeviceIDFilterSuggestions
var sigfoxDevicesCache []prompt.Suggest

func (s *SoracomCompleter) sigfoxDeviceIDFilterSuggestions(word string) []prompt.Suggest {
	c := make(chan []prompt.Suggest, 1024)
	if len(sigfoxDevicesCache) == 0 {
		go getSigfoxDevices(c, s.worker)
		select {
		case res := <-c:
			sigfoxDevicesCache = res
		case <-time.After(10 * time.Second):
			return []prompt.Suggest{{
				Text:        "Downloading device information in background",
				Description: "Hit space to see latest",
			}}
		}
	}
	return filterFunc(sigfoxDevicesCache, word, filterTextOrDescriptionFuzzy)
}

var getSigfoxDevices = func(c chan<- []prompt.Suggest, worker *SoracomWorker) {
	var r []prompt.Suggest

	result := worker.Execute("sigfox-devices list --fetch-all")
	if err := json.NewDecoder(strings.NewReader(result)).Decode(&sigfoxDevices); err != nil {
		c <- []prompt.Suggest{{
			Text:        "Error while running 'sigfox-devices list --fetch-all'",
			Description: err.Error(),
		}}
	}
	for _, device := range sigfoxDevices {
		r = append(r, prompt.Suggest{
			Text: device.DeviceID,
			Description: fmt.Sprintf("%-24s | %-8s | %-8s",
				trunc(device.Tags.Name, 24),
				device.Status,
				device.DeviceID,
			),
		})
	}
	c <- r
}

package sorashell

import (
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	"strings"
	"time"
)

// naive cache which holds subscribers data for imsiFilterSuggestions
var subscribersCache []prompt.Suggest

func (s *SoracomCompleter) imsiFilterSuggestions(word string) []prompt.Suggest {
	c := make(chan []prompt.Suggest, 1024)
	if len(subscribersCache) == 0 {
		go getSubscribers(c, s.worker)
		select {
		case res := <-c:
			subscribersCache = res
		case <-time.After(10 * time.Second):
			return []prompt.Suggest{{
				Text:        "Downloading IMSI in background",
				Description: "Hit space to see latest",
			}}
		}
	}
	return filterFunc(subscribersCache, word, filterTextOrDescriptionFuzzy)
}

var getSubscribers = func(c chan<- []prompt.Suggest, worker *SoracomWorker) {
	var r []prompt.Suggest

	result := worker.Execute("subscribers list --fetch-all")
	if err := json.NewDecoder(strings.NewReader(result)).Decode(&subscribers); err != nil {
		c <- []prompt.Suggest{{
			Text:        "Error while running 'subscribers list --fetch-all'",
			Description: err.Error(),
		}}
	}
	for _, subscriber := range subscribers {
		online := "offline"
		if subscriber.SessionStatus.Online {
			online = "online"
		}
		r = append(r, prompt.Suggest{
			Text: subscriber.Imsi,
			Description: fmt.Sprintf("%-12s | %-10s | %-7s | %-8s | %-11s | %s",
				trunc(subscriber.Subscription, 14),
				subscriber.Status,
				online,
				subscriber.ModuleType,
				trunc(subscriber.SpeedClass, 11),
				subscriber.Tags.Name,
			),
		})
	}
	c <- r
}

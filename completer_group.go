package sorashell

import (
	"encoding/json"
	"github.com/c-bata/go-prompt"
	"strings"
	"time"
)

// naive cache which holds group data for groupFilterSuggestions
var groupsCache []prompt.Suggest

func (s *SoracomCompleter) groupFilterSuggestions(word string) []prompt.Suggest {
	c := make(chan []prompt.Suggest, 1024)
	if len(groupsCache) == 0 {
		go getGroups(c, s.worker)
		select {
		case res := <-c:
			groupsCache = res
		case <-time.After(10 * time.Second):
			return []prompt.Suggest{{
				Text:        "Downloading groups in background",
				Description: "Hit space to see latest",
			}}
		}
	}
	return filterFunc(groupsCache, word, filterTextOrDescriptionFuzzy)
}

var getGroups = func(c chan<- []prompt.Suggest, worker *SoracomWorker) {
	var r []prompt.Suggest

	result := worker.Execute("groups list --fetch-all")
	if err := json.NewDecoder(strings.NewReader(result)).Decode(&groups); err != nil {
		c <- []prompt.Suggest{{
			Text:        "Error while running 'groups list'",
			Description: err.Error(),
		}}
	}
	for _, group := range groups {
		r = append(r, prompt.Suggest{
			Text:        group.GroupID,
			Description: group.Tags.Name,
		})
	}
	c <- r
}

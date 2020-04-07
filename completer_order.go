package sorashell

import (
	"encoding/json"
	"fmt"
	"github.com/c-bata/go-prompt"
	"strings"
	"time"
)

// naive cache which holds orders data for orderFilterSuggestions
var ordersCache []prompt.Suggest

func (s *SoracomCompleter) orderFilterSuggestions(word string) []prompt.Suggest {
	c := make(chan []prompt.Suggest, 1024)
	if len(ordersCache) == 0 {
		go getOrders(c, s.worker)
		select {
		case res := <-c:
			ordersCache = res
		case <-time.After(10 * time.Second):
			return []prompt.Suggest{{
				Text:        "Downloading orders in background",
				Description: "Hit space to see latest",
			}}
		}
	}
	return filterFunc(ordersCache, word, filterTextOrDescriptionFuzzy)
}

var getOrders = func(c chan<- []prompt.Suggest, worker *SoracomWorker) {
	var r []prompt.Suggest

	result := worker.Execute("orders list")
	if err := json.NewDecoder(strings.NewReader(result)).Decode(&orders); err != nil {
		c <- []prompt.Suggest{{
			Text:        "Error while running 'orders list'",
			Description: err.Error(),
		}}
	}
	for _, order := range orders.OrderList {
		products := make([]string, len(order.OrderItemList))
		for i, v := range order.OrderItemList {
			products[i] = fmt.Sprintf("%s (%d)", v.Product.ProductName, v.Product.Count)
		}
		r = append(r, prompt.Suggest{
			Text: order.OrderID,
			Description: fmt.Sprintf("%-10s | %s",
				order.OrderDateTime,
				trunc(strings.Join(products, " / "), 60),
			),
		})
	}
	c <- r
}

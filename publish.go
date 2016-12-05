package common

import (
	"net/http"

	"bytes"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
)

type Publisher struct {
	client *http.Client
	urls   []string
}

func NewPublisher(urls []string) Publisher {
	return Publisher{client: new(http.Client), urls: urls}
}

func (p Publisher) Publish(event interface{}) {
	json, err := json.Marshal(event)
	if err != nil {
		log.Error("Failed to convert event interface to json", err)
		return
	}
	reader := bytes.NewReader(json)
	for _, currUrl := range p.urls {
		if currUrl == "" {
			continue
		}
		request, err := http.NewRequest("POST", currUrl, reader)
		if err != nil {
			log.Errorf("Failed to create request. URL: %s Event: %v", currUrl, event)
			continue
		}
		request.Header.Add("Content-Type", "application/json")
		response, err := p.client.Do(request)
		if err != nil {
			log.Errorf("Failed to publish event. Error: %v URL:%s event: %v", err, currUrl, event)
			continue
		}
		log.Debugf("Event has been published. Response: %v URL: %s Event: %v", response, currUrl, event)
	}
}

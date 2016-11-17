package event

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
)

func TestPublishDockerEvent(t *testing.T) {
	const ID = "123qwe456"
	called := false
	ep := NewEventPublisher([]string{
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			assert.NoError(t, err)
			var event dockerclient.Event
			err = json.Unmarshal(body, &event)
			assert.NoError(t, err)
			assert.Equal(t, ID, event.ID)
			w.WriteHeader(http.StatusOK)
			called = true
		})).URL})
	ep.Publish(dockerclient.Event{ID: ID, Status: "Running", Type: "Container"})
	assert.True(t, called)
}

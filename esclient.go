package common

import (
	"errors"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
)

const illegalCharacters = "[:,?<>/\\*?| ]"

type ESClient struct {
	client *elastic.Client
	index  string
	typ    string
}

func NewESClient(url string) (*ESClient, error) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}

	return &ESClient{client, "", ""}, nil
}

func (esc *ESClient) Initialize(index, typ string) error {
	if strings.ContainsAny(index, illegalCharacters) {
		return errors.New(fmt.Sprintf("Index name contains illegal characters (%s)", index))
	}
	esc.index = index
	if strings.ContainsAny(typ, illegalCharacters) {
		return errors.New(fmt.Sprintf("Type contains illegal characters (%s)", typ))
	}
	esc.typ = typ

	return nil
}

func (esc *ESClient) Index(doc string) error {

	response, err := esc.client.Index().
		Index(esc.index).
		Type(esc.typ).
		BodyString(doc).
		Refresh("true").
		Do(context.Background())
	if err != nil {
		log.Errorf("Failed indexing document %s for index %s. %+v", doc, esc.index, err)
	} else {
		log.Infof("Indexed document ID: %s into index name: %s", response.Id, response.Index)
	}

	return err
}
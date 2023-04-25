package algolia

import (
	"encoding/json"
	"errors"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/opt"
	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/sirupsen/logrus"
)

var algoliaClient *search.Client

func Init(appID, apiSecret string) error {
	if algoliaClient != nil {
		return errors.New("algolia service already initiated")
	}
	algoliaClient = search.NewClient(appID, apiSecret)
	return nil
}

func InsertObject(indexName string, obj interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return err
	}
	m["objectID"] = m["id"]
	delete(m, "id")
	index := algoliaClient.InitIndex(indexName)
	if _, err := index.SaveObjects(m); err != nil {
		return err
	}
	return nil
}

func SearchObject(indexName, query string) ([]map[string]interface{}, error) {
	index := algoliaClient.InitIndex(indexName)
	results, err := index.Search(query, opt.AttributesToHighlight())
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}
	var m []map[string]interface{}
	results.UnmarshalHits(&m)
	for i := range m {
		m[i]["id"] = m[i]["objectID"]
		delete(m[i], "objectID")
	}
	return m, nil
}

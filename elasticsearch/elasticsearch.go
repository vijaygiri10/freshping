package elasticsearch

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
)

var (
	client *elastic.Client
)

func init() {
	var err error
	client, err = elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		fmt.Println("Error elastic NewClient : ", err)
	}
}

var Data string = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"doc":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
			"retweets":{
				"type":"long"
			},
				"tags":{
					"type":"keyword"
				},
				"location":{
					"type":"geo_point"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}
`

func InsertDataToElastic(Index, data string) {
	// Use the IndexExists service to check if a specified index exists.
	if exists, _ := client.IndexExists(Index).Do(context.Background()); !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(Index).BodyString(data).Do(context.Background())
		if err != nil {
			// Handle error
			fmt.Println("Error CreateIndex")
		}
		_ = createIndex.Acknowledged
	} else {
		if _, err := client.Index().Index(Index).Type("doc").BodyJson(data).Do(context.Background()); err != nil {
			fmt.Println("Error client.Index ")
		}
	}
}

package main

import (
	"fmt"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	// "context"
	// "time"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

//Must download one page and get a maximum of data to then save in the DB
// The commented code is work in progress that does not currently work. The db stay inaccessible for some reason.
// The download part works, but since I was trying to get acces to the DB to save the parts of the page, I commented
// this part as well.
func treatOneOffer(current_offer string) {
	fmt.Println("Link to download : " + current_offer)
	// response, err := http.Get(current_offer)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer response.Body.Close()

	// // Create a goquery document from the HTTP response
	// document, err := goquery.NewDocumentFromReader(response.Body)
	// if err != nil {
	// 	log.Fatal("Error loading HTTP response body. ", err)
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// defer cancel()
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI(
	// 	"mongodb://scarboni:<ghWG2m2aedqXrgvk>@scrap-test-btzyh.mongodb.net/test?retryWrites=true&w=majority",
	// ))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = client.Ping(ctx, readpref.Primary())
	// if err != nil {
	// 	fmt.Println("errrr")
	// 	log.Fatal(err)
	// }
	// defer client.Disconnect(ctx)
}

//must treat one standalone url from kafka topic at a time to extract and save datas into mongoDB
func main() {

	//Kafka use
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"web-adresses", "^aRegex.*[Tt]opic"}, nil)
	defer c.Close()

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			treatOneOffer(string(msg.Value))
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

}

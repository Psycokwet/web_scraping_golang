package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type Auto struct {
	Url    string
	Fields []Tuple
}

type Tuple struct {
	Field   string
	Content string
}

func initialize_auto(url string, fields_nbr int) Auto {
	return Auto{url, make([]Tuple, fields_nbr)}
}

//Must download one page and get a maximum of data to then save in the DB
// The commented code is work in progress that does not currently work. The db stay inaccessible for some reason.
// The download part works, but since I was trying to get acces to the DB to save the parts of the page, I commented
// this part as well.
func treatOneOffer(current_offer string, collection *mongo.Collection) {
	response, err := http.Get(current_offer)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}
	current_auto := initialize_auto(current_offer, 0)
	document.Find("div").Each(func(index int, element *goquery.Selection) {
		element_class, exists := element.Attr("class")
		if exists && strings.HasPrefix(element_class, "small-12 bg-box landing  columns") {
			// Produce messages to topic (asynchronously)
			// fill up current_auto with any fields found below this div, since it does englobe the whole offer.
			// The rest of the page is useless
			// WIP, simulating one data found for now
			current_auto.Fields = append(current_auto.Fields, Tuple{"fieldname", "fieldcontent"})
		}
	})

	insertResult, err := collection.InsertOne(context.TODO(), current_auto)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

//must treat one standalone url from kafka topic at a time to extract and save datas into mongoDB
func main() {

	//starting mongodb connection
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("scrap_db").Collection("test")

	defer client.Disconnect(ctx)

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
			treatOneOffer(string(msg.Value), collection)
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

}

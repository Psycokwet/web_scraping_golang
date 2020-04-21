package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func sendAllUrlsFromListingUrl(current_link string, topic string, p *kafka.Producer) string {
	response, err := http.Get(current_link)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Find and send offer URLs
	//must put in a black list http://www.autoreflex.com/voiture-occasion-professionnels.htm
	document.Find("a").Each(func(index int, element *goquery.Selection) {
		linkoffer, exists := element.Attr("href")
		if exists && strings.HasPrefix(linkoffer, "voiture") {
			// Produce messages to topic (asynchronously)
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte("http://www.autoreflex.com/" + linkoffer),
			}, nil)
		}
	})
	// Find the link to the next listing page to process
	var next_link string
	fmt.Println("Hey [" + next_link + "]")
	document.Find("a").Each(func(index int, element *goquery.Selection) {
		title, exists := element.Attr("title")
		if exists && (title == "Suivant") {
			link, exists := element.Attr("href")
			if exists {
				next_link = "http://www.autoreflex.com/" + link
			}
		}
	})
	return next_link
}

// Must get all standalone page's offer from a listing page and send them through kafka topic, one at a time
func main() {

	// create kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	topic := "web-adresses"

	// start with the first adresse of the listing, and keep on with the next ones
	current_link := "http://www.autoreflex.com/137.0.-1.-1.-1.0.999999.1900.999999.-1.99.0.1?fulltext=&amp;geoban=M137R99"
	for current_link = sendAllUrlsFromListingUrl(current_link, topic, p); current_link != ""; current_link = sendAllUrlsFromListingUrl(current_link, topic, p) {

	}
	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}

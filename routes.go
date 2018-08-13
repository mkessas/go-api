package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/buger/jsonparser"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
)

var creds = Creds{username: "admin", password: "admin"}

var paths = []Path{
	{path: "/api/test1", handler: test1Handler, method: "GET", protected: true},
	{path: "/api/test2", handler: test2Handler, method: "GET", protected: false},
	{path: "/api/pump", handler: pump, method: "GET", protected: false},
	{path: "/api/do", handler: doHandler, method: "POST", protected: false},
	{path: "/api/kafkaProducer", handler: kafkaProducerHandler, method: "GET", protected: false},
	{path: "/api/hello/{name}", handler: helloHandler, method: "GET", protected: false},
	{path: "/api/temperature", handler: temperatureHandler, method: "GET", protected: false},
}

func test(body string, topic string) {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "127.0.0.1:9092"})
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

	// Produce messages to topic (asynchronously)
	//topic := "myTopic"
	// for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client"} {
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(body),
	}, nil)
	// }

	// Wait for message deliveries
	p.Flush(15 * 1000)
}
func doHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := r.FormValue("url")
	resp, err := http.Get(url)
	//------------------------------
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%s", err)
		return
	}

	defer r.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		test(bodyString, "getResponces")
	}

	data, _ := ioutil.ReadAll(resp.Body)
	title, _, _, err := jsonparser.Get(data, "title", "")
	fmt.Fprintf(w, ":%s:", title)

}
func test1Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test 1 handler!")
}
func kafkaProducerHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "kicking of kafka producer ")
	test("test text", "simple_test")
}
func test2Handler(w http.ResponseWriter, r *http.Request) {
	//w.WriteHeader(413)
	w.Header().Set("x-Yes", "Cool")
	fmt.Fprintf(w, "Test 2 handler!\n")

	for k, v := range r.Header {
		fmt.Fprintf(w, "- %s : %s\n", k, v)

	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Hello %s!", vars["name"])
}

func pump(w http.ResponseWriter, r *http.Request) {
	configfileName := "job.cfg"
	if len(os.Args) > 1 {
		configfileName = os.Args[1]
	}
	fmt.Fprintf(w, "Using configuration file:%s\n", configfileName)
	inFile, _ := os.Open(configfileName)
	ch := make(chan string)
	var urls []string

	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		url := scanner.Text()
		urls = append(urls, url)
	}

	for _, aaa := range urls {
		fmt.Fprintf(w, "\nattemting to GET %s \n", aaa)
		fmt.Fprintf(w, "%s\n", aaa)
		go MakeRequest(aaa, ch)
	}
	for range urls[1:] {
		fmt.Fprintf(w, "%s", <-ch)
	}
}

func MakeRequest(url string, ch chan<- string) {
	start := time.Now()
	resp, _ := http.Get(url)

	secs := time.Since(start).Seconds()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	responseString := string(body)
	fmt.Println(responseString)
	ch <- fmt.Sprintf("%.2f elapsed with response length: %d %s\n", secs, len(body), url)
}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get("http://alarm:3000/api/sensor/status")

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%s", err)
		return
	}

	defer r.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	temperature, _, _, err := jsonparser.Get(data, "details", "status", "temperature")

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%s", err)
	}

	fmt.Fprintf(w, "%s", temperature)

}

package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"example.com/greetings"
)

func main() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	rand.Seed(time.Now().UnixNano())

	names := []string{"Gladys", "Samantha", "Darrin"}
	messages, err := greetings.Hellos(names)
	if err != nil {
		log.Fatal(err)
	}
	var allNodesWaitGroup sync.WaitGroup
	for _, message := range messages {
		allNodesWaitGroup.Add(1)
		go func(message string) {
			defer allNodesWaitGroup.Done()
			sleep := rand.Intn(1000) + 500
			time.Sleep(time.Duration(sleep) * time.Millisecond)
			fmt.Println(message, sleep)
		}(message)
	}
	allNodesWaitGroup.Wait()
}

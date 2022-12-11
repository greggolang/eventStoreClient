package main

import (
	"fmt"
	"log"
	"time"

	dummy_client "github.com/greggolang/eventStoreClient/dummy/client"
)

// An example of how the client clould be integrated into the other project
// (such as Crospsoint or HomeServer).
func main() {
	// Parse the eventeur connection configuration.
	// TODO: define configuration.
	config := dummy_client.Config{
		Url: "localhost:8080",
	}

	// Initialize eventeur client with the provided configuration.
	eventsClient := dummy_client.NewDummyClient(config)
	go eventsClient.Connect()
	<-time.After(time.Second)
	log.Print("server started")

	// We are only printing the received messages here, but any logic
	// can be defined in the `onEvent` callback.
	teslaSub, err := eventsClient.Subscribe("tesla-events", onEvent)
	if err != nil {
		eventsClient.Publish("example-client-events", []byte("Failed to subscribe to tesla events"))
	}
	log.Print("sub to tesla events")

	smartthingsSub, err := eventsClient.Subscribe("smartthings-events", onEvent)
	if err != nil {
		eventsClient.Publish("example-client-events", []byte("Failed to subscribe to smartthings events"))
	}
	log.Print("sub to smartthings events")

	// Lets inform every client that is listening for `example-client-events` that
	// the client was started succesfully.
	eventsClient.Publish("example-client-events", []byte("Client succesfully initialized"))
	log.Print("publish to client events")

	// For this example we are unsubscribing instantly.
	teslaSub.Unsubscribe()
	smartthingsSub.Unsubscribe()
	log.Print("unsub to tesla and smartthings")

	// After we are done, the client `Close()` method informs the server, that
	// it won't listen to any events.
	eventsClient.Close()
}

func onEvent(clientID uint64, data []byte) {
	// Eventeur passes data without any assumptions about it.
	// The user decides what serialization method to use.
	// In this case MessagePack is used to minimize the size of the meta data.
	// TODO: use MessagePack.
	str := string(data)
	fmt.Printf("Got message from client [%d] with data: %s\n", clientID, str)
}

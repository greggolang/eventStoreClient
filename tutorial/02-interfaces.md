### Interfaces

## A quick example of how to define interface in go

In Go, an interface is a set of method signatures that defines a contract for types that implement the interface. In other words, an interface defines a set of behaviors that a type must implement in order to be considered a member of the interface. A type can be used interchangeably with other types that implement the same interface.

```
// Define a new interface type called "Shape"
// with a single method called "Area() float64"
type Shape interface {
    Area() float64
}

// Define a new struct type called "Rectangle"
// with fields for length and width
type Rectangle struct {
    length float64
    width  float64
}

// Implement the "Area() float64" method on the "Rectangle" type
// by returning the product of the length and width
func (r Rectangle) Area() float64 {
    return r.length * r.width
}

// Create a new "Rectangle" value with a length of 5 and a width of 10
rect := Rectangle{length: 5, width: 10}

// Create a new variable of the "Shape" interface type
// and assign the "rect" value to it
var s Shape = rect

// Use the "Area() float64" method on the "Shape" interface variable
fmt.Println(s.Area())  // Output: 50
```

In this example, the Shape interface defines a contract for any type that has an Area() float64 method. The Rectangle type satisfies this contract by implementing the Area() float64 method. This allows us to use the Shape interface as a type, and assign values of the Rectangle type to it. We can then use the Area() float64 method on the interface variable, which will call the method on the underlying Rectangle value.

## Interfaces in Eventeur

The code in package called eventeur defines several types and interfaces. The EventeurEventHandler type is a function type that takes in a clientID of type uint64 and a data byte slice, and returns no value.

The Subscription interface defines a single method named Unsubscribe, which returns an error. This method allows a type that implements the Subscription interface to unsubscribe from an event. We will be using the same exact interface to implement a multiple clients that can connect to different event sourcing backends.

The interface defined in `eventeur-server/client.go` only defines the methods that the user of a library should care about. The user should not worry about how to convert event data to some underlying data structure that is used between client and server. The implementation such as `eventeur-server/dummy` takes care of that.

This is an interface that various clients will need to implement:
```
package eventeur

type EventeurEventHandler func(clientID uint64, data []byte)

type Subscription interface {
	Unsubscribe() error
}

type EventeurClient interface {
	Connect() 
	Publish(topic string, data []byte) error
	Subscribe(topic string, callback EventeurEventHandler) (Subscription, error)
	Unsubscribe(topic string) error
	Close()
}
```

The EventeurClient interface defines several methods:
* Connect: This method establishes a connection to an event service. It takes no arguments and returns no value.
* Publish: This method publishes an event to a specified topic. It takes a topic string and a data byte slice as arguments, and returns an error if the event could not be published.
* Subscribe: This method subscribes to events on a specified topic. It takes a topic string, a callback function of type EventeurEventHandler, and returns a Subscription interface and an error if the subscription could not be established.
* Unsubscribe: This method unsubscribes from events on a specified topic. It takes a topic string as an argument and returns an error if the unsubscription could not be completed.
* Close: This method closes the connection to the event service. It takes no arguments and returns no value.

Here's an example of how you might use these types and interfaces:

```
package main

import (
	"fmt"
	"eventeur"
)

func main() {
	// Create a new instance of the EventeurClient interface.
	// This will be a type that implements the EventeurClient interface.
	client := newEventeurClient()

	// Connect to the event service.
	client.Connect()

	// Subscribe to events on the "golang" topic.
	subscription, err := client.Subscribe("golang", handleEvent)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Publish an event to the "golang" topic.
	err = client.Publish("golang", []byte("hello world"))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Unsubscribe from the "golang" topic.
	err = subscription.Unsubscribe()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Close the connection to the event service.
	client.Close()
}

// handleEvent is the event handler that will be called whenever an event is
// received on the "golang" topic.
func handleEvent(clientID uint64, data []byte) {
	fmt.Printf("Received event from client %d: %s\n", clientID, string(data))
}
```

By using `newEventeurClient()` we as a user don't care how the code underneath works. It can connect to our `Dummy` server, or it can connect to `Kafka` or `EventStore` backends.

In this example, we create a new instance of the EventeurClient interface. This will be a type that implements the EventeurClient interface. We then use the Connect method to establish a connection to the event service.

Next, we use the Subscribe method to subscribe to events on the "golang" topic. This method takes a topic string and an EventeurEventHandler function as arguments. The EventeurEventHandler function will be called whenever an event is received on the "golang" topic.

After subscribing to the "golang" topic, we use the Publish method to publish an event to the same topic. This sends the "hello world" message to the event service, which will be received by our subscription.

We then use the Unsubscribe method on the Subscription interface to unsubscribe from the "golang" topic. Finally, we use the Close method to close the connection to the event service.

This is just one example of how you can use the eventeur package and its associated types and interfaces to work with an event service in Go. There are many other ways you can use these types and interfaces to build event-driven applications in Go.


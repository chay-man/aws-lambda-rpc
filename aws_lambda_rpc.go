package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda/messages"
)

var (
	host    = flag.String("h", "localhost", "The ip the lambda is running on. This must match the environment variable AWS_LAMBDA_RUNTIME_API")
	port    = flag.Int("p", 8080, "The port the lambda is listening to. This must match the environment variable _LAMBDA_SERVER_PORT")
	timeout = flag.Int64("t", 86400, "The time to wait in seconds until sending the done signal to the lambda. Default is 1 day")
	file    = flag.String("f", "", "The file to send in the InvokeRequest. Must be valid json")
	debug   = flag.Bool("debug", false, "Enables debug logging")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	if *file == "" {
		log.Fatal("file is required to invoke the lambda")
	}

	bytes, err := os.ReadFile(*file)
	if err != nil {
		log.Fatalf("error opening file: %s", err.Error())
	}
	errInv := Invoke(bytes)
	if errInv != nil {
		log.Fatal(errInv)
	}

	if *debug {
		log.Println("no errors while invoking lambda")
	}
}

func Invoke(bytes []byte) error {
	address := fmt.Sprintf("%s:%d", *host, *port)
	if *debug {
		log.Printf("Using address: %s\n", address)
	}

	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return err
	}

	req := messages.InvokeRequest{
		Payload: bytes,
		Deadline: messages.InvokeRequest_Timestamp{
			Seconds: time.Now().Unix() + *timeout,
		},
	}

	var res messages.InvokeResponse
	err = client.Call("Function.Invoke", req, &res)
	if err != nil {
		return fmt.Errorf("failed calling Funcion.Invoke: %v", err)
	}
	if res.Error != nil {
		return fmt.Errorf("failure in invoking lambda: %v", res.Error)
	}

	if *debug {
		log.Println("successfully invoked lambda")
	}

	var payload map[string]any
	if err = json.Unmarshal(res.Payload, &payload); err != nil {
		return fmt.Errorf("failure unmarshaling response: %v", err)
	}

	if *debug {
		log.Println("printing response...")
		for key, value := range payload {
			log.Printf("key=%v value=%v", key, value)
		}
	}

	return nil
}

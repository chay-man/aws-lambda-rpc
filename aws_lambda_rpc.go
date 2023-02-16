package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda/messages"
	"log"
	"net/rpc"
	"os"
)

var (
	host    = flag.String("h", "localhost", "The ip the lambda is running on. This must match the environment variable AWS_LAMBDA_RUNTIME_API")
	port    = flag.Int("p", 8080, "The port the lambda is listening to. This must match the environment variable _LAMBDA_SERVER_PORT")
	timeout = flag.Int64("t", 300, "The time to wait until sending the done signal to the context")
	data    = flag.String("f", "", "The data to send in the InvokeRequest. Must be valid json")
)

func main() {
	flag.Parse()

	if *data == "" {
		log.Fatal("data is required to invoke the lambda")
	}

	bytes, err := os.ReadFile(*data)
	if err != nil {
		log.Fatalf("error opening file: %s", err.Error())
	}
	errInv := Invoke(bytes)
	if errInv != nil {
		log.Fatalf("failed to invoke: %s", errInv.Error())
	}
}

func Invoke(bytes []byte) error {
	address := fmt.Sprintf("%s:%d", *host, *port)
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return err
	}

	req := messages.InvokeRequest{
		Payload: bytes,
		Deadline: messages.InvokeRequest_Timestamp{
			Seconds: *timeout,
		},
	}

	var res messages.InvokeResponse
	return client.Call("Function.Invoke", req, &res)
}

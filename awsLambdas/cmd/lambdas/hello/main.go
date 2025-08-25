package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

// TODO: uses the event apropriate in place of interface if necessary, such as:
// events.APIGatewayProxyRequest, events.CloudWatchEvent, events.SQSEvent, events.KafkaEvent
func Handler(ctx context.Context, ievent interface{}) {
	fmt.Println("Hello world golang in localstack")
}

func main() {
	lambda.Start(Handler)
}

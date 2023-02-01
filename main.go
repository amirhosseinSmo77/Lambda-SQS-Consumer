package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Transaction struct {
	Id     int    `json:"id"`
	Amount int    `json:"amount"`
	State  string `json:"state"`
}

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	/* Create aws config, in our case we should connect to local aws */
	config := aws.Config{Endpoint: aws.String("http://localhost:4566")}
	/* creates DynamoDB session */
	sess, err := session.NewSession(&config)
	if err != nil {
		fmt.Errorf("failed to create session: %v", err)
		return err
	}

	ddb := dynamodb.New(sess)

	/* Process SQS events to store our data in DynamoDG */
	for _, record := range sqsEvent.Records {
		body := record.Body

		/* Parse string body as Json and map it our Transaction model */
		var transaction Transaction
		err := json.Unmarshal([]byte(body), &transaction)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return err
		}

		/* current timestamp nanoSecond as string */
		currentTime := strconv.FormatInt(time.Now().UnixNano(), 10)

		/* Create the object we are going to store into the DynamoDB */
		input := &dynamodb.PutItemInput{
			TableName: aws.String("PaymentsStates"),
			Item: map[string]*dynamodb.AttributeValue{
				"id": {
					N: aws.String(strconv.Itoa(transaction.Id)),
				},
				"timestamp": {
					S: aws.String(currentTime),
				},
				"amount": {
					N: aws.String(strconv.Itoa(transaction.Amount)),
				},
				"state": {
					S: aws.String(transaction.State),
				},
			},
		}

		_, err = ddb.PutItem(input)
		if err != nil {
			fmt.Errorf("failed to put item into DynamoDB: %v", err)
			return err
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}

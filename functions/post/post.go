package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty"`
}

var (
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorCouldNotMarshalItem     = "could not marshal item"
	ErrorInvalidDeviceData       = "invalid device data! some of the device information is not filled properly"
	ErrorCouldNotDynamoPutItem   = "could not dynamo put item"
)

type Device struct {
	Id          string `json:"id"`
	DeviceModel string `json:"deviceModel"`
	Name        string `json:"name"`
	Note        string `json:"note"`
	Serial      string `json:"serial"`
}

func createDevice(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*Device,
	error,
) {
	var d Device

	if err := json.Unmarshal([]byte(req.Body), &d); err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	if len(d.Id) < 1 || len(d.Name) < 1 || len(d.DeviceModel) < 1 || len(d.Note) < 1 || len(d.Serial) < 1 {

		return nil, errors.New(ErrorInvalidDeviceData)
	}

	av, err := dynamodbattribute.MarshalMap(d)

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &d, nil
}
func apiResponse(status int, body interface{}) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/json"}}
	resp.StatusCode = status

	stringBody, _ := json.Marshal(body)
	resp.Body = string(stringBody)
	return &resp, nil
}

const tableName = "DevicesList"

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	result, err := createDevice(req, tableName, dynaClient)

	if err != nil {
		return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
	}

	return apiResponse(http.StatusCreated, result)
}

func init() {
	region := os.Getenv("AWS_REGION")
	fmt.Println("region= ", region)
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)})

	if err != nil {
		return
	}
	dynaClient = dynamodb.New(awsSession)
}

func main() {
	lambda.Start(handler)
}

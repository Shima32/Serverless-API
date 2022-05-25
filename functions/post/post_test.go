package main

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
)

type stubDynamoDB struct {
	dynamodbiface.DynamoDBAPI
}

func (m *stubDynamoDB) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {

	Id := dynamodb.AttributeValue{}
	Id.SetS("id1")
	DeviceModel := dynamodb.AttributeValue{}
	DeviceModel.SetS("deviceModel1")
	Name := dynamodb.AttributeValue{}
	Name.SetS("sensor")
	Note := dynamodb.AttributeValue{}
	Note.SetS("testing")
	Serial := dynamodb.AttributeValue{}
	Serial.SetS("000A251")

	resp := make(map[string]*dynamodb.AttributeValue)
	resp["id"] = &Id
	resp["deviceModel"] = &DeviceModel
	resp["name"] = &Name
	resp["note"] = &Note
	resp["serial"] = &Serial

	return nil, nil
}

func TestCreateDevice(t *testing.T) {
	svc := &stubDynamoDB{}
	//var d_test Device
	testsCreate := []struct {
		name    string
		request events.APIGatewayProxyRequest
		expect  *Device
		err     error
	}{
		{
			name:    "test_post1",
			request: events.APIGatewayProxyRequest{Body: `{"id":"id242", "deviceModel":"M242",  "Name" : "sen2",  "Note": "testing", "serial" :"0000012"}`},
			expect: &Device{
				Id:          "id242",
				DeviceModel: "M242",
				Name:        "sen2",
				Note:        "testing",
				Serial:      "0000012",
			},
			err: nil,
		},
		{
			name:    "test_post2",
			request: events.APIGatewayProxyRequest{Body: `{"id":"", "deviceModel":"M22",  "Name" : "sen22",  "Note": "testing", "serial" :"0000012"}`},
			expect:  nil,
			err:     errors.New(ErrorInvalidDeviceData),
		},
		{
			name:    "test_post3",
			request: events.APIGatewayProxyRequest{Body: `{"id""id3", "deviceModel":"M23",  "Name" : "sen22",  "Note": "testing", "serial" :"0000012"}`},
			expect:  nil,
			err:     errors.New(ErrorFailedToUnmarshalRecord),
		},
	}

	for _, test := range testsCreate {
		d_test, err := createDevice(test.request, "DevicesList", svc)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, d_test)
	}

}

func TestHandler(t *testing.T) {
	dynaClient = &stubDynamoDB{}
	var d_test Device
	testsHand := []struct {
		name    string
		request events.APIGatewayProxyRequest
		expect  Device
		err     error
	}{
		{
			name:    "test_post1",
			request: events.APIGatewayProxyRequest{Body: `{"id":"id2", "deviceModel":"M2",  "Name" : "sen2",  "Note": "testing", "serial" :"0000012"}`},
			expect: Device{
				Id:          "id2",
				DeviceModel: "M2",
				Name:        "sen2",
				Note:        "testing",
				Serial:      "0000012",
			},
			err: nil,
		},
	}

	for _, test := range testsHand {
		resp, err := handler(test.request)
		json.Unmarshal([]byte(resp.Body), &d_test)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, d_test)
	}

}

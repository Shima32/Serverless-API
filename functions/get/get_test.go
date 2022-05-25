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

func (m *stubDynamoDB) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {

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

	reqId := input.Key
	var reqKey string
	for _, db := range reqId {
		reqKey = db.GoString()
	}

	if reqKey != Id.GoString() {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	output := &dynamodb.GetItemOutput{
		Item: resp,
	}

	return output, nil
}

func TestFetchDevice(t *testing.T) {
	testsFetch := []struct {
		name          string
		id            string
		response      *Device
		expectedError error
	}{
		{
			name: "test_Table1",
			id:   "id1",
			response: &Device{
				Id:          "id1",
				DeviceModel: "deviceModel1",
				Name:        "sensor",
				Note:        "testing",
				Serial:      "000A251",
			},
			expectedError: nil,
		},
		{
			name:          "test_Table2",
			id:            "id2",
			response:      nil,
			expectedError: errors.New(ErrorFailedToFetchRecord),
		},
		{
			name:          "test_Table3",
			id:            "",
			response:      nil,
			expectedError: errors.New(ErrorDeviceDoesNotExist),
		},
	}
	svc := &stubDynamoDB{}

	for _, test := range testsFetch {
		d_test, err := fetchDevice(test.id, "DevicesList", svc)
		assert.IsType(t, test.expectedError, err)
		assert.Equal(t, test.response, d_test)
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
			name:    "test_get1",
			request: events.APIGatewayProxyRequest{Body: `{"id":"id1"}`},
			expect: Device{
				Id:          "id1",
				DeviceModel: "deviceModel1",
				Name:        "sensor",
				Note:        "testing",
				Serial:      "000A251",
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

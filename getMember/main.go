package main

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Member struct {
	Id       string `json:"Id" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Email    int    `json:"email" validate:"required"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//dynamoDB session connect
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := dynamodb.New(sess)

	id := request.PathParameters["id"]

	// Create GetItemInput instance
	params := &dynamodb.GetItemInput{
		TableName: aws.String("member-crud-dev"),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
	}

	// Execute GetItem
	result, err := svc.GetItem(params)
	if err != nil {
		log.Warn(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       string(err.Error()),
			StatusCode: 500,
		}, nil
	}
	member := Member{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &member); err != nil {
		log.WithError(err).Warnf("Request body unmarshal error.")
		return events.APIGatewayProxyResponse{
			Body:       string("This request is invalid"),
			StatusCode: 400,
		}, nil
	}
	body, _ := json.Marshal(member)
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}

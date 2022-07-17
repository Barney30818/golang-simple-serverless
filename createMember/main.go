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
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Member struct {
	Id       string `json:"Id" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

type MailInfo struct {
	Id    string `json:"Id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//load credentials from the shared credentials file, ~/.aws/credentials and the default region from ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	//dynamoDB connect
	svc := dynamodb.New(sess)

	log.Info("DB connect successfully")

	// parse request body
	var requestBody Member
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		log.WithError(err).Warnf("Request body unmarshal error.")
		return events.APIGatewayProxyResponse{
			Body:       string("This request is invalid"),
			StatusCode: 400,
		}, nil
	}

	memberInfo := Member{
		Id:       requestBody.Id,
		Password: requestBody.Password,
		Name:     requestBody.Name,
		Email:    requestBody.Email,
	}

	av, err := dynamodbattribute.MarshalMap(memberInfo)
	if err != nil {
		log.WithError(err).Warnf("Got error while marshalling")
		return events.APIGatewayProxyResponse{
			Body:       string("This request is invalid"),
			StatusCode: 400,
		}, nil
	}

	// Create PutItemInput instance
	params := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("member-crud-dev"),
	}

	// Execute PutItem
	result, err := svc.PutItem(params)
	if err != nil {
		log.Warn(err.Error())
		return events.APIGatewayProxyResponse{
			Body:       string(err.Error()),
			StatusCode: 500,
		}, nil
	}

	log.Infof("DB insert successfully,the result is: %s", result)

	//send to SQS
	mailInfo := MailInfo{
		Id:    memberInfo.Id,
		Name:  memberInfo.Name,
		Email: memberInfo.Email,
	}
	var sqsURL = "https://sqs.us-west-2.amazonaws.com/801659726931/SendEmailQueue"
	log.Infof("sqs url is: %s", sqsURL)
	err = SendMessage(sess, sqsURL, mailInfo)
	if err != nil {
		log.Warnf("Got an error while trying to send message to queue: %v", err)
		return events.APIGatewayProxyResponse{
			Body:       string("fail to send SQS"),
			StatusCode: 400,
		}, nil
	}
	accountId, IdOk := result.Attributes["Id"]
	if IdOk {
		log.Info("Insert data successfully, accountId: %s", accountId)
	}
	body, _ := json.Marshal(memberInfo)
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func SendMessage(sess *session.Session, queueUrl string, messageBody MailInfo) error {
	sqsClient := sqs.New(sess)

	body, _ := json.Marshal(messageBody)
	_, err := sqsClient.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: aws.String(string(body)),
	})

	return err
}

func main() {
	lambda.Start(Handler)
}

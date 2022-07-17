package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	gomail "gopkg.in/gomail.v2"
)

type MailInfo struct {
	Id    string `json:"Id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {

	var requestBody MailInfo
	err := json.Unmarshal([]byte(sqsEvent.Records[0].Body), &requestBody)
	if err != nil {
		log.Warnf("Request body unmarshal error.", err)
	}
	log.Infof("Id is: %s,name is: %s, email is: %s", requestBody.Id, requestBody.Name, requestBody.Email)
	//send gmail
	m := gomail.NewMessage()
	m.SetHeader("From", "barney30818@gmail.com")
	m.SetHeader("To", requestBody.Email)
	m.SetHeader("Subject", "Hello "+requestBody.Id)
	m.SetBody("text/html", "Dear"+requestBody.Name+",\n很高興您已經加入會員\n祝您有個美好時光")

	d := gomail.NewDialer("smtp.gmail.com", 587, "barney30818@gmail.com", "panrvrwfhpjjpybg")

	// Send the email to member
	if err := d.DialAndSend(m); err != nil {
		log.Warn(err)
		panic(err)
	}

	log.Info("send mail successfully")

	return nil
}

func main() {
	lambda.Start(Handler)
}

package main

import (
	"bytes"
	"context"
	"encoding/json"

	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/joho/godotenv/autoload"
	"github.com/line/line-bot-sdk-go/linebot"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

var (
	client               *linebot.Client
	err                  error
	CHANNEL_SECRET       = "f2aeb1487d679a1655fb230bb516db94"
	CHANNEL_ACCESS_TOKEN = "1fOCBx1XiL91H2gvdMHzrA5kPZ6R31A7+1O8xh071vvKDUL//mD2wUSaDy9lCMfgZWzvThILfmV3diqqnMwvIMywXWg1U//PdLdLItwgt2w5ZNsT0uWBKBE//A3l+4z9AcsdKFNq+iYlvvZNnjtipwdB04t89/1O/w1cDnyilFU="
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {

	//complement line robot

	client, err = linebot.New(CHANNEL_SECRET, CHANNEL_ACCESS_TOKEN)

	log.Println("success")
	if err != nil {
		log.Println(err.Error())
	}

	//http.HandleFunc("/callback", callbackHandler)

	//log.Fatal(http.ListenAndServe(":84", nil))

	//api response
	var buf bytes.Buffer

	body, err := json.Marshal(map[string]interface{}{
		"message": "Okay so your other function also executed successfully!",
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "world-handler",
		},
	}

	return resp, nil
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	// 接收請求
	events, err := client.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}

		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				// 回覆訊息
				if _, err = client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Println(err.Error())
				}
			}
		}
	}
}

func main() {
	lambda.Start(Handler)
}

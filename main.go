package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"slack-bot/bot"
	"slack-bot/routes"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda
var r *gin.Engine

func init() {
	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Gin cold start")
	r = gin.Default()

	r.POST("/slack-event", routes.SlackEventHandler)

	ginLambda = ginadapter.New(r)
}

func apiGatewayProxyRequestHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	_, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func PollEventHandler(ctx context.Context, req bot.PollEvent) (bot.PollEvent, error) {
	_, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	return bot.PollEvent{
		Name:   req.Name,
		Status: "Success",
	}, nil
}

type CustomHandle struct{}

func (handler CustomHandle) Invoke(ctx context.Context, payload []byte) ([]byte, error) {

	fmt.Println("RAW REQUEST")
	fmt.Println(string(payload[:]))

	var apiGatewayProxyRequest events.APIGatewayProxyRequest
	json.Unmarshal(payload, &apiGatewayProxyRequest)

	// Check if event was actually an API Gateway Event by checking if HTTPMethod exists
	if apiGatewayProxyRequest.HTTPMethod != "" {
		fmt.Println("Received an API Gateway Proxy Event.")
		response, err := apiGatewayProxyRequestHandler(ctx, apiGatewayProxyRequest)
		if err != nil {
			return nil, err
		}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}
		return responseBytes, nil
	}

	var pollEvent bot.PollEvent
	json.Unmarshal(payload, &pollEvent)

	if PpllEvent.Name != "" {
		fmt.Println("Received a Poll Request Event.")

		response, err := PollEventHandler(ctx, pollEvent)
		if err != nil {
			return nil, err
		}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}
		return responseBytes, nil
	}

	return []byte(`{"error":"Unable to handle this event type."}`), nil
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") != "" {
		handle := new(CustomHandle)
		lambda.StartHandler(handle)
	} else {
		r.Run()
	}
}

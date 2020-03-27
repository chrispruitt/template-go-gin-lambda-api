package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

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

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	json, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Not an API gateway event.")
		log.Fatal(err)
	}

	fmt.Println(req.Body)
	fmt.Println(req)
	fmt.Println(json)

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") != "" {
		lambda.Start(Handler)
	} else {
		r.Run()
	}
}

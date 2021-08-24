package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/teris-io/shortid"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	LinksTableName = "UrlShortenerLinks"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	ShortURL string `json:"short_url"`
}

// Link describes a DynamoDB table item - ShortURL Is used as the primary key
type Link struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	req := Request{}
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{
			Body: "encountered an error when parsing the URL request param",
		}, err
	}

	cfg := aws.Config{
		Region: aws.String(endpoints.EuWest2RegionID),
	}

	sess := session.Must(session.NewSession(&cfg))

	service := dynamodb.New(sess)

	shortURL := shortid.MustGenerate()
	for shortURL == "shorten" {
		shortURL = shortid.MustGenerate()
	}

	link := &Link{
		ShortURL: shortURL,
		LongURL:  req.URL,
	}

	av, err := dynamodbattribute.MarshalMap(link)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "encountered an error when marshall via aws",
		}, err
	}

	input := dynamodb.PutItemInput{Item: av, TableName: aws.String(LinksTableName)}

	if _, err := service.PutItem(&input); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	r := Response{ShortURL: shortURL}

	response, err := json.Marshal(r)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil
}

func main() {
	lambda.Start(Handler)
}

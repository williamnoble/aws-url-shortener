package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"net/http"
)

const (
	LinksTableName = "UrlShortenerLinks"
)

type Link struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	shortURL, _ := request.PathParameters["short_url"]

	cfg := &aws.Config{
		Region: aws.String(endpoints.EuWest2RegionID),
	}

	sess := session.Must(session.NewSession(cfg))

	service := dynamodb.New(sess)

	input := &dynamodb.GetItemInput{
		TableName: aws.String(LinksTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"short_url": {
				S: aws.String(shortURL),
			},
		},
	}

	result, err := service.GetItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	link := Link{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &link); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	header := map[string]string{
		"location": link.LongURL,
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusPermanentRedirect,
		Headers:    header,
	}, nil

}

func main() {
	lambda.Start(Handler)
}

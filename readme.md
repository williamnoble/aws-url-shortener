# AWS URL Shortener


## Add Project Dependencies
```shell
# for Lambda
$ go get github.com/aws/aws-lambda-go

# for AWS (session/dynamodb)
$ go get github.com/aws/aws-sdk-go
```

## Setup DynamoDB
```shell
# this will create our UrlShortenerLinks table
$ aws dynamodb create-table --cli-input-json file://dynamodb/urlshortenerlinks_table.json
```

##### Build project for upload to Lambda
```shell
# Build our executable for AWS
 GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main main.go
 
# Zip our executable for AWS
zip -j main.zip main
```

##### Test the function
```shell
# Invoke our function through Lambda directly does not work because within our Handler it's expecting an events.APIGatewayProxyRequest!
$ aws lambda invoke --function-name ShortenFunction --cli-binary-format raw-in-base64-out --payload '{"url":"www.google.com/aVeryLongURL"}' out

$ curl {staging_url} -H "Content-Type: application/json" -X POST -d '{"url":"https://www.google.com""}'
{"short_url":"yCNoz67nR"}⏎

Alternative:
$ export STAGING_URL=https://www.something.execute-api.eu-west-2.amazonaws.com/prod/shorten
$ curl $STAGING_URL -H "Content-Type: application/json" -X POST -d '{"url":"https://www.google.com"}'
 => example Resposne: {"short_url":"yCNoz67nR"}⏎

```




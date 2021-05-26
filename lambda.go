package lambda

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func InvokeAuthorizer(functionName string, awsRegion string, request events.APIGatewayCustomAuthorizerRequestTypeRequest) (*events.APIGatewayCustomAuthorizerResponse, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	client := lambda.New(sess, &aws.Config{Region: aws.String(awsRegion)})

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	result, err := client.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		InvocationType: aws.String("RequestResponse"),
		LogType:        aws.String("Tail"),
		Payload:        payload,
	})
	if err != nil {
		return nil, err
	}

	var resp events.APIGatewayCustomAuthorizerResponse
	err = json.Unmarshal(result.Payload, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Context == nil {
		err = errors.New("Unauthorized.")
		return nil, err
	}

	return &resp, nil
}

func Invoke(functionName string, awsRegion string, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	client := lambda.New(sess, &aws.Config{Region: aws.String(awsRegion)})

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	result, err := client.Invoke(&lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		InvocationType: aws.String("RequestResponse"),
		LogType:        aws.String("Tail"),
		Payload:        payload,
	})
	if err != nil {
		return nil, err
	}

	var resp events.APIGatewayProxyResponse

	err = json.Unmarshal(result.Payload, &resp)
	if err != nil {
		return nil, err
	}

	// If the status code is NOT 200, the call failed
	if resp.StatusCode != 200 {
		log := base64.StdEncoding.EncodeToString([]byte(*result.LogResult))
		err := errors.New("Error StatusCode: " + strconv.Itoa(resp.StatusCode) + log)
		return nil, err
	}

	return &resp, nil
}
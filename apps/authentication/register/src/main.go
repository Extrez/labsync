package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type RequestBody struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"fullName"`
	PhoneNumber string `json:"phoneNumber"`
}

type ResponseBody struct {
	Message     string
	Success     bool
	AccessToken string
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest, cognitoClient *cognitoidentityprovider.CognitoIdentityProvider) (events.APIGatewayProxyResponse, error) {
	var requestBody RequestBody
	err := json.Unmarshal([]byte(request.Body), &requestBody)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, err
	}

	signUpInput := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String("Cognito_client_id"),
		Username: aws.String(requestBody.Username),
		Password: aws.String(requestBody.Password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(requestBody.Email),
			},
			{
				Name:  aws.String("fullName"),
				Value: aws.String(requestBody.FullName),
			},
			{
				Name:  aws.String("phoneNumber"),
				Value: aws.String(requestBody.PhoneNumber),
			},
		},
	}

	_, err = cognitoClient.SignUp(signUpInput)
	
	if err != nil {
		log.Printf("Failed at creating new user: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}	

	responseBody := ResponseBody{
		Message:     "Registration successful",
		Success:     true,
		AccessToken: "your_access_token",
	}
	resBodyBytes, _ := json.Marshal(responseBody)

	return events.APIGatewayProxyResponse{
		Body:       string(resBodyBytes),
		StatusCode: http.StatusCreated,
	}, nil
}

func main() {
	// Initialize AWS Session
    sess, err := session.NewSession()
    if err != nil {
        panic(err)
    }

    // Create Cognito client
    cognitoClient := cognitoidentityprovider.New(sess)

    // Start the lambda function with the real Cognito client
    lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
        return Handler(ctx, request, cognitoClient)
    })
}

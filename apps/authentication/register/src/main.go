package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
	Message     string `json:"message"`
	Success     bool   `json:"success"`
	AccessToken string `json:"access_token"`
}

func ComputeSecretHash(clientId, clientSecret, username string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest, cognitoClient *cognitoidentityprovider.CognitoIdentityProvider) (events.APIGatewayProxyResponse, error) {
	var requestBody RequestBody
	err := json.Unmarshal([]byte(request.Body), &requestBody)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}

	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	clientSecret := os.Getenv("COGNITO_CLIENT_SECRET")

	secretHash := ComputeSecretHash(userPoolID, clientSecret, requestBody.Username)

	signUpInput := &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(userPoolID),
		SecretHash: aws.String(secretHash),
		Username:   aws.String(requestBody.Username),
		Password:   aws.String(requestBody.Password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(requestBody.Email),
			},
			{
				Name:  aws.String("name"),
				Value: aws.String(requestBody.FullName),
			},
			{
				Name:  aws.String("phone_number"),
				Value: aws.String(requestBody.PhoneNumber),
			},
		},
	}

	_, err = cognitoClient.SignUp(signUpInput)

	if err != nil {
		responseBody := ResponseBody{
			Message:     "Registration successful",
			Success:     false,
			AccessToken: "Internal Server Error",
		}

		resBodyBytes, _ := json.Marshal(responseBody)

		if awsErr, ok := err.(awserr.Error); ok {
			// Check if the error is UsernameExistsException
			if awsErr.Code() == cognitoidentityprovider.ErrCodeUsernameExistsException {
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusConflict,
					Body:       string(resBodyBytes),
				}, nil
			}
		}

		log.Printf("Failed at creating new user: %v", err)

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(resBodyBytes),
		}, err
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
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	cognitoClient := cognitoidentityprovider.New(sess)

	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return Handler(ctx, request, cognitoClient)
	})
}

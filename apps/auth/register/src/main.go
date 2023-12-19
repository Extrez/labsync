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
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type RequestBody struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FullName    string `json:"fullName" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
}

type ResponseBody struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type CognitoClient interface {
	SignUp(*cognitoidentityprovider.SignUpInput) (*cognitoidentityprovider.SignUpOutput, error)
}

func ComputeSecretHash(clientId, clientSecret, username string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func isValidPassword(password string) bool {
	var (
		hasMinLen    = len(password) >= 8
		hasUppercase = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLowercase = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasNumber    = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasSpecial   = regexp.MustCompile(`[\W_]`).MatchString(password)
	)

	return hasMinLen && hasUppercase && hasLowercase && hasNumber && hasSpecial
}

func getEnvVar(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s must be set", key)
	}
	return value
}

func createResponse(response interface{}, statusCode int) (events.APIGatewayProxyResponse, error) {
	resBodyBytes, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(resBodyBytes),
	}, nil
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest, cognitoClient CognitoClient) (events.APIGatewayProxyResponse, error) {
	var requestBody RequestBody
	err := json.Unmarshal([]byte(request.Body), &requestBody)

	if err != nil {
		log.Printf("Error unmarshalling request: %v", err)
		responseBody := ResponseBody{
			Message: "Invalid request body",
			Success: false,
		}

		return createResponse(responseBody, http.StatusBadRequest)
	}

	if !isValidEmail(requestBody.Email) {
		responseBody := ResponseBody{
			Message: "Invalid email",
			Success: false,
		}

		return createResponse(responseBody, http.StatusBadRequest)
	}

	if !isValidPassword(requestBody.Password) {
		responseBody := ResponseBody{
			Message: "Invalid password",
			Success: false,
		}

		return createResponse(responseBody, http.StatusBadRequest)
	}

	clientID := getEnvVar("COGNITO_CLIENT_ID")
	clientSecret := getEnvVar("COGNITO_CLIENT_SECRET")

	secretHash := ComputeSecretHash(clientID, clientSecret, requestBody.Username)

	signUpInput := &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(clientID),
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

	if err == nil {
		responseBody := ResponseBody{
			Message: "Registration successful",
			Success: true,
		}

		return createResponse(responseBody, http.StatusCreated)
	}

	// Log and handle different types of errors
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == cognitoidentityprovider.ErrCodeUsernameExistsException {
			responseBody := ResponseBody{
				Message: "User already exists",
				Success: false,
			}

			return createResponse(responseBody, http.StatusConflict)
		}
	}

	log.Printf("Failed at creating new user: %v", err)

	responseBody := ResponseBody{
		Message: "Internal server error",
		Success: false,
	}

	return createResponse(responseBody, http.StatusInternalServerError)
}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	cognitoClient := cognitoidentityprovider.New(sess)

	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return Handler(ctx, request, cognitoClient)
	})
}

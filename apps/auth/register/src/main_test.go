package main

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/stretchr/testify/assert"
)

type MockCognitoClient struct {
	cognitoidentityprovider.CognitoIdentityProvider
	// Add additional fields or methods if needed
}

func (m *MockCognitoClient) SignUp(input *cognitoidentityprovider.SignUpInput) (*cognitoidentityprovider.SignUpOutput, error) {
	// Simulate the behavior when an invalid email is provided
	if input.Username != nil && *input.Username == "nonewuser" {
		return nil, awserr.New(cognitoidentityprovider.ErrCodeUsernameExistsException, "", nil)
	}
	// Handle other cases or default behavior
	return &cognitoidentityprovider.SignUpOutput{}, nil
}

func TestSuccessfulRegistration(t *testing.T) {
	os.Setenv("COGNITO_CLIENT_SECRET", "client_secret")
	os.Setenv("COGNITO_CLIENT_ID", "client_id")
	cognitoClient := &MockCognitoClient{}

	t.Run("Successful Registration", func(t *testing.T) {
		request := events.APIGatewayProxyRequest{
			Body: `{"username":"testuser", "email":"test@example.com", "password":"Password123!", "fullName":"Test User", "phoneNumber":"+11234567890"}`,
		}

		response, err := Handler(context.Background(), request, cognitoClient)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})

	t.Run("Invalid Email Format", func(t *testing.T) {
		request := events.APIGatewayProxyRequest{
			Body: `{"username":"testuser", "email":"invalid_email", "password":"Password123!", "fullName":"Test User", "phoneNumber":"+11234567890"}`,
		}

		response, err := Handler(context.Background(), request, cognitoClient)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Weak Password", func(t *testing.T) {
		request := events.APIGatewayProxyRequest{
			Body: `{"username":"testuser", "email":"test@example.com", "password":"password123", "fullName":"Test User", "phoneNumber":"+11234567890"}`,
		}

		response, err := Handler(context.Background(), request, cognitoClient)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("User Already Exists", func(t *testing.T) {
		request := events.APIGatewayProxyRequest{
			Body: `{"username":"nonewuser", "email":"test@example.com", "password":"Password123!", "fullName":"Test User", "phoneNumber":"+11234567890"}`,
		}

		response, err := Handler(context.Background(), request, cognitoClient)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, response.StatusCode)
	})

	t.Run("Invalid Input Data", func(t *testing.T) {
		request := events.APIGatewayProxyRequest{
			Body: `{"username":"nonewuser", "email":"test@example.com", "password":"password123", "fullName":"Test User", "phoneNumber":"+11234567890"}`,
		}

		response, err := Handler(context.Background(), request, cognitoClient)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
}

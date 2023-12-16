package main

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/awstesting/mock"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulRegistration(t *testing.T) {
    sess := mock.Session
    sess.Config.Region = aws.String("us-west-2")
    cognitoClient := cognitoidentityprovider.New(sess)

    t.Run("Successful Registration", func(t *testing.T) {
        request := events.APIGatewayProxyRequest{
            Body: `{"username":"testuser", "email":"test@example.com", "password":"password123", "fullName":"Test User", "phoneNumber":"1234567890"}`,
        }

        response, err := Handler(context.Background(), request, cognitoClient)

        assert.NoError(t, err)
        assert.Equal(t, http.StatusCreated, response.StatusCode)
    })
}
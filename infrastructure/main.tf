terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "us-east-1"
}

# data "aws_iam_policy_document" "assume_role" {
#   statement {
#     effect = "Allow"

#     principals {
#       type        = "Service"
#       identifiers = ["lambda.amazonaws.com"]
#     }

#     actions = ["sts:AssumeRole"]
#   }
# }

# resource "aws_iam_role" "iam_for_lambda" {
#   name               = "iam_for_lambda"
#   assume_role_policy = data.aws_iam_policy_document.assume_role.json
# }

# resource "aws_lambda_function" "test_lambda" {
#   filename = "lambda_function.zip"
#   //source_code_hash = filebase64sha256("lambda_function.zip")
#   function_name = "test_lambda"
#   handler       = "main"
#   runtime       = "go1.x"

#   role = aws_iam_role.iam_for_lambda.arn

#   environment {
#     variables = {
#       COGNITO_USER_POOL_ID = "your_user_pool_id"
#       COGNITO_CLIENT_ID    = "your_client_id"
#     }
#   }
# }

# resource "aws_api_gateway_rest_api" "api" {
#   name        = "GoLambdaAPI"
#   description = "API Gateway for Go Lambda Function"

#   endpoint_configuration {
#     types = ["REGIONAL"]
#   }
# }

# resource "aws_api_gateway_resource" "api_resource" {
#   rest_api_id = aws_api_gateway_rest_api.api.id
#   parent_id   = aws_api_gateway_rest_api.api.root_resource_id
#   path_part   = "register"
# }

# resource "aws_api_gateway_deployment" "example" {
#   rest_api_id = aws_api_gateway_rest_api.api.id

#   triggers = {
#     redeployment = sha1(jsonencode(aws_api_gateway_rest_api.api.body))
#   }

#   lifecycle {
#     create_before_destroy = true
#   }
# }

# resource "aws_api_gateway_stage" "example" {
#   deployment_id = aws_api_gateway_deployment.example.id
#   rest_api_id   = aws_api_gateway_rest_api.api.id
#   stage_name    = "example"
# }

# resource "aws_api_gateway_method" "api_method" {
#   rest_api_id   = aws_api_gateway_rest_api.api.id
#   resource_id   = aws_api_gateway_resource.api_resource.id
#   http_method   = "POST"
#   authorization = "NONE"
# }

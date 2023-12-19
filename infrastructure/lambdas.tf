data "aws_iam_policy_document" "lambda" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

# Register lambda

data "aws_iam_policy_document" "register" {
  statement {
    effect    = "Allow"
    actions   = ["cognito-idp:AdminCreateUser"]
    resources = [aws_cognito_user_pool.user_pool.arn]
  }
}

resource "aws_iam_role" "iam_for_lambda" {
  name               = "iam_for_lambda"
  assume_role_policy = data.aws_iam_policy_document.lambda.json

  inline_policy {
    name   = "cognito_policy"
    policy = data.aws_iam_policy_document.register.json
  }
}

data "archive_file" "register_lambda" {
  type        = "zip"
  source_file = "${path.module}/../dist/auth/register/bootstrap"
  output_path = "${path.module}/../dist/auth/register/register.zip"
}

resource "aws_lambda_function" "register" {
  filename      = data.archive_file.register_lambda.output_path
  function_name = "${var.environment}_register"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "bootstrap"

  source_code_hash = data.archive_file.register_lambda.output_base64sha256

  runtime       = "provided.al2023"
  architectures = ["x86_64"]

  environment {
    variables = {
      COGNITO_CLIENT_ID     = aws_cognito_user_pool_client.user_pool_client.id
      COGNITO_CLIENT_SECRET = aws_cognito_user_pool_client.user_pool_client.client_secret
    }
  }
}

resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.register.function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.api_gateway.execution_arn}/*"
}

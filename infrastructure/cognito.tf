resource "aws_cognito_user_pool" "user_pool" {
  name = "${var.environment}_user_pool"

  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_numbers   = true
    require_symbols   = true
    require_uppercase = true
  }

  tags = {
    environment = var.environment
  }
}


resource "aws_cognito_user_pool_client" "user_pool_client" {
  name         = "${var.environment}_user_pool_client"
  user_pool_id = aws_cognito_user_pool.user_pool.id

  // Define other app client settings as needed, for example:
  // allowed_oauth_flows_user_pool_client = true
  // allowed_oauth_flows = ["code", "implicit"]
  // allowed_oauth_scopes = ["phone", "email", "openid", "profile", "aws.cognito.signin.user.admin"]
  // callback_urls = ["https://www.example.com/callback"]
  // logout_urls = ["https://www.example.com/logout"]

  generate_secret = true
}

// Optionally, define a domain for the user pool
# resource "aws_cognito_user_pool_domain" "my_user_pool_domain" {
#   domain       = "my-unique-app-domain"
#   user_pool_id = aws_cognito_user_pool.my_user_pool.id
# }

terraform {
  backend "s3" {
    bucket         = "darklab-eu-west-1-staging-terraform-state"
    dynamodb_table = "darklab-eu-west-1-staging-terraform-state"
    encrypt        = true
    key            = "darkbot.tfstate"
    region         = "eu-west-1"
  }
}

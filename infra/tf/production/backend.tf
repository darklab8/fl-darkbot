terraform {
  backend "s3" {
    bucket         = "darklab-eu-west-1-production-terraform-state"
    dynamodb_table = "darklab-eu-west-1-production-terraform-state"
    encrypt        = true
    key            = "darkbot.tfstate"
    region         = "eu-west-1"
  }
}

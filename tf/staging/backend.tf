terraform {
  backend "s3" {
    bucket         = "darklab-eu-west-1-global-terraform-state"
    dynamodb_table = "darklab-eu-west-1-global-terraform-state"
    encrypt        = true
    key            = "darkbot.staging.tfstate"
    region         = "eu-west-1"
  }
}

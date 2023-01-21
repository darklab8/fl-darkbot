terraform {
  backend "s3" {
    bucket         = "darklab-eu-west-1-repo-terraform-state"
    dynamodb_table = "darklab-eu-west-1-repo-terraform-state"
    encrypt        = true
    key            = "darkbot.tfstate"
    region         = "eu-west-1"
  }
}

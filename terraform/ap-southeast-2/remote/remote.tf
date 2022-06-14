# This is the provider version we're using for this deployment
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.17.1"
    }
  }
}

# This is the aws profile and aws region we're using for this deployment
provider "aws" {
  profile = "jponc"
  region  = "ap-southeast-2"
}

# This is where we store our shared terraform state
resource "aws_s3_bucket" "terraform_state_bucket" {
  bucket = "estimatex-terraform-state-ap-southeast-2"
}

resource "aws_s3_bucket_server_side_encryption_configuration" "terraform_state_bucket_encryption" {
  bucket = aws_s3_bucket.terraform_state_bucket.bucket

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# This is where we store our shared terraform lock
resource "aws_dynamodb_table" "terraform_locks" {
  name         = "estimatex-terraform-locks-ap-southeast-2"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}

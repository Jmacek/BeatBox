terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}

# Set the aws region closest to you
provider "aws" {
  region = "us-west-2"
}

# Create DynamoDB tables
resource "aws_dynamodb_table" "User" {
  name           = "User"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "Email"

  attribute {
    name = "Email"
    type = "S"
  }
}
resource "aws_dynamodb_table" "Artist" {
  name           = "Artist"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "ArtistID"

  attribute {
    name = "ArtistID"
    type = "S"
  }
  attribute {
    name = "ArtistName"
    type = "S"
  }
  global_secondary_index {
    name               = "ArtistName-ArtistID-index"
    hash_key           = "ArtistName"
    write_capacity     = 10
    read_capacity      = 10
    projection_type    = "KEYS_ONLY"
  }
}
resource "aws_dynamodb_table" "Album" {
  name           = "Album"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "AlbumID"

  attribute {
    name = "AlbumID"
    type = "S"
  }
  attribute {
    name = "AlbumName"
    type = "S"
  }
  global_secondary_index {
    name               = "AlbumName-AlbumID-index"
    hash_key           = "AlbumName"
    write_capacity     = 10
    read_capacity      = 10
    projection_type    = "KEYS_ONLY"
  }
}
resource "aws_dynamodb_table" "Song" {
  name           = "Song"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "SongID"

  attribute {
    name = "SongID"
    type = "S"
  }
  attribute {
    name = "SongName"
    type = "S"
  }
  global_secondary_index {
    name               = "SongName-SongID-index"
    hash_key           = "SongName"
    write_capacity     = 10
    read_capacity      = 10
    projection_type    = "KEYS_ONLY"
  }
}

# Create S3 bucket
resource "aws_s3_bucket" "s3_bucket" {
  bucket = "beatbox-songbucket"

  tags = {
    Name = "Bucket to store songs for beatbox service"
  }
}

# set ACL for S3 bucket
resource "aws_s3_bucket_acl" "s3_bucket_acl" {
  bucket = aws_s3_bucket.s3_bucket.id
  acl    = "private"
}

##### Below is cloudfront distribution. Uncomment to create (note: not implimented in code to pull from cloudfront)
# locals {
#   s3_origin_id = "beatbox-songbucket"
# }
#
# resource "aws_cloudfront_distribution" "s3_distribution" {
#   origin {
#     domain_name              = aws_s3_bucket.s3_bucket.bucket_regional_domain_name
#     origin_id                = local.s3_origin_id
#   }
#
#   enabled             = true
#   is_ipv6_enabled     = true
#   comment             = "CDN for S3 beatbox-bucket"
#
#   default_cache_behavior {
#     allowed_methods  = ["GET", "HEAD"]
#     cached_methods   = ["GET", "HEAD"]
#     target_origin_id = local.s3_origin_id
#
#     forwarded_values {
#       query_string = false
#
#       cookies {
#         forward = "none"
#       }
#     }
#
#     viewer_protocol_policy = "allow-all"
#     min_ttl                = 0
#     default_ttl            = 3600
#     max_ttl                = 86400
#   }
#
#   restrictions {
#     geo_restriction {
#       restriction_type = "whitelist"
#       locations        = ["US", "CA", "GB", "DE"]
#     }
#   }
#
#   tags = {
#     Environment = "dev"
#   }
#
#   viewer_certificate {
#     cloudfront_default_certificate = true
#   }
# }
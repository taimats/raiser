terraform {
  required_version = "1.11.4"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
    http = {
      source  = "hashicorp/http"
      version = "~> 3.4"
    }
  }
}

provider "aws" {
  profile = var.profile
  region  = var.region
}

provider "http" {}

variable "project" {
  type        = string
  description = "project名"
}

variable "env" {
  type        = string
  description = "インフラの環境名"
}

variable "profile" {
  type        = string
  description = "awsアクセス時のユーザー名"
}

variable "region" {
  type        = string
  description = "awsのリージョン"
}
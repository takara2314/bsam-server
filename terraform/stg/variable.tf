variable "project" {
  description = "Google Cloud project ID"
  default     = "bsam-stg"
}

variable "environment" {
  description = "Environment name, 'stg' or 'prd'"
  default     = "stg"
}

variable "location" {
  description = "Google Cloud location"
  default     = "asia-northeast1"
}

variable "api_service_image" {
  description = "B-SAM API Service Docker image URL"
  type        = string
}

variable "game_service_image" {
  description = "B-SAM Game Service Docker image URL"
  type        = string
}

variable "auth_service_image" {
  description = "B-SAM Auth Service Docker image URL"
  type        = string
}

variable "domain_name" {
  description = "B-SAM domain name"
  default     = "bsam.app"
}

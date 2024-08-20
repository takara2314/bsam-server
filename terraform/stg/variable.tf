variable "project" {
  type    = string
  default = "bsam-stg"
}

variable "location" {
  type    = string
  default = "asia-northeast1"
}

variable "api_service_image" {
  type        = string
  description = "B-SAM API Service Docker image URL"
}

variable "game_service_image" {
  type        = string
  description = "B-SAM Game Service Docker image URL"
}

variable "auth_service_image" {
  type        = string
  description = "B-SAM Auth Service Docker image URL"
}

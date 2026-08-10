variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "region" {
  description = "GCP Region (e.g., asia-southeast1)"
  type        = string
  default     = "asia-southeast1"
}

variable "environment" {
  description = "Deployment Environment (e.g., dev, prod)"
  type        = string
}

variable "service_name" {
  description = "ชื่อของ Cloud Run Service และ Repository"
  type        = string
}

variable "github_owner" {
  description = "ชื่อเจ้าของ GitHub Account หรือ Organization"
  type        = string
}

variable "github_repo" {
  description = "ชื่อของ Repository ใน GitHub"
  type        = string
}

variable "branch_name" {
  description = "Branch ที่ต้องการให้ Trigger ทำงานเมื่อมีการ Push (เช่น ^main$)"
  type        = string
  default     = "^main$"
}

variable "db_password_value" {
  description = "Sensitive Database Password"
  type        = string
  sensitive   = true 
}
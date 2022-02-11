terraform {
  experiments = [module_variable_optional_attrs]
}

variable "group" {
  type        = string
  description = "A unique identifier for the EBS group"
}

variable "attached_ebs" {
  type        = any
  description = "Map of the EBS objects to allocate"
}

variable "ebs_bootstrap_binary_url" {
  default     = null
  description = "Custom URL from which to download the ebs_bootstrap binary"
}

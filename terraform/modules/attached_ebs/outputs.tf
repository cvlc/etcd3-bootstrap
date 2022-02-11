output "iam_role_policy_arn" {
  description = "IAM role policy ARN to assign to ASG instance role"
  value       = aws_iam_policy.data.arn
}

output "userdata_snippets_by_az" {
  description = "Map of snippets of userdata to assign to ASG instances by availability zone"
  value       = local.user_data_snippets_by_az
}

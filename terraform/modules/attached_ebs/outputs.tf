output "iam_role_policy_document" {
  description = "IAM role policy document to assign to ASG instance role"
  value       = data.aws_iam_policy_document.ebs.json
}

output "userdata_snippets_by_az" {
  description = "Map of snippets of userdata to assign to ASG instances by availability zone"
  value       = local.user_data_snippets_by_az
}

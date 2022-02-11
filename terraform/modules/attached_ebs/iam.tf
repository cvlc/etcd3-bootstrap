resource "aws_iam_policy" "data" {
  name        = "${var.group}-data"
  path        = "/"
  description = "Allow data volume management for nodes from ${var.group}"
  policy      = data.aws_iam_policy_document.ebs.json
}

data "aws_region" "current" {}

data "aws_iam_policy_document" "ebs" {
  statement {
    sid = "allowec2read"
    actions = [
      "ec2:DescribeAddresses",
      "ec2:DescribeInstances",
      "ec2:DescribeVolumes",
      "ec2:DescribeVolumeStatus"
    ]
    resources = [
      "*"
    ]
  }
  statement {
    sid = "allowebsreassignment"
    actions = [
      "ec2:AttachVolume",
      "ec2:DetachVolume"
    ]
    resources = [
      "*"
    ]
    condition {
      test     = "StringEquals"
      values   = [var.group]
      variable = "ec2:ResourceTag/Group"
    }
    condition {
      test     = "StringEquals"
      values   = [data.aws_region.current.name]
      variable = "ec2:Region"
    }
  }
}

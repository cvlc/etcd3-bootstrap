resource "aws_iam_role" "dlm_lifecycle_role" {
  name               = "dlm-${var.group}"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "dlm.amazonaws.com"
      },
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "dlm_lifecycle" {
  name = "dlm-snapshot-lifecycle-policy-${var.group}"
  role = aws_iam_role.dlm_lifecycle_role.id

  policy = <<EOF
{
   "Version": "2012-10-17",
   "Statement": [
      {
         "Effect": "Allow",
         "Action": [
            "ec2:CreateSnapshot",
            "ec2:DeleteSnapshot",
            "ec2:CreateTags"
         ],
         "Resource": "*",
         "Condition": {
         "ForAllValues:StringEquals": {
            "ec2:ResourceTag/snap-daily": "true",
            "ec2:Region": "${data.aws_region.current.name}"
          }
        }
      },
      {
         "Effect": "Allow",
         "Action": [
            "ec2:DescribeVolumes",
            "ec2:DescribeSnapshots"
         ],
         "Resource": "*"
      }
   ]
}
EOF
}

resource "aws_dlm_lifecycle_policy" "automatic_snapshots" {
  description        = "DLM snapshot lifecycle policy"
  execution_role_arn = aws_iam_role.dlm_lifecycle_role.arn
  state              = "ENABLED"
  policy_details {
    resource_types = ["VOLUME"]

    schedule {
      name = "week of daily snapshots"

      create_rule {
        interval      = 24
        interval_unit = "HOURS"
        times         = ["02:15"]
      }

      retain_rule {
        count = 7
      }

      tags_to_add = {
        SnapshotCreator = "DLM"
      }

      copy_tags = true
    }

    target_tags = {
      "snap-daily" = "true"
    }
  }
}

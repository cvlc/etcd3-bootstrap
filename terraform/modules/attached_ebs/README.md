# attached_ebs
## Introduction
attached_ebs is a Terraform module used to generate persistent EBS volumes and attach them to auto-scaled instances, ensuring that snapshots are taken of them daily.

## Usage
Set the input variable `group` to a unique identifier to use to refer to the group of EBS disks. Make sure to tag all of your instances with `Group=XXXX` (with `XXXX` as the value you set for `group`). This is used by the IAM policy to enable permissions for attaching volumes within the group. 

The input variable `attached_ebs` takes a map of volume definitions to attach to instances on boot:
```
module "attached-ebs" {
  source = "github.com/ondat/etcd3-bootstrap/terraform/modules/attached_ebs"
  attached_ebs = { 
    "ondat_data_1": {
      size = 100 # required
      availability_zone = eu-west-1a # required
      encrypted = true
      volume_type = gp3
      block_device_aws = "/dev/xvda1"
      block_device_os = "/dev/nvme0n1"
      block_device_mount_path = "/var/lib/data0"
    }
    "ondat_data_2": {
      size = 100
      availability_zone = eu-west-1a
      dependson = ["ondat_data_1"]
      encrypted = true
      restore_snapshot = ""
      iops = 3000
      volume_type = io2
      throughput = 150000
      kms_key_id = "arn:aws::kms/..."
      block_device_aws = /dev/xvda2
      block_device_os = /dev/nvme1n1
      block_device_mount_path = /var/lib/data1
    }
  }
}
```

For airgapped or private environments, the variable `ebs_bootstrap_binary_url` can be used to provide an HTTP/S address from which to retrieve the necessary binary.

Use the output `iam_role_policy_document` to generate and assign the policy to your ASG node's role.
Use the output `userdata_snippets_by_az` to embed in your ASG's userdata - it's a map of AZ to snippets.

## Appendix

### Requirements

No requirements.

### Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | n/a |

### Modules

No modules.

### Resources

| Name | Type |
|------|------|
| [aws_dlm_lifecycle_policy.automatic_snapshots](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/dlm_lifecycle_policy) | resource |
| [aws_ebs_volume.ssd](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume) | resource |
| [aws_iam_role.dlm_lifecycle_role](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role) | resource |
| [aws_iam_role_policy.dlm_lifecycle](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/iam_role_policy) | resource |
| [aws_iam_policy_document.ebs](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/iam_policy_document) | data source |
| [aws_region.current](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/region) | data source |

### Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_attached_ebs"></a> [attached\_ebs](#input\_attached\_ebs) | Map of the EBS objects to allocate | `any` | n/a | yes |
| <a name="input_ebs_bootstrap_binary_url"></a> [ebs\_bootstrap\_binary\_url](#input\_ebs\_bootstrap\_binary\_url) | Custom URL from which to download the ebs\_bootstrap binary | `any` | `null` | no |
| <a name="input_group"></a> [group](#input\_group) | A unique identifier for the EBS group | `string` | n/a | yes |

### Outputs

| Name | Description |
|------|-------------|
| <a name="output_iam_role_policy_document"></a> [iam\_role\_policy\_document](#output\_iam\_role\_policy\_document) | IAM role policy document to assign to ASG instance role |
| <a name="output_userdata_snippets_by_az"></a> [userdata\_snippets\_by\_az](#output\_userdata\_snippets\_by\_az) | Map of snippets of userdata to assign to ASG instances by availability zone |


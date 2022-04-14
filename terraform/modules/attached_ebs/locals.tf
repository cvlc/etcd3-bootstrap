locals {
  availability_zones = distinct([for key, value in var.attached_ebs :
  value.availability_zone])
  user_data_snippets_by_az = {
    for az in local.availability_zones :
    az => join("\n", [
      for key, value in var.attached_ebs : value["availability_zone"] == az ? templatefile("${path.module}/cloudinit/userdata-snippet.sh", {
        region          = data.aws_region.current.name,
        ebs_volume_name = "${key}-${value["availability_zone"]}",
        ebs_bootstrap_unit = templatefile("${path.module}/cloudinit/ebs_bootstrap_unit", {
          region                      = data.aws_region.current.name
          depends                     = lookup(value, "dependson", false) ? join("\n", [for e in value["dependson"] : "Requires=${e}-${value["availability_zone"]}.service\nAfter=${e}-${value["availability_zone"]}.service"]) : ""
          ebs_volume_name             = "${key}-${value["availability_zone"]}"
          ebs_block_device_aws        = value["block_device_aws"]
          ebs_block_device_os         = value["block_device_os"]
          ebs_block_device_mount_path = value["block_device_mount_path"]
          ebs_bootstrap_binary_url    = local.ebs_bootstrap_binary_url
        }),
    }) : ""])
  }
  ebs_bootstrap_binary     = "https://github.com/ondat/etcd3-bootstrap/releases/download/v0.1.1/ebs-bootstrap-linux-amd64"
  ebs_bootstrap_binary_url = var.ebs_bootstrap_binary_url == null ? local.ebs_bootstrap_binary : var.ebs_bootstrap_binary_url

}

resource "aws_ebs_volume" "ssd" {
  for_each = { for key, value in var.attached_ebs : key => value }

  snapshot_id       = lookup(each.value, "restore_snapshot", null)
  availability_zone = each.value["availability_zone"]
  size              = each.value["size"]
  iops              = lookup(each.value, "iops", null)
  type              = lookup(each.value, "volume_type", "gp3")
  kms_key_id        = lookup(each.value, "kms_key_id", null)
  encrypted         = lookup(each.value, "encrypted", true)
  throughput        = lookup(each.value, "throughput", null)

  tags = merge({
    Name         = "${each.key}-${each.value["availability_zone"]}"
    Group        = var.group
    "snap-daily" = "true"
    DependsOn    = join(",", lookup(each.value, "dependson", ["None"]))
  }, lookup(each.value, "tags", {}))
}

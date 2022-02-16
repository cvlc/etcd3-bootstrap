# ebs-bootstrap
## Introduction
This tool, written in Golang and based on [Monzo's etcd3-bootstap](https://github.com/monzo/etcd3-bootstrap) bootstraps an EBS volume for an EC2 instance. 

## Terraform
Included is a [Terraform module](terraform/modules/attached_ebs) allowing for 1..m EBS volumes to be automatically attached to autoscaling instances for convenience and automation. Documentation can be found [here](terraform/modules/attached_ebs/README.md).

package main

import (
	"context"
	"flag"
	"io"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/mvisonneau/go-ebsnvme/pkg/ebsnvme"
)

var (
	useEBS                    bool
	ebsVolumeName             string
	mountPoint                string
	blockDeviceAWS            string
	blockDeviceOS             string
	awsRegion                 string
	fileSystemFormatType      string
	fileSystemFormatArguments string
)

func init() {
	flag.StringVar(&awsRegion, "aws-region", "eu-west-1", "AWS region this instance is on")
	flag.StringVar(&ebsVolumeName, "ebs-volume-name", "", "EBS volume to attach to this node")
	flag.StringVar(&mountPoint, "mount-point", "/var/lib/etcd", "EBS volume mount point")
	flag.StringVar(&blockDeviceAWS, "block-device-aws", "/dev/xvdf", "Block device to attach as from AWS's perspective")
	flag.StringVar(&blockDeviceOS, "block-device-os", "/dev/nvme1n1", "Block device to attach as from OS's perspective")
	flag.StringVar(&fileSystemFormatType, "filesystem-type", "ext4", "Linux filesystem format type")
	flag.StringVar(&fileSystemFormatArguments, "filesystem-arguments", "", "Linux filesystem format arguments")
	flag.BoolVar(&useEBS, "use-ebs", true, "Use EBS instead of instance store")
	flag.Parse()
}

func main() {
	// Initialize AWS session
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("Error loading AWS config")
	}

	// Create ec2 and metadata svc clients with specified region
	ec2SVC := ec2.NewFromConfig(awsConfig)
	metadataSVC := imds.NewFromConfig(awsConfig)

	// obtain current AZ, required for finding volume
	azQueryOutput, err := metadataSVC.GetMetadata(context.TODO(), &imds.GetMetadataInput{Path: "placement/availability-zone"})
	if err != nil {
		panic("Error querying AWS metadata")
	}
	availabilityZoneByte, err := io.ReadAll(azQueryOutput.Content)
	if err != nil {
		panic(err)
	}

	if useEBS {
		volume, err := volumeFromName(ec2SVC, ebsVolumeName, string(availabilityZoneByte))
		if err != nil {
			panic(err)
		}

		instanceIDQueryOutput, err := metadataSVC.GetMetadata(context.TODO(), &imds.GetMetadataInput{Path: "instance-id"})
		if err != nil {
			panic(err)
		}
		instanceIDByte, err := io.ReadAll(instanceIDQueryOutput.Content)

		err = attachVolume(ec2SVC, string(instanceIDByte), volume)
		if err != nil {
			panic(err)
		}

		devices, err := os.ReadDir("/dev/")
		if err != nil {
			panic(err)
		}

		for _, d := range devices {
			valid, err := regexp.MatchString(`/dev/nvme[0-9]+`, d.Name())
			if err != nil {
				panic(err)
			}
			if valid {
				device, err := ebsnvme.ScanDevice(d.Name())
				if err != nil {
					panic(err)
				}

				if device.VolumeID == *volume.VolumeId {
					blockDeviceOS = device.Name
					break
				}
			}
		}
	}

	if err := ensureVolumeInited(blockDeviceOS, fileSystemFormatType, fileSystemFormatArguments); err != nil {
		panic(err)
	}

	if err := ensureVolumeMounted(blockDeviceOS, mountPoint); err != nil {
		panic(err)
	}

}

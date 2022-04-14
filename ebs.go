package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/jpillora/backoff"
)

func volumeFromName(
	svc *ec2.Client,
	volumeName string,
	az string,
) (*types.Volume, error) {
	b := &backoff.Backoff{
		Min:    5 * time.Second,
		Max:    350 * time.Second,
		Factor: 2,
		Jitter: false,
	}

	for {
		time.Sleep(b.Duration())
		input := &ec2.DescribeVolumesInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("tag:Name"),
					Values: []string{*aws.String(volumeName)},
				},
				{
					Name:   aws.String("availability-zone"),
					Values: []string{*aws.String(az)},
				},
			},
		}

		result, err := svc.DescribeVolumes(context.TODO(), input)
		if err != nil {
			return nil, err
		}

		if len(result.Volumes) == 0 {
			log.Printf("cannot find volume-id with name: %s, retrying...", volumeName)
			continue
		}

		log.Printf("Resolved volume %s to %s\n", volumeName, *result.Volumes[0].VolumeId)
		return &result.Volumes[0], nil
	}
}

func attachVolume(svc *ec2.Client, instanceID string, volume *types.Volume) error {
	log.Printf("Will attach volume %s to instance id %s\n", *volume.VolumeId, instanceID)

	// check if volume is already attached to this instance (ie, reboot)
	if len(volume.Attachments) > 0 && *volume.Attachments[0].InstanceId == instanceID {
		log.Printf("Volume %s is already attached to instance %s as device %s\n",
			*volume.VolumeId, instanceID, *volume.Attachments[0].Device)
		return nil
	}

	input := &ec2.AttachVolumeInput{
		Device:     aws.String(blockDeviceAWS),
		InstanceId: aws.String(instanceID),
		VolumeId:   volume.VolumeId,
	}

	_, err := svc.AttachVolume(context.TODO(), input)
	if err != nil {
		return err
	}

	b := &backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    100 * time.Second,
		Factor: 2,
		Jitter: false,
	}

	for {
		time.Sleep(b.Duration())
		volumeDescs, err := svc.DescribeVolumes(context.TODO(), &ec2.DescribeVolumesInput{
			VolumeIds: []string{*volume.VolumeId},
		})
		if err != nil {
			return err
		}

		volumes := volumeDescs.Volumes
		if len(volumes) == 0 {
			continue
		}

		if len(volumes[0].Attachments) == 0 {
			continue
		}

		if volumes[0].Attachments[0].State == types.VolumeAttachmentStateAttached {
			break
		}

		log.Printf(
			"Waiting for attachment to complete. Current state: %s",
			volumes[0].Attachments[0].State,
		)
	}

	log.Printf("Attached volume %s to instance %s as device %s\n",
		*volume.VolumeId, instanceID, blockDeviceAWS)

	return nil
}

func ensureVolumeInited(
	blockDeviceOS, fileSystemFormatType, fileSystemFormatArguments string,
) error {
	b := &backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    100 * time.Second,
		Factor: 2,
		Jitter: false,
	}

	log.Printf("Checking for existing filesystem on device: %s\n", blockDeviceOS)

	if err := exec.Command("sudo", "blkid", blockDeviceOS).Run(); err == nil {
		log.Println("Found existing filesystem")
		return nil
	}

	log.Println("Filesystem not present")

	// format volume here
	for {
		time.Sleep(b.Duration())
		cmd := exec.Command("sudo", "mkfs."+fileSystemFormatType, blockDeviceOS)
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		if err := cmd.Run(); err != nil {
			log.Println(err)
		} else {
			break
		}
	}

	return nil
}

func ensureVolumeMounted(blockDeviceOS, mountPoint string) error {
	log.Printf("Mounting device %s at %s\n", blockDeviceOS, mountPoint)

	// ensure mount point exists
	cmd := exec.Command("sudo", "mkdir", "-p", mountPoint)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("sudo", "mount", blockDeviceOS, mountPoint)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err == nil {
		log.Printf("Device %s successfully mounted at %s\n", blockDeviceOS, mountPoint)
		return nil
	}

	// mount failed, double-check as this may result from a previous mount
	log.Println("Mount failed. perhaps already mounted, will double check")
	out, err := exec.Command("mount").Output()
	if err != nil {
		return err
	}

	if strings.Contains(string(out), fmt.Sprintf("%s on %s", blockDeviceOS, mountPoint)) {
		log.Printf("Device %s successfully mounted at %s\n", blockDeviceOS, mountPoint)
		return nil
	}

	return err
}

package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2Instance struct {
	ID   string
	Name string
}

func ListEC2Instances(region string) ([]EC2Instance, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	// Create EC2 client
	client := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeInstancesInput{}
	output, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("unable to describe EC2 instances: %v", err)
	}

	var instances []EC2Instance

	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			var name string
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}
			instances = append(instances, EC2Instance{
				ID:   *instance.InstanceId,
				Name: name,
			})
		}
	}
	return instances, nil
}

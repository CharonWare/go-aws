package aws

import (
	"context"
	"fmt"

	"github.com/CharonWare/go-aws/internal/shared"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type EC2Instance struct {
	ID   string
	Name string
}

func ListEC2Instances(region string) ([]EC2Instance, error) {
	// Load AWS configuration
	cfg, err := shared.LoadAWSConfig(region)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	// Create EC2 client
	client := ec2.NewFromConfig(cfg)
	input := &ec2.DescribeInstancesInput{}
	var instances []EC2Instance

	// Use a paginator to ensure we see all the results
	paginator := ec2.NewDescribeInstancesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("unable to describe EC2 instances: %v", err)
		}
		for _, reservation := range page.Reservations {
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
	}

	return instances, nil
}

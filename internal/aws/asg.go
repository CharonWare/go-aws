package aws

import (
	"context"
	"fmt"

	"github.com/CharonWare/go-aws/internal/shared"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
)

type ASG struct {
	Name            string
	MinSize         int32
	MaxSize         int32
	DesiredCapacity int32
}

func DescribeASGs(region string) ([]ASG, error) {
	cfg, err := shared.LoadAWSConfig(region)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	client := autoscaling.NewFromConfig(cfg)
	input := &autoscaling.DescribeAutoScalingGroupsInput{}
	paginator := autoscaling.NewDescribeAutoScalingGroupsPaginator(client, input)

	var groups []ASG

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("unable to describe autoscaling groups: %v", err)
		}
		for _, AutoScalingGroups := range page.AutoScalingGroups {
			groups = append(groups, ASG{
				Name:            *AutoScalingGroups.AutoScalingGroupName,
				MinSize:         *AutoScalingGroups.MinSize,
				MaxSize:         *AutoScalingGroups.MaxSize,
				DesiredCapacity: *AutoScalingGroups.DesiredCapacity,
			})
		}
	}

	return groups, nil
}

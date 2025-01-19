package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/CharonWare/go-aws/internal/shared"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type ASG struct {
	Name            string
	MinSize         int32
	MaxSize         int32
	DesiredCapacity int32
	AVGCPU          float64
}

func DescribeASGs(region string) ([]ASG, error) {
	cfg, err := shared.LoadAWSConfig(region)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	client := autoscaling.NewFromConfig(cfg)
	input := &autoscaling.DescribeAutoScalingGroupsInput{}
	paginator := autoscaling.NewDescribeAutoScalingGroupsPaginator(client, input)

	// Time vars for cloudwatch queries
	startTime := time.Now().Add(-5 * time.Minute)
	endTime := time.Now()

	var groups []ASG
	var failedASGs []string // To track ASGs with no datapoints

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("unable to describe autoscaling groups: %v", err)
		}
		for _, AutoScalingGroups := range page.AutoScalingGroups {
			// Use the ASG names are the cloudwatch dimension
			dimensions := []types.Dimension{
				{Name: aws.String("AutoScalingGroupName"), Value: aws.String(*AutoScalingGroups.AutoScalingGroupName)},
			}
			// Define the struct for the GetMetricStats function
			metric := &MetricStats{
				Namespace:  "AWS/EC2",
				MetricName: "CPUUtilization",
				Dimensions: dimensions,
				StartTime:  &startTime,
				EndTime:    &endTime,
				Period:     300,
				Statistics: []types.Statistic{types.StatisticAverage},
			}
			// Get the average CPU usage for the last 5 min
			output, err := GetMetricStats(region, metric)
			if err != nil {
				// Log error and skip the ASG
				failedASGs = append(failedASGs, *AutoScalingGroups.AutoScalingGroupName)
				continue
			} else {
				groups = append(groups, ASG{
					Name:            *AutoScalingGroups.AutoScalingGroupName,
					MinSize:         *AutoScalingGroups.MinSize,
					MaxSize:         *AutoScalingGroups.MaxSize,
					DesiredCapacity: *AutoScalingGroups.DesiredCapacity,
					AVGCPU:          output,
				})
			}
		}
	}
	// log failed ASGs for user visibility
	if len(failedASGs) > 0 {
		fmt.Printf("Warning: No datapoints found for the following Auto Scaling Groups: %v\n", failedASGs)
	}

	return groups, nil
}

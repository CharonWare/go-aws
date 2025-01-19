package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/CharonWare/go-aws/internal/shared"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type MetricStats struct {
	Namespace  string
	MetricName string
	Dimensions []types.Dimension
	StartTime  *time.Time
	EndTime    *time.Time
	Period     int32
	Statistics []types.Statistic
}

func GetMetricStats(region string, m *MetricStats) (float64, error) {
	cfg, err := shared.LoadAWSConfig(region)
	if err != nil {
		return 0, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	client := cloudwatch.NewFromConfig(cfg)
	input := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String(m.Namespace),
		MetricName: aws.String(m.MetricName),
		Dimensions: m.Dimensions,
		StartTime:  m.StartTime,
		EndTime:    m.EndTime,
		Period:     &m.Period,
		Statistics: m.Statistics,
	}

	output, err := client.GetMetricStatistics(context.TODO(), input)
	if err != nil {
		return 0, fmt.Errorf("unable to get metric statistics from CloudWatch: %v", err)
	}

	if len(output.Datapoints) == 0 {
		return 0, fmt.Errorf("no datapoints found for metric %s with dimensions %v", *input.MetricName, input.Dimensions)
	}

	var sum float64
	for _, dp := range output.Datapoints {
		if dp.Average != nil {
			sum += *dp.Average
		}
	}

	average := sum / float64(len(output.Datapoints))
	return average, nil
}

package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/CharonWare/go-aws/internal/shared"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type describeCluster struct {
	Name           string
	ContainerHosts int32
	RunningTasks   int32
	PendingTasks   int32
	Services       int32
	AVGCPU         float64
}

type describeService struct {
	Name       string
	Desired    int32
	Running    int32
	Pending    int32
	LaunchType string
	AVGCPU     float64
}

func newECSClient(cfg aws.Config) *ecs.Client {
	return ecs.NewFromConfig(cfg)
}

func ListClusters(region string) (clusters []string, err error) {
	cfg, err := shared.LoadAWSConfig(region)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	// Create ECS client
	client := newECSClient(cfg)
	input := &ecs.ListClustersInput{}

	// Use a paginator to ensure we see all the results
	paginator := ecs.NewListClustersPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("unable to list ECS clusters: %v", err)
		}
		clusters = append(clusters, page.ClusterArns...)
	}

	return clusters, nil
}

func DescribeCluster(region, cluster string) ([]describeCluster, error) {
	cfg, err := shared.LoadAWSConfig(region)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	client := newECSClient(cfg)
	input := &ecs.DescribeClustersInput{
		Clusters: []string{cluster},
	}

	// Time vars for cloudwatch queries
	startTime := time.Now().Add(-5 * time.Minute)
	endTime := time.Now()

	output, err := client.DescribeClusters(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("unable to describe cluster: %v", err)
	}

	var chosenCluster []describeCluster

	for _, clusters := range output.Clusters {
		// Use the clusterName as the metric dimension
		dimensons := []types.Dimension{
			{Name: aws.String("ClusterName"), Value: aws.String(*clusters.ClusterName)},
		}
		// Define the struct for the GetMetricStats function
		metric := &MetricStats{
			Namespace:  "AWS/ECS",
			MetricName: "CPUUtilization",
			Dimensions: dimensons,
			StartTime:  &startTime,
			EndTime:    &endTime,
			Period:     300,
			Statistics: []types.Statistic{types.StatisticAverage},
		}
		// Get the average CPU usage for the last 5 min
		output, err := GetMetricStats(region, metric)
		if err != nil {
			return nil, fmt.Errorf("error: %v", err)
		}
		chosenCluster = append(chosenCluster, describeCluster{
			Name:           *clusters.ClusterName,
			ContainerHosts: clusters.RegisteredContainerInstancesCount,
			RunningTasks:   clusters.RunningTasksCount,
			Services:       clusters.ActiveServicesCount,
			AVGCPU:         output,
		})
	}

	return chosenCluster, nil
}

func ListServices(region, cluster string) (services []string, err error) {
	cfg, err := shared.LoadAWSConfig(region)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	// Create ECS client
	client := newECSClient(cfg)
	input := &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	}

	// Use a paginator to ensure we see all the results
	paginator := ecs.NewListServicesPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("unable to list ECS services: %v", err)
		}
		services = append(services, page.ServiceArns...)
	}

	return services, nil
}

func DescribeService(region, cluster, service string) ([]describeService, error) {
	cfg, err := shared.LoadAWSConfig(region)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	client := newECSClient(cfg)
	input := &ecs.DescribeServicesInput{
		Cluster:  aws.String(cluster),
		Services: []string{service},
	}

	// Time vars for cloudwatch queries
	startTime := time.Now().Add(-5 * time.Minute)
	endTime := time.Now()

	output, err := client.DescribeServices(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("unable to describe service: %v", err)
	}

	parts := strings.Split(cluster, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid cluster ARN: %s", cluster)
	}
	clusterName := parts[len(parts)-1]

	var chosenService []describeService

	for _, i := range output.Services {
		// Use the clusterName as the metric dimension
		dimensons := []types.Dimension{
			{Name: aws.String("ClusterName"), Value: aws.String(clusterName)},
			{Name: aws.String("ServiceName"), Value: aws.String(*i.ServiceName)},
		}
		// Define the struct for the GetMetricStats function
		metric := &MetricStats{
			Namespace:  "AWS/ECS",
			MetricName: "CPUUtilization",
			Dimensions: dimensons,
			StartTime:  &startTime,
			EndTime:    &endTime,
			Period:     300,
			Statistics: []types.Statistic{types.StatisticAverage},
		}
		// Get the average CPU usage for the last 5 min
		output, err := GetMetricStats(region, metric)
		if err != nil {
			return nil, fmt.Errorf("error: %v", err)
		}
		chosenService = append(chosenService, describeService{
			Name:       *i.ServiceName,
			Desired:    i.DesiredCount,
			Running:    i.RunningCount,
			Pending:    i.PendingCount,
			LaunchType: string(i.LaunchType),
			AVGCPU:     output,
		})
	}
	return chosenService, nil
}

func ListTasks(region, cluster, service string) (tasks []string, err error) {
	cfg, err := shared.LoadAWSConfig(region)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	// Create ECS client
	client := newECSClient(cfg)
	input := &ecs.ListTasksInput{
		Cluster:     aws.String(cluster),
		ServiceName: aws.String(service),
	}

	// Use a paginator to ensure we see all the results
	paginator := ecs.NewListTasksPaginator(client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("unable to list ECS tasks: %v", err)
		}
		tasks = append(tasks, page.TaskArns...)
	}

	return tasks, nil
}

type taskInfo struct {
	Name              string
	TaskDefinitionArn string
}

func DescribeTasks(region, cluster, task string) ([]taskInfo, error) {
	cfg, err := shared.LoadAWSConfig(region)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	// Create ECS client
	client := newECSClient(cfg)
	input := &ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   []string{task},
	}

	output, err := client.DescribeTasks(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("unable to describe tasks: %v", err)
	}

	var availableContainers []taskInfo

	for _, outputTask := range output.Tasks {
		taskDefinitionArn := outputTask.TaskDefinitionArn
		for _, container := range outputTask.Containers {
			availableContainers = append(availableContainers, taskInfo{
				TaskDefinitionArn: *taskDefinitionArn,
				Name:              *container.Name,
			})
		}
	}

	return availableContainers, nil
}

func DescribeTaskDefinition(region, taskDefinition string) (string, error) {
	cfg, err := shared.LoadAWSConfig(region)
	if err != nil {
		return "", fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	client := newECSClient(cfg)
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &taskDefinition,
	}

	output, err := client.DescribeTaskDefinition(context.Background(), input)
	if err != nil {
		return "", fmt.Errorf("unable to describe task definition: %v", err)
	}

	outputJSON, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return "", fmt.Errorf("unable to marshal task definition output: %v", err)
	}

	return string(outputJSON), nil
}

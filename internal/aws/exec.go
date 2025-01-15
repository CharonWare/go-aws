package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func loadAWSConfig(region string) (aws.Config, error) {
	return config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
}

func newECSClient(cfg aws.Config) *ecs.Client {
	return ecs.NewFromConfig(cfg)
}

func ListClusters(region string) (clusters []string, err error) {
	cfg, err := loadAWSConfig(region)
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

func ListServices(region, cluster string) (services []string, err error) {
	cfg, err := loadAWSConfig(region)
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

func ListTasks(region, cluster, service string) (tasks []string, err error) {
	cfg, err := loadAWSConfig(region)
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

func DescribeTasks(region, cluster, task string) (availableContainers []string, err error) {
	cfg, err := loadAWSConfig(region)
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

	for _, outputTask := range output.Tasks {
		for _, container := range outputTask.Containers {
			availableContainers = append(availableContainers, aws.ToString(container.Name))
		}
	}

	return availableContainers, nil
}

package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func ListClusters(region string) (clusters []string, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	// Create ECS client
	client := ecs.NewFromConfig(cfg)
	input := &ecs.ListClustersInput{}
	output, err := client.ListClusters(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("unable to list ECS clusters: %v", err)
	}

	clusters = append(clusters, output.ClusterArns...)
	return clusters, nil
}

func ListServices(region, cluster string) (services []string, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	// Create ECS client
	client := ecs.NewFromConfig(cfg)
	input := &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	}
	output, err := client.ListServices(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("unable to list ECS services: %v", err)
	}

	services = append(services, output.ServiceArns...)
	return services, nil
}

func ListTasks(region, cluster, service string) (tasks []string, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	// Create ECS client
	client := ecs.NewFromConfig(cfg)
	input := &ecs.ListTasksInput{
		Cluster:     aws.String(cluster),
		ServiceName: aws.String(service),
	}
	output, err := client.ListTasks(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("unable to list ECS tasks: %v", err)
	}

	tasks = append(tasks, output.TaskArns...)
	return tasks, nil
}

func DescribeTasks(region, cluster, task string) (availableContainers []string, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS configuration: %v", err)
	}

	// Create ECS client
	client := ecs.NewFromConfig(cfg)
	input := &ecs.DescribeTasksInput{
		Cluster: aws.String(cluster),
		Tasks:   []string{task},
	}

	output, err := client.DescribeTasks(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("unable to describe tasks: %v", err)
	}

	for _, task := range output.Tasks {
		for _, container := range task.Containers {
			availableContainers = append(availableContainers, aws.ToString(container.Name))
		}
	}

	return availableContainers, nil
}

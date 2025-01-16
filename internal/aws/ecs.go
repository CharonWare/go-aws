package aws

import (
	"context"
	"fmt"

	"github.com/CharonWare/go-aws/internal/shared"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

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

type describeCluster struct {
	Name           string
	ContainerHosts int32
	RunningTasks   int32
	PendingTasks   int32
	Services       int32
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

	output, err := client.DescribeClusters(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("unable to describe cluster: %v", err)
	}

	var chosenCluster []describeCluster

	for _, clusters := range output.Clusters {
		chosenCluster = append(chosenCluster, describeCluster{
			Name:           *clusters.ClusterName,
			ContainerHosts: clusters.RegisteredContainerInstancesCount,
			RunningTasks:   clusters.RunningTasksCount,
			Services:       clusters.ActiveServicesCount,
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

type describeService struct {
	Name       string
	Desired    int32
	Running    int32
	Pending    int32
	LaunchType string
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

	output, err := client.DescribeServices(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("unable to describe service: %v", err)
	}

	var chosenService []describeService

	for _, i := range output.Services {
		chosenService = append(chosenService, describeService{
			Name:       *i.ServiceName,
			Desired:    i.DesiredCount,
			Running:    i.RunningCount,
			Pending:    i.PendingCount,
			LaunchType: string(i.LaunchType),
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

func DescribeTasks(region, cluster, task string) (availableContainers []string, err error) {
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

	for _, outputTask := range output.Tasks {
		for _, container := range outputTask.Containers {
			availableContainers = append(availableContainers, aws.ToString(container.Name))
		}
	}

	return availableContainers, nil
}

// func DescribeService(region, cluster, service string) error {
// 	cfg, err := shared.LoadAWSConfig(region)
// 	if err != nil {
// 		return err
// 	}

// 	client := newECSClient(cfg)
// 	input := &ecs.DescribeServicesInput{
// 		Cluster: aws.String(cluster),
// 		Services: []string{service},
// 	}

// 	output, err := client.DescribeServices(context.Background(), input)
// 	if err != nil {
// 		return fmt.Errorf("unable to describe ECS service: %v", err)
// 	}

// 	for outputService := range output.Services {}
// }

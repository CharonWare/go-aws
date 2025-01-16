/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/CharonWare/go-aws/internal/aws"
	"github.com/CharonWare/go-aws/internal/ui"
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var ecsCmd = &cobra.Command{
	Use:   "ecs",
	Short: "Start an ECS Exec session with a specified container",
	Long: `The exec command will present you a series of interactable and searchable menus
	that will allow you to select a cluster, service, task, and finally a container which
	will then be exec'd into. Use with: go-aws ecs`,
	RunE: func(cmd *cobra.Command, args []string) error {
		region := os.Getenv("AWS_DEFAULT_REGION")
		if region == "" {
			region = "eu-west-1" // Default region if the environment variable is not set
		}

		// Search for available ECS clusters in the chosen region
		clusters, err := aws.ListClusters(region)
		if err != nil {
			return err
		}

		// Prompt the user to select a cluster
		i, _, err := ui.CreatePrompt(clusters, "Select a cluster:")
		if err != nil {
			return err
		}

		selectedCluster := clusters[i]

		// Check if the describe-cluster flag is set and proceed based on that
		describeCluster, _ := cmd.Flags().GetBool("describe-cluster")
		if describeCluster {
			output, err := aws.DescribeCluster(region, selectedCluster)
			if err != nil {
				return err
			}
			for _, clusterInfo := range output {
				fmt.Printf(`
Name:            %s
Container Hosts: %d
Running Tasks:   %d
Pending Tasks:   %d
Services:        %d
`,
					clusterInfo.Name,
					clusterInfo.ContainerHosts,
					clusterInfo.RunningTasks,
					clusterInfo.PendingTasks,
					clusterInfo.Services,
				)
				os.Exit(0)
			}
		}

		// Pass the selected cluster to a list services call to see all services in that cluster
		services, err := aws.ListServices(region, selectedCluster)
		if err != nil {
			return err
		}

		// Prompt the user to select a service
		ii, _, err := ui.CreatePrompt(services, "Select a service:")
		if err != nil {
			return err
		}

		selectedService := services[ii]

		// Check if the describe-service flag is set and proceed based on that
		describeService, _ := cmd.Flags().GetBool("describe-service")
		if describeService {
			output, err := aws.DescribeService(region, selectedCluster, selectedService)
			if err != nil {
				return err
			}
			for _, serviceInfo := range output {
				fmt.Printf(`
Name:       %s
Desired:    %d
Running:    %d
Pending:    %d
LaunchType: %s
`,
					serviceInfo.Name,
					serviceInfo.Desired,
					serviceInfo.Running,
					serviceInfo.Pending,
					serviceInfo.LaunchType,
				)
				os.Exit(0)
			}
		}

		// if service flag is set, stop here and describe this service

		// Pass the selected cluster and service to a list tasks call to see all tasks in that service
		tasks, err := aws.ListTasks(region, selectedCluster, selectedService)
		if err != nil {
			return err
		}

		iii, _, err := ui.CreatePrompt(tasks, "Select a task:")
		if err != nil {
			return err
		}

		selectedTask := tasks[iii]

		// if task flag is set, stop here and describe this task

		// Tasks can have multiple containers so we need to describe them to find the container names
		containers, err := aws.DescribeTasks(region, selectedCluster, selectedTask)
		if err != nil {
			return err
		}

		// If there are no containers, return an error
		if len(containers) == 0 {
			return fmt.Errorf("no containers available for the selected task")
		}

		// If only one container exists, skip the prompt
		if len(containers) == 1 {
			return execToContainer(region, selectedCluster, selectedTask, containers[0])
		}

		// If there are multiple containers, prompt the user for selection
		iiii, _, err := ui.CreatePrompt(containers, "Select a container:")
		if err != nil {
			return err
		}

		selectedContainer := containers[iiii]
		return execToContainer(region, selectedCluster, selectedTask, selectedContainer)

	},
}

func init() {
	rootCmd.AddCommand(ecsCmd)

	ecsCmd.Flags().Bool("describe-cluster", false, "Describe the selected cluster")
	ecsCmd.Flags().Bool("describe-service", false, "Describe the selected service")
	ecsCmd.Flags().Bool("describe-task", false, "Describe the selected task")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func execToContainer(region, cluster, taskArn, container string) error {
	cmd := exec.Command("aws", "ecs", "execute-command",
		"--cluster", cluster,
		"--task", taskArn,
		"--container", container,
		"--interactive",
		"--command", "/bin/bash",
		"--region", region,
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Starting ECS Exec session for task: %s\n", taskArn)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start ECS Exec session: %w", err)
	}
	return nil
}

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/CharonWare/go-aws/internal/aws"
	"github.com/CharonWare/go-aws/internal/ui"
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var ecsCmd = &cobra.Command{
	Use:   "ecs",
	Short: "Start an ECS Exec session with a specified container, or describe ECS resources",
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
		describeClusterBool, _ := cmd.Flags().GetBool("describe-cluster")
		if describeClusterBool {
			describeCluster(region, selectedCluster)
			os.Exit(0)
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
		describeServiceBool, _ := cmd.Flags().GetBool("describe-service")
		if describeServiceBool {
			describeService(region, selectedCluster, selectedService)
			os.Exit(0)
		}

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

		// Tasks can have multiple containers so we need to describe them to find the container names
		containers, err := aws.DescribeTasks(region, selectedCluster, selectedTask)
		if err != nil {
			return err
		}

		// If there are no containers, return an error
		if len(containers) == 0 {
			return fmt.Errorf("no containers available for the selected task")
		}

		// If there are multiple containers, prompt the user for selection
		containerNames := make([]string, len(containers))
		for i, container := range containers {
			containerNames[i] = container.Name
		}

		// if task-definition flag is set, stop here and describe the task-definition for this task
		taskDefinitionBool, _ := cmd.Flags().GetBool("task-definition")
		if taskDefinitionBool {
			output, err := aws.DescribeTaskDefinition(region, containers[0].TaskDefinitionArn)
			if err != nil {
				return err
			}
			fmt.Println(output)
			os.Exit(0)
		}

		// // If only one container exists, skip the prompt
		if len(containerNames) == 1 {
			return execToContainer(region, selectedCluster, selectedTask, containers[0].Name)
		}

		// Otherwise, select a container to exec to
		iiii, _, err := ui.CreatePrompt(containerNames, "Select a container:")
		if err != nil {
			return err
		}

		selectedContainer := containers[iiii]

		return execToContainer(region, selectedCluster, selectedTask, selectedContainer.Name)

	},
}

func init() {
	rootCmd.AddCommand(ecsCmd)

	ecsCmd.Flags().BoolP("describe-cluster", "c", false, "Describe the selected cluster")
	ecsCmd.Flags().BoolP("describe-service", "s", false, "Describe the selected service")
	ecsCmd.Flags().BoolP("task-definition", "t", false, "Show the task definition for the selected task")

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

	// Create a signal channel to forward SIGINT / Ctrl+C to the container
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		for sig := range signalChannel {
			if sig == os.Interrupt {
				// Forward SIGINT to subprocess
				if cmd.Process != nil {
					_ = cmd.Process.Signal(os.Interrupt)
				}
			}
		}
	}()

	fmt.Printf("Starting ECS Exec session for task: %s\n", taskArn)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ECS Exec session: %w", err)
	}

	// Wait for the command to finish
	err := cmd.Wait()

	// Stop listening for signals after the command exits
	signal.Stop(signalChannel)
	close(signalChannel)

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && !exitError.Success() {
			return fmt.Errorf("container session exited with error: %v", exitError)
		}
		return fmt.Errorf("failed to execute command: %w", err)
	}
	return nil
}

func describeCluster(region, cluster string) error {
	output, err := aws.DescribeCluster(region, cluster)
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
AVG CPU (5 min): %.1f%%
`,
			clusterInfo.Name,
			clusterInfo.ContainerHosts,
			clusterInfo.RunningTasks,
			clusterInfo.PendingTasks,
			clusterInfo.Services,
			clusterInfo.AVGCPU,
		)
	}
	return nil
}

func describeService(region, cluster, service string) error {
	output, err := aws.DescribeService(region, cluster, service)
	if err != nil {
		return err
	}
	for _, serviceInfo := range output {
		fmt.Printf(`
Name:            %s
Desired:         %d
Running:         %d
Pending:         %d
LaunchType:      %s
AVG CPU (5 min): %.1f%%
`,
			serviceInfo.Name,
			serviceInfo.Desired,
			serviceInfo.Running,
			serviceInfo.Pending,
			serviceInfo.LaunchType,
			serviceInfo.AVGCPU,
		)
	}
	return nil
}

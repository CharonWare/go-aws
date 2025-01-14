/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/CharonWare/go-aws/internal/aws"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Start an ECS Exec session with a specified container",
	RunE: func(cmd *cobra.Command, args []string) error {
		region := os.Getenv("AWS_DEFAULT_REGION")
		if region == "" {
			region = "eu-west-1" // Default region if the environment variable is not set
		}
		clusters, err := aws.ListClusters(region)
		if err != nil {
			return err
		}

		// fmt.Println("Clusters:")
		// for i, cluster := range clusters {
		// 	fmt.Printf("%d. %s\n", i+1, cluster)
		// }

		prompt := promptui.Select{
			Label: "Select a cluster",
			Items: clusters,
			Searcher: func(input string, index int) bool {
				item := clusters[index]
				return containsIgnoreCase(item, input)
			},
		}

		i, _, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}

		selectedCluster := clusters[i]

		services, err := aws.ListServices(region, selectedCluster)
		if err != nil {
			return err
		}

		prompt2 := promptui.Select{
			Label: "Select a service",
			Items: services,
			Searcher: func(input string, index int) bool {
				item := services[index]
				return containsIgnoreCase(item, input)
			},
		}

		ii, _, err := prompt2.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}

		selectedService := services[ii]

		tasks, err := aws.ListTasks(region, selectedCluster, selectedService)
		if err != nil {
			return err
		}

		prompt3 := promptui.Select{
			Label: "Select a task",
			Items: tasks,
			Searcher: func(input string, index int) bool {
				item := tasks[index]
				return containsIgnoreCase(item, input)
			},
		}

		iii, _, err := prompt3.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}

		selectedTask := tasks[iii]

		containers, err := aws.DescribeTasks(region, selectedCluster, selectedTask)
		if err != nil {
			return err
		}

		prompt4 := promptui.Select{
			Label: "Select a container",
			Items: containers,
			Searcher: func(input string, index int) bool {
				item := containers[index]
				return containsIgnoreCase(item, input)
			},
		}

		iiii, _, err := prompt4.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}

		selectedContainer := containers[iiii]

		return execToContainer(region, selectedCluster, selectedTask, selectedContainer)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func execToContainer(region, cluster, taskArn, container string) error {
	cmd := exec.Command("aws", "ecs", "execute-command", "--cluster", cluster, "--task", taskArn, "--container", container, "--interactive", "--command", "/bin/bash", "--region", region)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Starting ECS Exec session for task: %s\n", taskArn)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start ECS Exec session: %w", err)
	}
	return nil
}

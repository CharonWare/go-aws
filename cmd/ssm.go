/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/CharonWare/go-aws/internal/aws"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// ssmCmd represents the ssm command
var ssmCmd = &cobra.Command{
	Use:   "ssm",
	Short: "Start an SSM session with an EC2 instance",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 && args[0] != "" {
			// If an instance ID is provided then start the SSM session directly
			return startSSMSession(args[0])
		}

		region := os.Getenv("AWS_DEFAULT_REGION")
		if region == "" {
			region = "eu-west-1" // Default region if the environment variable is not set
		}

		instances, err := aws.ListEC2Instances(region)
		if err != nil {
			return fmt.Errorf("error listing EC2 instances: %w", err)
		}

		if len(instances) == 0 {
			fmt.Println("No EC2 instances found")
			return nil
		}

		var options []string
		instanceMap := make(map[string]string)
		for _, inst := range instances {
			name := inst.Name
			if name == "" {
				name = "{no name}"
			}
			options = append(options, name)
			instanceMap[name] = inst.ID
		}

		prompt := promptui.Select{
			Label: "Select an EC2 instance",
			Items: options,
			Searcher: func(input string, index int) bool {
				item := options[index]
				return containsIgnoreCase(item, input)
			},
		}

		i, _, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}

		selectedName := options[i]
		selectedID := instanceMap[selectedName]

		fmt.Printf("%s chosen.\n", selectedName)
		return startSSMSession(selectedID)
	},
}

func containsIgnoreCase(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

func init() {
	rootCmd.AddCommand(ssmCmd)
	//ssmCmd.Flags().String("region", "", "AWS region to query (e.g., us-east-1)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ssmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ssmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func startSSMSession(instanceID string) error {
	cmd := exec.Command("aws", "ssm", "start-session", "--target", instanceID)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Starting SSM session for instance: %s\n", instanceID)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start SSM session: %w", err)
	}
	return nil
}

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/CharonWare/go-aws/internal/aws"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List EC2 instances in this account",
	RunE: func(cmd *cobra.Command, args []string) error {
		region := os.Getenv("AWS_DEFAULT_REGION")
		if region == "" {
			region = "eu-west-1" // Default region if the environment variable is not set
		}
		instances, err := aws.ListEC2Instances(region)
		if err != nil {
			return err
		}

		fmt.Println("Instances:")
		for i, instance := range instances {
			fmt.Printf("%d. %s\n", i+1, instance)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

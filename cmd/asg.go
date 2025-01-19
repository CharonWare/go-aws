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

// asgCmd represents the asg command
var asgCmd = &cobra.Command{
	Use:   "asg",
	Short: "Describes the scaling values of the ASGs in the current account and region",
	Long: `Provides Name, MinSize, MaxSize and DesiredCapacity for each autoscaling
	group in the current account and region.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		region := os.Getenv("AWS_DEFAULT_REGION")
		if region == "" {
			region = "eu-west-1" // Default region if the environment variable is not set
		}

		groups, err := aws.DescribeASGs(region)
		if err != nil {
			return fmt.Errorf("error describing ASGs: %w", err)
		}

		if len(groups) == 0 {
			fmt.Println("No ASGs found")
		}

		for _, asg := range groups {
			fmt.Printf(`

Name:            %s
MinSize:         %d
MaxSize:         %d
DesiredCapacity: %d
AVG CPU (5 min): %.1f%%

`,
				asg.Name,
				asg.MinSize,
				asg.MaxSize,
				asg.DesiredCapacity,
				asg.AVGCPU,
			)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(asgCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// asgCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// asgCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

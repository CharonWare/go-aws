# go-aws

`go-aws` makes it easier to use AWS SSM Session Manager and ECS Exec on the cli. These commands can be tedious because of needing to find out instance IDs and ECS task IDs in order to use them. This tool provides an interactive menu to select your EC2 instance, or your cluster->service->task->container (if there is only one container in your task it will skip this prompt).

This tool is still being improved.

## Usage
* Currently the tool defaults to `eu-west-1` but this will be overridden if you export your region with the `AWS_DEFAULT_REGION` variable
* `go-aws ssm` will present you with a list of EC2 instances, once selected it will attempt to connect you to this instance via a session manager connection
* If you have an instance ID already you can pass it into the ssm argument, `go-aws ssm i-112233445566`
* `go-aws exec` will allow you to select a cluster, service, task, and then container in that task, it will then attempt to connect via ECS Exec
* `go-aws asg` will describe the autoscaling groups in the current account and region, providing the name, MinSize, MaxSize, and DesirecCapcity

## Requirements
* go 1.23.1 or higher
* AWS account credentials
* Your EC2 instances must be capable of receiving session manager connections
* Your ECS tasks must be set up for ECS Exec

## Installation
1) Clone the repo
2) `go build .`
3) Ensure the `go-aws` binary is executable
4) Use the tool with `./go-aws` OR
5) `sudo mv go-aws /usr/local/bin` - assuming this is in your $PATH you can now use the tool with `go-aws`
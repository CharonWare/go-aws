# go-aws

`go-aws` makes it easier to use AWS SSM Session Manager and ECS Exec on the cli. These commands can be tedious because of needing to find out instance IDs and ECS task IDs in order to use them. This tool provides an interactive menu to select your EC2 instance, or your cluster->service->task->container (if there is only one container in your task it will skip this prompt).

This tool is still being improved.

## Usage
* Run the tool without arguments to see a breakdown of available commands
* Append `-h` to a command to get a help menu with a more detailed breakdown of the command and any available flags

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
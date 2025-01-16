# go-aws

`go-aws` provides an alternative to the AWS CLI and speeds up administration by utilising the AWS Go SDK. Session manager and ECS Exec connections can be initiated using `go-aws` without having to run multiple queries or recall all the details of your infrastructure. `go-aws` will prompt you with interactive menus that will allow you to select which clusters, services, autoscaling groups, instances etc to connect to, describe, or update.

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
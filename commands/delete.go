package commands

import (
	"log"

	"github.com/Skarlso/go-furnace/config"
	"github.com/Skarlso/go-furnace/utils"
	"github.com/Yitsushi/go-commander"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/fatih/color"
)

// Delete command.
type Delete struct {
}

// Execute defines what this command does.
func (c *Delete) Execute(opts *commander.CommandHelper) {

	stackname := opts.Arg(0)
	if len(stackname) < 1 {
		stackname = config.STACKNAME
	}

	cyan := color.New(color.FgCyan).SprintFunc()

	log.Printf("Deleting CloudFormation stack with name: %s\n", cyan(stackname))
	sess := session.New(&aws.Config{Region: aws.String("eu-central-1")})
	cfClient := cloudformation.New(sess, nil)
	params := &cloudformation.DeleteStackInput{
		StackName: aws.String(stackname),
	}
	_, err := cfClient.DeleteStack(params)
	utils.CheckError(err)
	describeStackInput := &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackname),
	}
	utils.WaitForFunctionWithStatusOutput("DELETE_COMPLETE", config.WAITFREQUENCY, func() {
		cfClient.WaitUntilStackDeleteComplete(describeStackInput)
	})
}

// NewDelete Creates a new Delete command.
func NewDelete(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &Delete{},
		Help: &commander.CommandDescriptor{
			Name:             "delete",
			ShortDescription: "Delete a stack",
			LongDescription:  `Delete a stack with a given name.`,
			Arguments:        "name",
			Examples:         []string{"delete", "delete MyStackName"},
		},
	}
}
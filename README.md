# Furnace

![Logo](img/logo.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/Skarlso/go-furnace)](https://goreportcard.com/report/github.com/Skarlso/go-furnace) [![Build Status](https://travis-ci.org/Skarlso/go-furnace.svg?branch=master)](https://travis-ci.org/Skarlso/go-furnace)

## Intro

AWS Cloud Formation hosting with Go. This project utilises the power of AWS CloudFormation and CodeDeploy in order to
simply deploy an application into a robust, self-healing, redundant environment. The environment is configurable through
the CloudFormation Template. A sample can be found in the `templates` folder.

The application to be deployed is currently handled via GitHub, but later on, S3 based deployment will also be supported.

A sample application is provider under the `furnace-codedeploy-app` folder.

## AWS

### CloudFormation

[CloudFormation](https://aws.amazon.com/cloudformation/) as stated in the AWS documentation is an
> ...easy way to create and manage a collection of related AWS resources, provisioning and updating them in an orderly and predictable fashion.

Meaning, that via a template file it is possible to provide a description of the environment we would like to launch
are application into. How many server we would like to have? Load Balancing, and Auto Scaling setup. Own, isolated
network with VPCs. CloudFormation brings all these elements together into a bundler project called a `Stack`.
This stack can be created, updated, deleted and queried for various information.

This is what `Furnace` aims ti abstract in order to provide a very easy interface to work with complex architecture.

### CodeDeploy

[CodeDeploy](http://docs.aws.amazon.com/codedeploy/latest/userguide/welcome.html), as the documentation states
> ...coordinates application deployments to Amazon EC2 instances

In short, once the stack is up, we would like to deploy our application to the stack for usage. CodeDeploy takes care of that.
We don't have to scp something to our instances, we don't have to care if an instance goes away, or if we would like to have
a copy of that same instance. CodeDeploy can be integrated with various other services, so once we described how to deploy
our application, we never have to worry about it again. A simple `furnace push` will install our app to every instance that
the `Stack` manages.

Don't forget to install the CodeDeploy agent to your instances for the CodeDeploy to work. For this, see an example in the
provided template.

## Go

The decision to use [Go](https://golang.org/) for this project came very easy considering the many benefits Go provides when
handling APIs and async requests. Downloading massive files from S3 in threads, or starting many services at once is a breeze.
Go also provides a single binary which is easy to put on the execution path and use `Furnace` from any location.

Go has ample libraries which come to aid with AWS and their own Go SDK is pretty mature and stable.

## Usage

### Make

This project is using a `Makefile` for it's build processes. The following commands will create a binary and
run all tests:

```bash
make build test
```

`make install` will install the binary in your `$GOHOME\bin` path. Though the binary will be named `go-furnace`.

For other targets, please consult the Makefile.

### Configuration

Furnace requires three environment properties. These are:

```bash
# git revision is the latest git commit to use
export FURNACE_GIT_REVISION=b80ea5b9dfefcd21e27a3e0f149ec73519d5a6f1
export FURNACE_GIT_ACCOUNT=skarlso/furnace-codedeploy-app
export FURNACE_REGION=eu-central-1
```

Furnace also requires the CloudFormation template to be placed under `~/.config/go-furnace`.

CodeDeploy further requires a IAM policy to the current user in order to be able to handle ASG and deploying to the EC2 instances.
For this, a regular IAM role can be created from the AWS console. The name of the IAM profile can be configured later when pushing,
if that is not set the default is used which is `CodeDeployServiceRole`. This default can be found under `config.CODEDEPLOYROLE`.

The CloudFormation template needs also be copied under `~/.config/go-furnace`.

### Commands

Furnace provides the following commands (which you can check by running `./furnace`):

```bash
➜  go-furnace git:(master) ✗ ./furnace
create name               Create a stack
delete name               Delete a stack
status name               Status of a stack.
push name                 Push to stack
delete-application name   Deletes an Application
help [command]            Display this help or a command specific help
```

Create and Delete will wait for the these actions to complete via a Waiter function. The waiters spinner type
can be set via the env property `FURNACE_SPINNER`. This is optional. The following spinners are available:

```go
// Spinners is a collection os spinner types
var Spinners = []string{`←↖↑↗→↘↓↙`,
	`▁▃▄▅▆▇█▇▆▅▄▃`,
	`┤┘┴└├┌┬┐`,
	`◰◳◲◱`,
	`◴◷◶◵`,
	`◐◓◑◒`,
	`⣾⣽⣻⢿⡿⣟⣯⣷`,
	`|/-\`}
```

The spinner defaults to `|/-\` which is # 7.

#### create

This will create the whole stack via the configuration provided under templates.

![Create](./img/create.png)

As you can see, furnace will ask for the parameters that reside in a template. If default is desired, simply
hit enter to continue using the default value.

#### delete

Deletes the whole stack complete with everything attached to the stack expect for the CodeDeploy application.

![Delete](./img/delete.png)

#### push

This is the command to get your application to be deployed onto all of your configured instances. This works via
two things. AutoScaling groups provided by the CloudFormation stack plus Tags that are put onto the instances called
`fu_stage`.

![Push](./img/push.png)

#### delete-application

Will delete the application and the deployment group completely.

#### successful push

If you are using the provided example and everything works, you should see the following output once you visit the
url provided by the load balancer.

![Success](./img/push_success.png)

## Plugins

Until Go's own Plugin system is fully supported ( which will take a while ), a rudimentary plugin system has been put in place.
There are four events currently for plugins:
```bash
- PRE_CREATE
- POST_CREATE
- PRE_DELETE
- POST_DELETE
```

In order to implement a plugin, place a file into the `plugins` folder, and implement the following interface:

```go
// Plugin interface defines the capabilities of a plugin
type Plugin interface {
	RunPlugin()
}
```

At the moment the plugin also needs to be registered by hand in `furnace.go` like this:

```go
// For now, the including of a plugin is done manually.
preCreatePlugins := []plugins.Plugin{
    plugins.MyAwesomePreCreatePlugin{Name: "SamplePreCreatePlugin"},
}
postCreatePlugins := []plugins.Plugin{
    plugins.MyAwesomePostCreatePlugin{Name: "SamplePostCreatePlugin"},
}
plugins.RegisterPlugin(config.PRECREATE, preCreatePlugins)
plugins.RegisterPlugin(config.POSTCREATE, postCreatePlugins)
```

This will later be replaced by putting a .so or .dylib file into the plugin folder, and no re-compile will be necessary.

The plugin system also needs some way to pass in control over the current environment. So it's very much under development.

## Testing

Testing the project for development is simply by executing `make test`.

## Contributions

Contributions are very welcomed, ideas, questions, remarks, please don't hesitate to submit a ticket. On what to do,
please take a look at the ROADMAP.md file.

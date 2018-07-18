# slack-on-call-badge

This program was designed to automatically set Slack chat status badges for designated users if they are on-call in an OpsGenie's schedule.

It is designed to run indefinitely inside docker container.

## Disclaimer

This is an alpha release. Use it on your own risk. If you liked it and/or see how it can be improved, please help me with your Pull Request. Thank you!

## Installation

### You can download a prepared image from Docker Hub

```shell
docker pull dragonsmith/slack-on-call-badge:0.0.3
```

### Or you can build a binary from sources

```shell
go get github.com/dragonsmith/slack-on-call-badge
cd $GOPATH/github.com/dragonsmith/slack-on-call-badge
make build
```

## Configuration

The configuration is done via ENV variables.

Required ones:

* `SLACK_TOKEN` - Slack API token with admin rights.
* `OPSGENIE_TOKEN` - OpsGenie API token.
* `OPSGENIE_ROTATION` - OpsGenie rotation name which should be used for reference.
* `ADMINS` - List of people we want to track. Example: `User1_OpsGenie_email:User1_Slack_id,User2_OpsGenie_email:User2_Slack_id,...`

And optional:

* `ON_CALL_ICON` - Slack icon name to use as a status icon. Default: `:on_call:`
* `ON_CALL_TEXT` - Slack status text to use. Default: `on call`

Command line options for debugging & testing:

* `--once` - Make program to check and set status badges once and exit.
* `--debug` - Make output more verbose.
* `--dru-run` - Do not actually set Slack status.

## Example

To run the program inside docker:

```shell
docker run --name slack-on-call-badge \
 -e SLACK_TOKEN=changeme \
 -e OPSGENIE_TOKEN=changeme \
 -e OPSGENIE_ROTATION=changeme \
 -e ADMINS="User1_OpsGenie_email:User1_Slack_id,User2_OpsGenie_email:User2_Slack_id" dragonsmith/slack-on-call-badge:0.0.3
```

To see its logs:

```shell
docker logs slack-on-call-badge
```

## Sponsor

[![Sponsored by Evil Martians](https://evilmartians.com/badges/sponsored-by-evil-martians.png)](https://evilmartians.com)

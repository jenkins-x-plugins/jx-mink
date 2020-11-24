# jx mink

[![Documentation](https://godoc.org/github.com/jenkins-x-plugins/jx-mink?status.svg)](https://pkg.go.dev/mod/github.com/jenkins-x-plugins/jx-mink)
[![Go Report Card](https://goreportcard.com/badge/github.com/jenkins-x-plugins/jx-mink)](https://goreportcard.com/report/github.com/jenkins-x-plugins/jx-mink)
[![Releases](https://img.shields.io/github/release-pre/jenkins-x/helmboot.svg)](https://github.com/jenkins-x-plugins/jx-mink/releases)
[![LICENSE](https://img.shields.io/github/license/jenkins-x/helmboot.svg)](https://github.com/jenkins-x-plugins/jx-mink/blob/master/LICENSE)
[![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](https://slack.k8s.io/)

`jx-mink` is a simple command line tool for using [mink](https://github.com/mattmoor/mink) with Jenkins X Pipelines to perform image builds and resolve image references in helm charts.


## Getting Started

Download the [jx-mink binary](https://github.com/jenkins-x-plugins/jx-mink/releases) for your operating system and add it to your `$PATH`.

## Configuring Kaniko

If you want to use vanilla Kaniko then you can use the following environment variables:

* **PLAIN_KANIO** set to `true` or `yes` to enable vanilla kaniko rather than using the mink logic to discover the container images to build and optionally YAML to resolve
* **KANIKO_CONTEXT** the context path for the build which defaults to the current directory. So this will typically be `/workspace/source` inside a tekton pipeline or the current working directory if invoked locally
* **KANIKO_DOCKERFILE** the location of the `Dockerfile` which is usually `Dockerfile` inside the context (see above)

## Commands

See the [jx-mink command reference](https://github.com/jenkins-x-plugins/jx-mink/blob/master/docs/cmd/jx-mink.md#jx-mink)


/**
qd (QuickDeploy)

Quickly deploy an image to K8s as a deployment.

usage:
qd run alpine:latest
qd exec ubuntu:latest
qd list
qd stop
**/

package main

import (
	"qd/cmd"
)

func main() {
	cmd.Run()
}

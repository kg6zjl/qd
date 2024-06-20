package docker

import (
	"fmt"
	"os/exec"
	"qd/utils"
)

func build(imageName string) (string, error) {
	cmd := exec.Command("podman", "build", "-t", imageName, ".")
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error: %v", err)
	}
	return imageName, nil
}

func BuildDeploy() (string, error) {
	name := utils.UniqName("docker-deploy")
	imageName, err := build(name)
	if err != nil {
		fmt.Println(err)
		return "", err
	} else {
		fmt.Printf("Docker image %s built successfully!\n", imageName)
		return imageName, nil
	}
}

package cmd

import (
	"fmt"
	"log"
	"qd/docker"
	"qd/kube"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

var kubeconfig string

var Version string

var rootCmd = &cobra.Command{
	Use:   "qd",
	Short: "qd is a quick deploy tool for Kubernetes",
	Long:  `qd is a CLI tool to quickly deploy images to Kubernetes`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Return qd version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version v%s", Version)
	},
}

var runCmd = &cobra.Command{
	Use:   "run [image:tag]",
	Short: "Run a new deployment",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// read config
		config := kube.ReadKubeConfig(kubeconfig)
		// setup client
		clientset, _ := kubernetes.NewForConfig(config)
		// get current context namespace
		namespace := kube.CurrentNamespace()
		// setup deployment
		deploymentClient := clientset.AppsV1().Deployments(namespace)
		_ = kube.Run(deploymentClient, args[0])
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all quickdeploy deployments",
	Run: func(cmd *cobra.Command, args []string) {
		// read config
		config := kube.ReadKubeConfig(kubeconfig)
		// setup client
		clientset, _ := kubernetes.NewForConfig(config)
		// get current context namespace
		namespace := kube.CurrentNamespace()
		// setup deployment
		deploymentClient := clientset.AppsV1().Deployments(namespace)
		kube.List(deploymentClient)
	},
}

var buildDeployCmd = &cobra.Command{
	Use:   "build",
	Short: "Docker build and then deploy",
	Run: func(cmd *cobra.Command, args []string) {
		// read config
		config := kube.ReadKubeConfig(kubeconfig)
		// setup client
		clientset, _ := kubernetes.NewForConfig(config)
		// get current context namespace
		namespace := kube.CurrentNamespace()
		// setup deployment
		deploymentClient := clientset.AppsV1().Deployments(namespace)
		imageName, err := docker.BuildDeploy()
		if err != nil {
			log.Fatalf("Failed to build Docker image: %s", err)
		}
		kube.Run(deploymentClient, imageName)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all quickdeploy deployments",
	Run: func(cmd *cobra.Command, args []string) {
		// read config
		config := kube.ReadKubeConfig(kubeconfig)
		// setup client
		clientset, _ := kubernetes.NewForConfig(config)
		// get current context namespace
		namespace := kube.CurrentNamespace()
		// setup deployment
		deploymentClient := clientset.AppsV1().Deployments(namespace)
		kube.Stop(deploymentClient)
	},
}

var execCmd = &cobra.Command{
	Use:   "exec [image:tag]",
	Short: "Deploy and exec into a pod",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// read config
		config := kube.ReadKubeConfig(kubeconfig)
		// setup client
		clientset, _ := kubernetes.NewForConfig(config)
		// get current context namespace
		namespace := kube.CurrentNamespace()
		// setup deployment
		deploymentClient := clientset.AppsV1().Deployments(namespace)

		kube.RunAndExec(deploymentClient, args[0], namespace, config, clientset)
	},
}

func Run() {
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "~/.kube/config", "absolute path to the kubeconfig file")
	rootCmd.AddCommand(versionCmd, runCmd, listCmd, stopCmd, execCmd, buildDeployCmd)
	rootCmd.Execute()
}

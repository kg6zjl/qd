package kube

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"time"

	"qd/utils"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	appstypedv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

func Stop(deploymentClient appstypedv1.DeploymentInterface) {
	deployments, _ := deploymentClient.List(context.TODO(), v1.ListOptions{})
	for _, deployment := range deployments.Items {
		if deployment.Annotations["quickdeploy"] == "true" {
			deploymentClient.Delete(context.TODO(), deployment.Name, v1.DeleteOptions{})
			fmt.Printf("Stopped deployment %s\n", deployment.Name)
		}
	}
}

func List(deploymentClient appstypedv1.DeploymentInterface) {
	found := false
	deployments, _ := deploymentClient.List(context.TODO(), v1.ListOptions{})
	for _, deployment := range deployments.Items {
		if deployment.Annotations["quickdeploy"] == "true" {
			fmt.Printf("%s\n", deployment.Name)
			found = true
		}
	}
	if !found {
		fmt.Println("No qd deployments found in namespace")
	}
}

func getPodsinDeploy(clientset *kubernetes.Clientset, namespace string, deployName string) ([]string, error) {
	// get the pods in the given namespace
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var podNames []string
	for _, pod := range pods.Items {
		if pod.Labels["qd-app"] == deployName {
			podNames = append(podNames, pod.Name)
		}
	}

	return podNames, nil
}

func waitForDeploymentReady(clientset *kubernetes.Clientset, deployName string, namespace string) error {
	fmt.Printf("Waiting for %s deployment to be created\n", deployName)
	for {
		deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deployName, v1.GetOptions{})
		if err != nil {
			return err
		}

		if *deployment.Spec.Replicas == deployment.Status.ReadyReplicas {
			return nil
		}

		time.Sleep(time.Second * 2)
	}
}

func RunAndExec(deployclient appstypedv1.DeploymentInterface, image string, namespace string, restconfig *rest.Config, clientset *kubernetes.Clientset) error {
	// deploy
	deployName := Run(deployclient, image)

	// wait for pod to come up
	err := waitForDeploymentReady(clientset, deployName, namespace)
	if err != nil {
		log.Fatalf("Error waiting for pod to be ready: %s", err.Error())
	}

	podNames, err := getPodsinDeploy(clientset, namespace, deployName)
	if err != nil {
		return err
	}

	// exec into pod
	exec(clientset, podNames[0], namespace, []string{"/bin/sh"}, restconfig)

	return nil
}

func Run(deploymentClient appstypedv1.DeploymentInterface, image string, names ...string) string {
	var name string
	if len(names) > 0 {
		name = names[0]
	} else {
		// split the image name for use in deployment/pod name
		imageBase := strings.Split(image, ":")
		name = utils.UniqName(imageBase[0])
	}

	// build the deployment spec
	deployment := &appsv1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				"quickdeploy": "true", // this is how we identify what deployments/pods we have to clean up later
				"qd-app":      name,
				"app":         name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32Ptr(1),
			Selector: &v1.LabelSelector{
				MatchLabels: map[string]string{
					"app":         name,
					"qd-app":      name,
					"quickdeploy": "true",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						"app":         name,
						"qd-app":      name,
						"quickdeploy": "true",
					},
				},
				Spec: corev1.PodSpec{
					DNSPolicy: corev1.DNSClusterFirst,
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: image,
							Command: []string{
								"/bin/sh",
								"-c",
								"sleep infinity",
							},
						},
					},
				},
			},
		},
	}

	// deploy to current context cluster and namespace
	// TODO ns and cluster could be passed in via args if needed
	fmt.Println("Creating deployment...")
	result, err := deploymentClient.Create(context.TODO(), deployment, v1.CreateOptions{})
	if err != nil {
		log.Fatalf("Error creating deployment: %s", err.Error())
	}
	fmt.Printf("Created deployment %q\n", result.GetObjectMeta().GetName())

	return result.GetObjectMeta().GetName()
}

func exec(clientset *kubernetes.Clientset, podName string, namespace string, command []string, config *rest.Config) error {
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command: command,
			Stdin:   true,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:             os.Stdin,
		Stdout:            os.Stdout,
		Stderr:            os.Stderr,
		Tty:               true,
		TerminalSizeQueue: nil,
	})

	return err
}

func ReadKubeConfig(kubeconfig string) *rest.Config {
	// Check if KUBECONFIG environment variable is set
	if os.Getenv("KUBECONFIG") != "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	} else {
		// Expand the "~" to the actual home directory
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("Cannot get current user: %s", err.Error())
		}
		kubeconfig = strings.Replace(kubeconfig, "~", usr.HomeDir, 1)
	}

	// Check if the file exists
	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		log.Fatalf("Kubeconfig file does not exist: %s", kubeconfig)
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	return config
}

func CurrentNamespace() string {
	// Get the namespace from the current context
	namespace, _, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	).Namespace()
	if err != nil {
		log.Fatalf("Error getting namespace from current context: %s", err.Error())
	}

	return namespace
}

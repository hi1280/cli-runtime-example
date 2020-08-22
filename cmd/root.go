package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

var (
	configFlags          *genericclioptions.ConfigFlags
	resourceBuilderFlags *genericclioptions.ResourceBuilderFlags
)

var rootCmd = &cobra.Command{
	Use: "cli-runtime-example",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := configFlags.ToRESTConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		rawConfig, err := configFlags.ToRawKubeConfigLoader().RawConfig()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		namespace := rawConfig.Contexts[rawConfig.CurrentContext].Namespace
		if len(*configFlags.Namespace) > 0 {
			namespace = *configFlags.Namespace
		}
		if *resourceBuilderFlags.AllNamespaces {
			namespace = ""
		}
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		for _, pod := range pods.Items {
			fmt.Println(pod.Name)
		}
	},
}

// Execute is entrypoint
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	configFlags = genericclioptions.NewConfigFlags(true)
	resourceBuilderFlags = genericclioptions.NewResourceBuilderFlags()
	resourceBuilderFlags.WithAllNamespaces(false)
	configFlags.AddFlags(rootCmd.PersistentFlags())
	resourceBuilderFlags.AddFlags(rootCmd.PersistentFlags())
}

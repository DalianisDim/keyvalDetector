/*
Copyright © 2023 Dimitris Dalianis <dimitris@dalianis.gr>
This file is part of CLI application keyvalDetector
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var cfgFile string
var version bool
var defaultKubeConfigPath = "/.kube/config"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "keyvalDetector",
	Version: "0.0.1",
	Short:   "Scan your k8s cluster for unused ConfigMaps and Secrets",
	Long: `keyvalDetector will scan your Kubernetes cluster for
ConfigMaps and Secrets that are not used by Pods.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		keyvalDetector()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "Config file (default is $HOME/.keyvalDetector.yaml)")
	rootCmd.Flags().BoolVarP(&version, "version", "v", version, "Print the version and exit.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".keyvalDetector" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".keyvalDetector")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func keyvalDetector() error {

	// uses the current context in kubeconfig
	// path-to-kubeconfig -- for example, /root/.kube/config
	config, _ := clientcmd.BuildConfigFromFlags("", homedir.HomeDir()+defaultKubeConfigPath)

	// creates the clientset
	clientset, _ := kubernetes.NewForConfig(config)

	// List all the namespaces in the cluster.
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// Foreach namespace, get configmaps, secrets and pods.
	// Then check pod's mounts

	// Foreach namespace
	for _, namespace := range namespaces.Items {
		// Slice to store mounted ConfigMap names
		namespaceMountedConfigMaps := []string{}

		// Slice to store mounted Secrets names
		namespaceMountedSecrets := []string{}

		// List all pods in current namespace
		pods, err := clientset.CoreV1().Pods(namespace.GetName()).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}

		// List all configmaps in current namespace
		configmaps, err := clientset.CoreV1().ConfigMaps(namespace.GetName()).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}

		// List all secrets in current namespace
		secrets, err := clientset.CoreV1().Secrets(namespace.GetName()).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}

		// If there's no pod in current namespace any configmap/secret are not used (except defaults)
		if len(pods.Items) == 0 {

			// Print the configmaps' names
			for _, configmap := range configmaps.Items {
				if !isDefaultConfigMap(configmap.GetName()) {
					fmt.Println("ConfigMap \"" + configmap.GetName() + "\" in namespace \"" + namespace.GetName() + "\" not used")
				}
			}

			// Print the secrets' names.
			for _, secret := range secrets.Items {
				if !isDefaultSecret(secret.GetName()) {
					fmt.Println("Secret \"" + secret.GetName() + "\" in namespace \"" + namespace.GetName() + "\" not used")
				}
			}
		}

		// If there is any pod in current namespace
		// Iterate over the pods
		// Check each pod's Volumes and store the configmaps and secrets mounted names on a separate slice
		// Compare lists of configmaps/secrets with mounted ones
		if len(pods.Items) != 0 {
			// Foreach pod
			for _, pod := range pods.Items {
				// Foreach pod's volume
				for _, volume := range pod.Spec.Volumes {
					if volume.ConfigMap != nil {
						namespaceMountedConfigMaps = append(namespaceMountedConfigMaps, volume.ConfigMap.Name)
					}
					if volume.Secret != nil {
						namespaceMountedSecrets = append(namespaceMountedSecrets, volume.Secret.SecretName)
					}
				}
			}
		}

		// I need check if configmaps/secrets of this namespace
		// are in namespaceMountedConfigMaps/namespaceMountedSecrets

		for _, configmap := range configmaps.Items {
			if !contains(namespaceMountedConfigMaps, configmap.GetName()) {
				if !isDefaultConfigMap(configmap.GetName()) {
					println("ConfigMap \"" + configmap.GetName() + "\" in namespace \"" + namespace.GetName() + "\" not used")
				}
			}
		}

		for _, secret := range secrets.Items {
			if !isDefaultSecret(secret.GetName()) {
				if !contains(namespaceMountedSecrets, secret.GetName()) {
					println("Secret \"" + secret.GetName() + "\" in namespace \"" + namespace.GetName() + "\" not used")
				}
			}
		}

	} // END - Foreach namespace

	return nil
}

// func test(clientset kubernetes.Clientset) {

// }

func isDefaultConfigMap(configmap string) bool {
	defaultConfigMaps := []string{"kube-root-ca.crt", "foobar"}

	for _, value := range defaultConfigMaps {
		if configmap == value {
			return true
		}
	}
	return false
}

func isDefaultSecret(secret string) bool {
	defaultSecrets := []string{"foobar"}

	for _, value := range defaultSecrets {
		if secret == value {
			return true
		}
	}
	return false
}
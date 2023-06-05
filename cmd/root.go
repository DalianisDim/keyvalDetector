/*
Copyright © 2023 Dimitris Dalianis <dimitris@dalianis.gr>
This file is part of CLI application keyvalDetector
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var cfgFile string
var version bool
var buildTimeVersion string
var defaultKubeConfigPath = homedir.HomeDir() + "/.kube/config"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "keyvalDetector",
	Version: buildTimeVersion,
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

	// rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "Config file (default is $HOME/.keyvalDetector.yaml)")
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

	var configMapsOut [][]string
	var secretsOut [][]string

	colorPrint(32, "keyvalDetector version ")
	fmt.Println(buildTimeVersion)
	colorPrint(33, "Current k8s context name: ")
	fmt.Println(getCurrentK8sContext(defaultKubeConfigPath) + "\n")

	// Build the spinner and start it
	s := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
	s.Prefix = "Scanning cluster for unused ConfigMaps and Secrets"
	s.FinalMSG = "Scanning cluster for unused ConfigMaps and Secrets...Complete!\n\n"
	s.Start()

	// uses the current context in kubeconfig
	// path-to-kubeconfig -- for example, /root/.kube/config
	config, _ := clientcmd.BuildConfigFromFlags("", defaultKubeConfigPath)

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
		// Slice to store used ConfigMap names
		namespaceUsedConfigMaps := []string{}

		// Slice to store used Secrets names
		namespaceUsedSecrets := []string{}

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

		// Iterate over the pods
		// Check each pod's Volumes and store the configmaps and secrets mounted names on a separate slice
		// Compare lists of configmaps/secrets with mounted ones
		for _, pod := range pods.Items {
			// Foreach pod's volume
			for _, volume := range pod.Spec.Volumes {
				if volume.ConfigMap != nil {
					namespaceUsedConfigMaps = append(namespaceUsedConfigMaps, volume.ConfigMap.Name)
				}
				if volume.Secret != nil {
					namespaceUsedSecrets = append(namespaceUsedSecrets, volume.Secret.SecretName)
				}
			}

			// Foreach container in pod
			for _, container := range pod.Spec.Containers {
				for _, envFrom := range container.EnvFrom {
					if envFrom.ConfigMapRef != nil {
						namespaceUsedConfigMaps = append(namespaceUsedConfigMaps, envFrom.ConfigMapRef.Name)
					}
					if envFrom.SecretRef != nil {
						namespaceUsedSecrets = append(namespaceUsedSecrets, envFrom.SecretRef.Name)
					}
				}
				for _, env := range container.Env {
					if env.ValueFrom != nil {
						if env.ValueFrom.ConfigMapKeyRef != nil {
							namespaceUsedConfigMaps = append(namespaceUsedConfigMaps, env.ValueFrom.ConfigMapKeyRef.Name)
						}
						if env.ValueFrom.SecretKeyRef != nil {
							namespaceUsedSecrets = append(namespaceUsedSecrets, env.ValueFrom.SecretKeyRef.Name)
						}
					}
				}
			}
		} // END - Foreach pod

		// Check if configmaps/secrets of this namespace
		// are in namespaceUsedConfigMaps/namespaceUsedSecrets
		for _, configmap := range configmaps.Items {
			if !contains(namespaceUsedConfigMaps, configmap.GetName()) {
				if !isSystemConfigMap(configmap.GetName()) {
					configMapsOut = append(configMapsOut, []string{configmap.GetName(), namespace.GetName()})
				}
			}
		}

		for _, secret := range secrets.Items {
			if !isSystemSecret(secret.GetName()) {
				if !contains(namespaceUsedSecrets, secret.GetName()) {
					secretsOut = append(secretsOut, []string{secret.GetName(), namespace.GetName()})
				}
			}
		}

	} // END - Foreach namespace

	s.Stop() // stop the spinner

	// Render to stdout
	colorPrint(31, "Unused ConfigMaps: \n")
	printTable(configMapsOut)

	colorPrint(31, "\nUnused Secrets: \n")
	printTable(secretsOut)

	return nil
}

// func test(clientset kubernetes.Clientset) {

// }

func isSystemConfigMap(configmap string) bool {
	systemConfigMaps := []string{"kube-root-ca.crt", "cluster-info", "kubelet-config", "kubeadm-config"}
	return contains(systemConfigMaps, configmap)
}

func isSystemSecret(secret string) bool {
	systemSecrets := []string{"foobar"}
	return contains(systemSecrets, secret)
}

func getCurrentK8sContext(kubeConfigPath string) string {
	config, _ := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: "",
		}).RawConfig()

	return config.CurrentContext
}

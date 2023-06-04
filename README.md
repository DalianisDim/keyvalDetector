

# kubevalDetector

**keyvalDetector** will scan your Kubernetes cluster for ConfigMaps and Secrets that are not used by Pods.

It uses the current Kubernetes context by default, therefore in its simplest form it can be ran without any flag:

```
kubevalDetector
```

## Usage

```
keyvalDetector [flags]

Flags:
  -c, --config string   Config file (default is $HOME/.keyvalDetector.yaml)
  -h, --help            help for keyvalDetector
  -v, --version         Print the version and exit.
```
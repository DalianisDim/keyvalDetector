

# keyvalDetector

keyvalDetector scans your Kubernetes cluster to identify ConfigMaps and Secrets that are not utilized by any running pods. It displays the names of ConfigMaps and Secrets that are not mounted by any pod and those whose contents are not being used to populate environment variables using [envFrom](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#configure-all-key-value-pairs-in-a-configmap-as-container-environment-variables) or [valueFrom](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#use-configmap-defined-environment-variables-in-pod-commands).

## Installation

### Binary
Binaries for Linux or OS X can be found in [GitHub releases page](https://github.com/DalianisDim/keyvalDetector/releases). You can use `curl` or `wget` to download it. Don't forget to `chmod +x` the file!


## Usage

`keyvalDetector` uses the current Kubernetes context by default, therefore in its simplest form it can be ran without any flag:

```
keyvalDetector
```


![Screenshot 2023-06-05 at 12 39 58](https://github.com/DalianisDim/keyvalDetector/assets/17311561/be8e355d-6e0f-4591-a9d9-7aa6d19b2404)


## Available flags
```
keyvalDetector [flags]

Flags:
  -h, --help            help for keyvalDetector
  -v, --version         Print the version and exit.
```


## Future plans

- Allow selection of context from available contexts in `~/.kube/config`
- Delete unused ConfigMaps & Secrets

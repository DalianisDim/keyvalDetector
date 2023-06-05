

# keyvalDetector

**keyvalDetector** will scan your Kubernetes cluster for ConfigMaps and Secrets that are not used by Pods.

It uses the current Kubernetes context by default, therefore in its simplest form it can be ran without any flag:

```
kubevalDetector
```


![Screenshot 2023-06-05 at 12 39 58](https://github.com/DalianisDim/keyvalDetector/assets/17311561/be8e355d-6e0f-4591-a9d9-7aa6d19b2404)

## Usage

```
keyvalDetector [flags]

Flags:
  -h, --help            help for keyvalDetector
  -v, --version         Print the version and exit.
```

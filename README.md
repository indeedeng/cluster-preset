# ClusterPreset

ClusterPreset is a [Kubernetes admission controller](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/) that injects common environment variables, volumes, and volume mounts into all pods running in the cluster.
It functions similar to `PodPresets` except that is not limited to a single namespace.
Unlike other solutions that exist, ClusterPresets do not require pods to be labeled for data injection.

## Status

This is currently an alpha project.
Ideally, we want to help contribute this functionality [upstream](https://github.com/kubernetes/kubernetes/issues/48180).
While we work with the community on driving a canonical solution forward, we've open sourced the one we're using.

![GitHub](https://img.shields.io/github/license/indeedeng/cluster-preset.svg)
[![Build Status](https://travis-ci.com/indeedeng/cluster-preset.svg?branch=master)](https://travis-ci.com/indeedeng/cluster-preset)
[![Image Layers](https://images.microbadger.com/badges/image/indeedoss/cluster-preset.svg)](https://microbadger.com/images/indeedoss/cluster-preset)
[![Image Version](https://images.microbadger.com/badges/version/indeedoss/cluster-preset.svg)](https://microbadger.com/images/indeedoss/cluster-preset)
![OSS Lifecycle](https://img.shields.io/osslifecycle/indeedeng/cluster-preset.svg)

### Supported Architectures

The cluster-preset container is built using docker-buildx to provide support across multiple architectures. The following os/architecture pairs are currently supported:

* `linux/amd64`
* `linux/arm64`
* `linux/arm/v7`

## Getting Started

The current set of Kubernetes configurations requires use of cert-manager.
[cert-manager](https://docs.cert-manager.io) is used to mint and issue TLS certificates using custom resource definitions.
Installing the system is pretty easy.

```bash
$ kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v0.8.0/cert-manager.yaml
```

Once the system is up and running, you'll be able to spin up cluster-preset pretty quickly.
First, let's apply the default configuration to get things up and running.

```bash
$ kubectl apply -f https://raw.githubusercontent.com/indeedeng/cluster-preset/master/k8s/certificate.yaml
$ kubectl apply -f https://raw.githubusercontent.com/indeedeng/cluster-preset/master/k8s/manifest.yaml
$ kubectl apply -f https://raw.githubusercontent.com/indeedeng/cluster-preset/master/k8s/webhook.yaml
```

Once everything is up and running, you can quickly test out the defaults:

```bash
$ kubectl run -it --rm --restart=Never alpine --image=alpine -- sh
# echo "${CLUSTER}"
us-west-1

# echo "${STAGING_LEVEL}" 
production
```

From there, you can update the ConfigMap used by the `cluster-preset`.

```bash
$ cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: cluster-presets-config
  namespace: kube-system
  labels:
    app: cluster-presets
data:
  presets.yaml: |
    env:
    - name: YOUR_VARIABLE
      value: your_value
    - name: YOUR_VARIABLE
      value: your_value
EOF
```

Once you've updated the config map, it can take up to a minute for the configuration to be reloaded.
This can be configured to be a smaller interval with the `--reload <Duration>` option.

## Contributing

See our [Contributing](CONTRIBUTING.md) guide for how you can help contribute to this project!

## Contact Us

Current project code owners:

* [mjpitz](https://github.com/mjpitz)
* [wmgroot](https://github.com/wmgroot)

We are working on providing a channel for communication.
In the meantime, open a ticket, watch, or star the repository.

## Code of Conduct

ClusterPreset is governed by the [Contributor Covenant v1.4.1](CODE_OF_CONDUCT.md).

## License

ClusterPreset is licensed under the [MIT License](LICENSE).

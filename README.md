# ClusterPreset

ClusterPreset provides an admission controller that enables administrators of [Kubernetes](https://kubernetes.io) Clusters to inject environment variables, volumes, and volume mounts across an entire cluster.
It functions similar to `PodPresets` except they are not limited to a single namespace.
Unlike other solutions that exist, ClusterPresets do not require pods to be labeled for data injection.

## Status

This project is an alpha status project currently.

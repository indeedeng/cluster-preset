package mutation

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/settings/v1alpha1"
)

type Patch struct {
	Op		string		`json:"op"`
	Path	string		`json:"path"`
	Value	interface{} `json:"value,omitempty"`
}

const (
	envPatchTemplate = "/spec/containers/%d/env"
)

func PatchPod(spec *v1alpha1.PodPresetSpec, pod *corev1.Pod) []*Patch {
	patches := make([]*Patch, 0)

	for i, container := range pod.Spec.Containers {
		envPatch := PatchEnvVar(container.Env, spec.Env, fmt.Sprintf(envPatchTemplate, i))
		patches = append(patches, envPatch)
	}

	return patches
}

func PatchEnvVar(source, added []corev1.EnvVar, base string) *Patch {
	idx := make(map[string]bool)
	for _, src := range source {
		idx[src.Name] = true
	}

	envVars := make([]corev1.EnvVar, 0)

	for _, add := range added {
		if _, exists := idx[add.Name]; exists {
			// already exists on source, skip
			continue
		}
		idx[add.Name] = true

		envVars = append(envVars, add)
	}

	envVars = append(envVars, source...)

	return &Patch{
		Op: "add",
		Path: base,
		Value: envVars,
	}
}

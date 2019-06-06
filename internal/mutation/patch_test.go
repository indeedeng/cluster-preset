package mutation_test

import (
	"testing"

	"github.com/mjpitz/cluster-preset/internal/mutation"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/settings/v1alpha1"
)

func Test_PatchPod_empty(t *testing.T) {
	spec := &v1alpha1.PodPresetSpec{}
	pod := &corev1.Pod{}

	patches := mutation.PatchPod(spec, pod)
	require.Len(t, patches, 0)
}

func Test_PatchEnvVar(t *testing.T) {
	existing := make([]corev1.EnvVar, 0)

	added := []corev1.EnvVar{
		{ Name: "A", Value: "EXPECTED_A" },
		{ Name: "A", Value: "RESET" },
		{ Name: "B", Value: "EXPECTED_B" },
	}

	patch := mutation.PatchEnvVar(existing, added, "")
	require.Equal(t, "add", patch.Op)
	require.Equal(t, "", patch.Path)

	value := patch.Value
	envVars, ok := value.([]corev1.EnvVar)

	require.True(t, ok)
	require.Len(t, envVars, 2)

	first := envVars[0]
	require.Equal(t, "A", first.Name)
	require.Equal(t, "EXPECTED_A", first.Value)

	second := envVars[1]
	require.Equal(t, "B", second.Name)
	require.Equal(t, "EXPECTED_B", second.Value)
}

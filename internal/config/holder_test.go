package config_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/indeedeng/cluster-preset/internal/config"
	"github.com/stretchr/testify/require"
)

const emptyContent = `
env:
`

const reloadContent = `
env:
- name: "TEST"
  value: "RELOAD"
`

func Test_ReloadingConfig(t *testing.T) {
	filename := t.Name() + ".yaml"

	_ = os.Remove(filename)
	err := ioutil.WriteFile(filename, []byte(emptyContent), 0644)
	require.Nil(t, err)

	holder, err := config.NewReloadingConfig(filename, &config.ReloadConfig{
		FailureRetryInterval: time.Second,
		ReloadInterval: time.Second,
	})
	require.Nil(t, err)

	preset := holder.Get()
	require.NotNil(t, preset)

	require.Len(t, preset.Env, 0)
	require.Len(t, preset.EnvFrom, 0)
	require.Len(t, preset.Volumes, 0)
	require.Len(t, preset.VolumeMounts, 0)

	err = ioutil.WriteFile(filename, []byte(reloadContent), 0644)
	require.Nil(t, err)

	time.Sleep(time.Second)	// wait for reload to happen

	preset = holder.Get()
	require.NotNil(t, preset)

	require.Len(t, preset.Env, 1)
	require.Len(t, preset.EnvFrom, 0)
	require.Len(t, preset.Volumes, 0)
	require.Len(t, preset.VolumeMounts, 0)

	env := preset.Env[0]
	require.Equal(t, "TEST", env.Name)
	require.Equal(t, "RELOAD", env.Value)

	_ = os.Remove(filename)
}

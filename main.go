package main

import (
	"fmt"
	"github.com/mjpitz/cluster-preset/internal/config"
	"github.com/mjpitz/cluster-preset/internal/mutation"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var longDescription = `
ClusterPreset functions similarly to PodPreset except that it is not bound by
namespaces. This enables administrators to apply a set of default environment
variables, volumes, and volume mounts to pods across their entire cluster.
`

func main() {
	preset := "preset.yaml"
	port := 8080
	failureRetryInterval := time.Second * 30
	reloadInterval := time.Minute

	cmd := &cobra.Command{
		Use: "cluster-preset",
		Long: longDescription,
		Run: func(cmd *cobra.Command, args []string) {
			holder, err := config.NewReloadingConfig(preset, &config.ReloadConfig{
				FailureRetryInterval: failureRetryInterval,
				ReloadInterval: reloadInterval,
			})

			if err != nil {
				panic(err.Error())
			}

			server := gin.New()

			mutation.RegisterMutateWebhook(server, holder)

			address := fmt.Sprintf("0.0.0.0:%d", port)
			if err := http.ListenAndServe(address, server); err != nil {
				panic(err.Error())
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&preset, "preset", preset, "(optional) path to the config file")
	flags.IntVar(&port, "port", port, "(optional) the port to bind to")
	flags.DurationVar(&failureRetryInterval, "retry-interval", failureRetryInterval, "(optional) specify the duration between reloads on failure")
	flags.DurationVar(&reloadInterval, "reload-interval", reloadInterval, "(optional) specify the duration between reloads on success")

	if err := cmd.Execute(); err != nil {
		panic(err.Error())
	}
}
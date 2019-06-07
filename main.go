package main

import (
	"fmt"
	"github.com/mjpitz/cluster-preset/internal/config"
	"github.com/mjpitz/cluster-preset/internal/mutation"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	certFile := "/etc/webhook/certs/cert.pem"
	keyFile := "/etc/webhook/certs/key.pem"

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

			go func() {
				address := fmt.Sprintf("0.0.0.0:%d", port)
				if err := http.ListenAndServeTLS(address, certFile, keyFile, server); err != nil {
					panic(err.Error())
				}
			}()

			// listening OS shutdown singal
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
			<-signalChan

			logrus.Info("received OS shutdown signal, shutting down webhook server gracefully...")
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&preset, "preset", preset, "(optional) path to the file container cluster presets")
	flags.IntVar(&port, "port", port, "(optional) the port to bind to")
	flags.DurationVar(&failureRetryInterval, "retry-interval", failureRetryInterval, "(optional) specify the duration between reloads on failure")
	flags.DurationVar(&reloadInterval, "reload-interval", reloadInterval, "(optional) specify the duration between reloads on success")
	flags.StringVar(&certFile, "cert", certFile, "(optional) file containing the x509 Certificate for HTTPS")
	flags.StringVar(&keyFile, "key", keyFile, "(optional) file containing the x509 private key to --cert")

	if err := cmd.Execute(); err != nil {
		panic(err.Error())
	}
}

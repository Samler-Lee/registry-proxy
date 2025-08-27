package cmd

import (
	"os"
	"os/signal"
	"registry-proxy/internal/server"
	"registry-proxy/pkg/console"
	"syscall"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the registry proxy server",
	Run: func(cmd *cobra.Command, args []string) {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
		go func() {
			sig := <-sigChan
			console.Log().Info("[core] 收到信号: %s, 准备关闭进程", sig)

			server.Stop()
		}()

		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

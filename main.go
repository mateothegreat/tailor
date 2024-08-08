package main

import (
	"log"
	"strings"
	"sync"

	ssh "github.com/mateothegreat/tailer/ssh"
	"github.com/mateothegreat/tailer/tail"
	"github.com/mateothegreat/tailer/util"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "tailer",
	Short: "Tailer is a tool to tail logs from multiple servers.",
	Long:  "Tailer is a tool to tail logs from multiple servers.",
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
}

func init() {
	tailCmd := &cobra.Command{
		Use:   "tail",
		Short: "Tail logs from multiple servers.",
		Long:  "Tail logs from multiple servers.",
		Run: func(cmd *cobra.Command, args []string) {
			addresses, err := cmd.Flags().GetString("addresses")
			if err != nil {
				log.Printf("error getting addresses: %v", err)
			}
			addrs := strings.Split(addresses, ",")

			files, err := cmd.Flags().GetString("files")
			if err != nil {
				log.Printf("error getting files: %v", err)
			}

			c := make([]ssh.HostConfig, len(addrs))
			for i, addr := range addrs {
				c[i] = ssh.HostConfig{
					Hostname: addr,
				}
			}

			manager, err := ssh.NewManager(c)
			if err != nil {
				log.Printf("error getting manager: %v", err)
			}

			// Tail logs on each server.
			var wg sync.WaitGroup
			for i, server := range addrs {
				wg.Add(1)
				go func(server string) {
					defer wg.Done()

					manager.Sessions[server].Connect()
					if err != nil {
						log.Printf("error getting session: %v", err)
					}

					err = tail.Run(manager.Sessions[server], tail.TailConfig{
						Color:   util.GetByInt(i),
						Command: "tail -f " + strings.ReplaceAll(files, ",", " "),
					})
					if err != nil {
						log.Printf("error tailing logs on server %s: %v", server, err)
					}

					manager.Sessions[server].Close()
				}(server)
			}

			wg.Wait()
		},
	}

	tailCmd.Flags().StringP("addresses", "a", "", "addresses to tail from")
	tailCmd.MarkFlagRequired("addresses")

	tailCmd.Flags().StringP("files", "f", "", "files to tail")
	tailCmd.MarkFlagRequired("files")

	root.AddCommand(tailCmd)
}

func main() {
	if err := root.Execute(); err != nil {
		log.Fatalf("error executing root command: %v", err)
	}
}

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/acogdev/action-target/monitor"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "action-target",
	Short: "Begin monitoring the given hosts",
	Long: `Begin monitoring the given hosts
For example:
To monitor host1, host2, and host3 on port 8080 every 5 seconds run

action-target --hosts host1,host2,host3 -p 8080 -i 5

`,
	Run: func(cmd *cobra.Command, args []string) {
		monitor.Monitor(Hosts, Port, Interval, ConfigFile)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var Hosts []string
var Port string
var Interval int
var ConfigFile string

func init() {

	rootCmd.Flags().StringSliceVar(&Hosts, "hosts", nil, "List of hosts to monitor in format host1,host2,host3")
	rootCmd.Flags().StringVarP(&Port, "port", "p", "80", "Port of hosts")
	rootCmd.Flags().IntVarP(&Interval, "interval", "i", 5, "Interval to run checks in seconds")
	rootCmd.Flags().StringVarP(&ConfigFile, "config", "c", "", "Path to config file")
}

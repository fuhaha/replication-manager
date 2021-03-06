// replication-manager - Replication Manager Monitoring and CLI for MariaDB
// Authors: Guillaume Lefranc <guillaume.lefranc@mariadb.com>
//          Stephane Varoqui  <stephane.varoqui@mariadb.com>
// This source code is licensed under the GNU General Public License, version 3.
// Redistribution/Reuse of this code is permitted under the GNU v3 license, as
// an additional term, ALL code must carry the original Author(s) credit in comment form.
// See LICENSE in this directory for the integral text.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tanji/replication-manager/config"
)

const repmgrVersion string = "0.7"

var (
	conf     config.Config
	cfgFile  string
	cfgGroup string
)

var (
	Version string
	Build   string
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Configuration file (default is config.toml)")
	rootCmd.PersistentFlags().StringVar(&cfgGroup, "config-group", "", "Configuration group (default is none)")
	rootCmd.PersistentFlags().StringVar(&conf.User, "user", "", "User for MariaDB login, specified in the [user]:[password] format")
	rootCmd.PersistentFlags().StringVar(&conf.Hosts, "hosts", "", "List of MariaDB hosts IP and port (optional), specified in the host:[port] format and separated by commas")
	rootCmd.PersistentFlags().StringVar(&conf.RplUser, "rpluser", "", "Replication user in the [user]:[password] format")
	rootCmd.Flags().StringVar(&conf.KeyPath, "keypath", "/etc/replication-manager/.replication-manager.key", "Encryption key file path")
	rootCmd.PersistentFlags().BoolVar(&conf.Verbose, "verbose", false, "Print detailed execution info")
	rootCmd.PersistentFlags().IntVar(&conf.LogLevel, "log-level", 0, "Log verbosity level")

	viper.BindPFlags(rootCmd.PersistentFlags())

	if conf.Verbose == true && conf.LogLevel == 0 {
		conf.LogLevel = 1
	}
	if conf.Verbose == false && conf.LogLevel > 0 {
		conf.Verbose = true
	}
}

func initConfig() {
	viper.SetConfigType("toml")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/replication-manager/")
		viper.AddConfigPath(".")
	}
	viper.SetEnvPrefix("MRM")
	err := viper.ReadInConfig()
	if err == nil {
		log.Println("INFO : Using config file:", viper.ConfigFileUsed())
	}
	if _, ok := err.(viper.ConfigParseError); ok {
		log.Fatalln("ERROR: Could not parse config file:", err)
	}
	if cfgGroup != "" {
		log.Println("INFO : Using configuration group", cfgGroup)
		cf2 := viper.Sub(cfgGroup)
		if cf2 == nil {
			log.Fatalln("ERROR: Could not parse configuration group", cfgGroup)
		}
		cf2.Unmarshal(&conf)
	} else {
		viper.Unmarshal(&conf)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "replication-manager",
	Short: "MariaDB Replication Manager Utility",
	Long: `replication-manager allows users to monitor interactively MariaDB 10.x GTID replication health
and trigger slave to master promotion (aka switchover), or elect a new master in case of failure (aka failover).`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the replication manager version number",
	Long:  `All software has versions. This is ours`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("MariaDB Replication Manager version", repmgrVersion)
		fmt.Println("Full Version: ", Version)
		fmt.Println("Build Time: ", Build)
	},
}

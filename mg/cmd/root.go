/*
Copyright 2018 The MetaGraph Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"os"
	"github.com/spf13/viper"
	"github.com/mitchellh/go-homedir"
)

const Banner string = "mg (metaGraf) -"

// Viper cfg file
var cfgFile string

// Flags
var Verbose	bool

var RootCmd = &cobra.Command{
	Use:   "mg",
	Short: "mg operates on collections of metaGraf's objects.",
	Long:  Banner + `is a utility that understands the metaGraf
datastructure and help you generate kubernetes primitives`,
	//Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	//},
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/mg/mg.yaml)")
	RootCmd.PersistentFlags().BoolVar(&Verbose,"verbose", false, "verbose output")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home+".config/mg/")
		viper.SetConfigName("mg.yaml")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

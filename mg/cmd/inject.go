/*
Copyright 2018-2019 The MetaGraph Authors

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
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"metagraf/pkg/modules"
	"os"
)

func init() {
	RootCmd.AddCommand(injectCmd)
	injectCmd.AddCommand(injectAnnotationsCmd)
	injectAnnotationsCmd.Flags().StringSliceVar(&CVars, "values", []string{}, "Slice of key=value pairs, seperated by ,")
}

var injectCmd = &cobra.Command{
	Use:   "inject",
	Short: "inject operations",
	Long:  Banner + ` inject `,
}

var injectAnnotationsCmd = &cobra.Command{
	Use: "annotations",
	Short: "inject annotations",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			glog.Info(StrActiveProject, viper.Get("namespace"))
			glog.Error(StrMissingMetaGraf)
			os.Exit(1)
		}

		if len(Namespace) == 0 {
			Namespace = viper.GetString("namespace")
			if len(Namespace) == 0 {
				glog.Error(StrMissingNamespace)
				os.Exit(1)
			}
		}
		FlagPassingHack()

		mg := metagraf.Parse(args[0])

		if len(modules.NameSpace) == 0 {
			modules.NameSpace = Namespace
		}

		if BaseEnvs {
			modules.BaseEnvs = BaseEnvs
		}

		modules.GenBuildConfig(&mg)
	},
}



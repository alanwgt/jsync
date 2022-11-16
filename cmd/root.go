// Copyright © 2022 Alan Weingartner <hi@alanwgt.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//  1. Redistributions of source code must retain the above copyright notice,
//     this list of conditions and the following disclaimer.
//
//  2. Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package cmd

import (
	jCfg "github.com/alanwgt/jsync/config"
	"github.com/alanwgt/jsync/internal/config"
	"github.com/alanwgt/jsync/internal/jsync"
	"github.com/alanwgt/jsync/log"
	"github.com/rs/zerolog"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const version = "v0.0.0"

var cfg *config.JetimobCfg
var jSync *jsync.JSync

var rootCmd = &cobra.Command{
	Use:     "jsync",
	Version: version,
	Short:   "Ferramenta para sincronização de dados dentro da Jetimob com um banco de dados",
	Long: `O jsync sincroniza informações de uma imobiliária, dentro da Jetimob, através da rota de webservice.
Os dados replicados são: imóveis, condomínios, corretores e banners.

Normalmente o comando é executado apenas da seguinte forma: jsync sync all

Isso cuidará do download de todos os dados e a sincronização dos mesmos no banco de dados local.

Acesse https://github.com/alanwgt/jsync para mais informações.`,
}

func Execute() {
	t := time.Now()
	err := rootCmd.Execute()
	if err != nil {
		log.Error().Err(err)
	}
	log.Log.Debug().Msgf("tempo de execução: %s", time.Now().Sub(t).Round(time.Millisecond).String())
}

func init() {
	rootCmd.SetUsageTemplate(`Uso:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Exemplos:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Comandos disponíveis:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Flags globais:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [comando] --help" para mais informações.{{end}}
`)

	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "aumenta a verbosidade dos logs")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "arquivo de configuração (padrão é $HOME/.jsync.yaml)")
	rootCmd.PersistentFlags().StringVarP(&tenantId, "tenant", "t", "", "faz a sincronização apenas para o id do tenant fornecido")

	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".jsync")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err = viper.ReadInConfig()
	if err != nil {
		var cfgPath string
		if cfgFile != "" {
			cfgPath = cfgFile
		} else {
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)

			cfgPath = path.Join(home, ".jsync.yaml")
		}

		log.Info().Msg("arquivo de configuração não encontrado, criando")
		f, err := os.Create(cfgPath)
		cobra.CheckErr(err)
		defer f.Close()
		_, err = f.WriteString(jCfg.EmbeddedConfig)
		cobra.CheckErr(err)
		log.Info().Str("path", cfgPath).Msg("o arquivo de configuração foi criado, por favor, configure conforme sua necessidade")
		os.Exit(0)
	}

	cobra.CheckErr(err)
	cfg = &config.JetimobCfg{}
	cobra.CheckErr(viper.Unmarshal(cfg))
	cfg.CmdCfg = config.CmdCfg{
		TenantId:       tenantId,
		IgnoreLastSync: ignoreLastSync,
		MaxPages:       maxPages,
	}

	var lvl zerolog.Level
	if verbose {
		lvl = zerolog.DebugLevel
	} else {
		lvl = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(lvl)

	jSync, err = jsync.New(cfg, version)
	cobra.CheckErr(err)
}

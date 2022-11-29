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
	"errors"
	"github.com/alanwgt/jsync/internal/shell"
	"github.com/alanwgt/jsync/log"
	"github.com/spf13/cobra"
	"math"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"s"},
	Short:   "Sincroniza recursos da jetimob com o banco de dados local",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !ignoreDbLock && truncate && cfg.LastSync != nil && !ignoreLastSync {
			return errors.New(`se as tabelas forem truncadas e houver uma data de sincronização anterior, dados serão perdidos!
Remova a opção "truncate", ou utilize a flag --ignore-last-sync para buscar todos os imóveis ignorando a data de sincronização`)
		}

		if cfg.TenantMapping != nil {
			if len(cfg.TenantMapping) == 0 && cfg.WebserviceKey == nil {
				return errors.New("a chave de webservice ou o mapeamento de tenants deve ser configurado")
			}

			if len(cfg.TenantMapping) > 0 && cfg.WebserviceKey != nil && *cfg.WebserviceKey != "" {
				return errors.New("a chave de webservice OU o mapeamento de tenants deve ser configurado, NÃO os dois")
			}

			if len(cfg.TenantMapping) > 0 && cfg.TenantDiscriminatorColumn == nil || *cfg.TenantDiscriminatorColumn == "" {
				return errors.New("a coluna discriminatória precisa estar configurada para ambiente multi tenancy")
			}
		} else if cfg.WebserviceKey == nil {
			return errors.New("a chave de webservice precisa ser especificada")
		}

		for _, m := range cfg.TenantMapping {
			if m.Identifier == "" || m.WebserviceKey == "" {
				return errors.New("a configuração de um dos tenants está vazia, por favor, remover a entrada ou incluir todas as chaves")
			}
		}

		if preHook != "" {
			log.Debug().Str("pre-hook", preHook).Msg("executando pre-hook")
			return shell.Exec(preHook)
		}

		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if postHook != "" {
			log.Debug().Str("post-hook", postHook).Msg("executando post-hook")
			return shell.Exec(postHook)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.PersistentFlags().IntVarP(&maxPages, "max-pages", "m", math.MaxInt, "número máximo de páginas requisitadas por recurso (útil para testes)")
	syncCmd.PersistentFlags().IntVar(&concurrentRequests, "concurrent-requests", 5, "máximo de requisições em paralelo (máximo 5)")
	syncCmd.PersistentFlags().BoolVarP(&ignoreLastSync, "ignore-last-sync", "i", false, "ignora a data da última sincronização, forçando a atualização de todos os dados")
	syncCmd.PersistentFlags().BoolVar(&truncate, "truncate", false, "trunca a(s) tabela(s) utilizada(s) durante a sincronização")
	syncCmd.PersistentFlags().StringVar(&preHook, "pre-hook", "", "comando para ser executado no shell antes de iniciar a sincronização")
	syncCmd.PersistentFlags().StringVar(&postHook, "post-hook", "", "comando para ser executado no shell após a sincronização bem sucedida")
}

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
	"database/sql"
	"github.com/doug-martin/goqu/v9"
	"github.com/spf13/cobra"
)

var ignoreDbLock bool
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Interação com o banco de dados",
}

var dbClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Limpa todas as tabelas",
	PreRun: func(cmd *cobra.Command, args []string) {
		ignoreDbLock = true
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		err := jSync.Db().ExecInTx(func(tx *sql.Tx) error {
			truncate := func(table string) error {
				s, _, err := goqu.Truncate(table).Cascade().ToSQL()
				if err != nil {
					return err
				}

				_, err = tx.Exec(s)
				return err
			}

			var err error
			if err = truncate(jSync.GetBannersTable()); err != nil {
				return err
			}

			if err = truncate(jSync.GetCondominiumsTable()); err != nil {
				return err
			}

			if err = truncate(jSync.GetPropertiesTable()); err != nil {
				return err
			}

			if err = truncate(jSync.GetBrokersTable()); err != nil {
				return err
			}

			return nil
		})

		if err == nil {
			jSync.L.Info().Msg("tabelas truncadas")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbClearCmd)
}

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

package jsync

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/alanwgt/jsync/internal/config"
	"github.com/alanwgt/jsync/internal/db"
	"github.com/alanwgt/jsync/internal/http"
	"github.com/alanwgt/jsync/internal/model"
	"github.com/alanwgt/jsync/log"
	"github.com/doug-martin/goqu/v9"
	"github.com/rs/zerolog"
	"reflect"
	"sort"
	"strings"
	"time"
)

type JSync struct {
	config        *config.JetimobCfg
	requester     *http.Requester
	db            *db.Db
	multiTenant   bool
	currentTenant *config.TenantMapping
	L             zerolog.Logger
}

func New(cfg *config.JetimobCfg, version string) (*JSync, error) {
	d := db.New(cfg.DB.ConnectionString)
	if err := d.Connect(); err != nil {
		return nil, err
	}

	return &JSync{
		config:      cfg,
		requester:   http.NewRequester(cfg.CmdCfg.MaxPages, cfg.CmdCfg.ConcurrentRequests),
		db:          d,
		multiTenant: len(cfg.TenantMapping) > 0,
		L:           log.Log.With().Str("version", version).Logger(),
	}, nil
}

func (j *JSync) SetCurrentTenant(t config.TenantMapping) {
	if j.multiTenant && t.Identifier == "" || t.WebserviceKey == "" {
		j.L.Panic().Str("identifier", t.Identifier).Str("webservice_key", t.WebserviceKey).Msg("falha de configuração em ambiente multi tenancy")
	}

	j.currentTenant = &t
	j.requester.SetWebserviceKey(t.WebserviceKey)

	if j.multiTenant {
		j.L = j.L.With().Str("tenant", t.Identifier).Logger()
	}
}

func (j JSync) Config() *config.JetimobCfg {
	return j.config
}

func (j JSync) Requester() *http.Requester {
	return j.requester
}

func (j JSync) Db() *db.Db {
	return j.db
}

func getDefaultTableName(opt *string, def string) string {
	if opt == nil {
		return def
	}

	return *opt
}

func (j JSync) GetPropertiesTable() string {
	return getDefaultTableName(j.config.Mappings.PropertiesTable, config.DefaultPropertiesTable)
}

func (j JSync) GetCondominiumsTable() string {
	return getDefaultTableName(j.config.Mappings.CondominiumsTable, config.DefaultCondominiumsTable)
}

func (j JSync) GetBannersTable() string {
	return getDefaultTableName(j.config.Mappings.BannersTable, config.DefaultBannersTable)
}

func (j JSync) GetBrokersTable() string {
	return getDefaultTableName(j.config.Mappings.BrokersTable, config.DefaultBrokersTable)
}

func (j JSync) GetTenants() []config.TenantMapping {
	if !j.multiTenant {
		return []config.TenantMapping{{
			Identifier:    "",
			WebserviceKey: *j.config.WebserviceKey,
		}}
	}

	if j.config.CmdCfg.TenantId != "" {
		j.L.Debug().Str("tenant", j.config.CmdCfg.TenantId).Msg("requsitado uso de configuração para apenas um tenant, encontrando configuração")
		for _, t := range j.config.TenantMapping {
			if t.Identifier == j.config.CmdCfg.TenantId {
				j.L.Debug().Str("tenant", j.config.CmdCfg.TenantId).Msg("tenant encontrado")
				return []config.TenantMapping{{
					Identifier:    t.Identifier,
					WebserviceKey: t.WebserviceKey,
				}}
			}
		}

		j.L.Panic().Str("tenant_id", j.config.CmdCfg.TenantId).Msg("tenant não encontrado para o id fornecido")
	}

	return j.config.TenantMapping
}

func sync[T model.Model](tx *sql.Tx, j JSync, values []T, colMap map[string]any, table string) error {
	l := j.L.With().Str("table", table).Logger()
	l.Debug().Msg("iniciando sincronização de dados")

	if len(values) == 0 {
		return nil
	}

	tagMap := extractTagMap(values[0])
	pks := make([]int, len(values))
	cols, i := make([]string, len(tagMap)), 0

	for k := range tagMap {
		cols[i] = fmt.Sprintf(`"%s"`, k)
		i++
	}

	sort.Strings(cols)

	var inserts []map[any]any
	for vi, v := range values {
		m := make(map[any]any)
		for _, col := range cols {
			col = strings.Replace(col, `"`, "", 2)
			colName, ok := tagMap[col]
			if !ok {
				return errors.New(fmt.Sprintf(`não há propriedade "%s" para a tabela "%s"`, col, table))
			}

			colSplit := strings.Split(col, ".")
			col = colSplit[len(colSplit)-1]

			remappedCol, ok := colMap[col]
			if !ok {
				return errors.New(fmt.Sprintf(`configuração de mapeamento para a coluna "%s" da tabela "%s" não encontrada`, col, table))
			}

			m[remappedCol] = reflect.ValueOf(v).FieldByName(colName.Name()).Interface()
		}

		if j.multiTenant {
			m[*j.config.TenantDiscriminatorColumn] = j.currentTenant.Identifier
		}

		pks[vi] = v.Identifier()
		inserts = append(inserts, m)
	}

	if j.config.TruncateAll {
		l.Warn().Bool("truncate", true).Msg("truncando tabela")
		exp := goqu.Delete(table)
		if j.multiTenant {
			exp = exp.Where(goqu.C(*j.config.TenantDiscriminatorColumn).Eq(j.currentTenant.Identifier))
		}
		q, _, err := exp.ToSQL()
		if err != nil {
			return err
		}

		_, err = tx.Exec(q)
		if err != nil {
			return err
		}

		l.Info().Msg("tabela truncada")
	} else {
		exp := goqu.Delete(table)
		if j.multiTenant {
			exp = exp.Where(goqu.C(*j.config.TenantDiscriminatorColumn).Eq(j.currentTenant.Identifier))
		}

		q, _, err := exp.
			Where(goqu.C("id").In(pks)).
			ToSQL()
		if err != nil {
			return err
		}

		_, err = tx.Exec(q)
		if err != nil {
			return err
		}

		l.Info().Ints("ids", pks).Msg("rows desatualizadas removidas da tabela")
	}

	q, _, err := goqu.
		Dialect("postgres").
		Insert(table).
		Rows(inserts).
		ToSQL()

	if err != nil {
		return err
	}

	_, err = tx.Exec(q)

	if err != nil {
		l.Error().Err(err).Msg("falha ao inserir dados no banco")
	} else {
		l.Info().Msg("dados sincronizados")
	}

	return err
}

func (j JSync) MarkPropertiesAsActive(tx *sql.Tx, ids []int) error {
	table := j.GetPropertiesTable()
	exp := goqu.Update(table)
	if j.multiTenant {
		exp = exp.Where(goqu.C(*j.config.TenantDiscriminatorColumn).Eq(j.currentTenant.Identifier))
	}

	q, _, err := exp.
		Set(goqu.Record{"active": false}).
		Where(goqu.C("active").Eq(true)).
		ToSQL()
	if err != nil {
		return err
	}

	j.L.Debug().Msg("marcando todos imóveis como inativos")
	if _, err = tx.Exec(q); err != nil {
		return err
	}

	exp = goqu.Update(table)
	if j.multiTenant {
		exp = exp.Where(goqu.C(*j.config.TenantDiscriminatorColumn).Eq(j.currentTenant.Identifier))
	}

	q, _, err = exp.
		Set(goqu.Record{"active": true}).
		Where(goqu.C("id").In(ids)).
		ToSQL()
	if err != nil {
		return err
	}

	j.L.Debug().Ints("ids", ids).Msg("marcando imóveis ativos")
	if _, err = tx.Exec(q); err != nil {
		return err
	}

	return nil
}

func (j JSync) SyncAll() error {
	return j.Db().ExecInTx(func(tx *sql.Tx) error {
		if err := j.SyncBanners(tx); err != nil {
			return err
		}

		if err := j.SyncBrokers(tx); err != nil {
			return err
		}

		if err := j.SyncCondominiums(tx); err != nil {
			return err
		}

		if err := j.SyncProperties(tx); err != nil {
			return err
		}

		return nil
	})
}

func syncSingle[T model.Model](tx *sql.Tx, j JSync, vs []T, colMap map[string]any, table string) error {
	if tx != nil {
		return sync(tx, j, vs, colMap, table)
	}

	return j.Db().ExecInTx(func(tx *sql.Tx) error {
		return sync(tx, j, vs, colMap, table)
	})
}

func (j JSync) SyncProperties(tx *sql.Tx) error {
	j.L.Info().Msg("iniciando sincronização de imóveis")
	var lastSync *time.Time

	if j.config.CmdCfg.IgnoreLastSync {
		j.L.Debug().Msg("ignorando última data de sincronização, requisitando todos os imóveis")
	} else if j.config.LastSync != nil {
		lastSync = j.config.LastSync
		j.L.Debug().Time("last_sync", *lastSync).Msg("utilizando útlima data de sincronização")
	}

	cs, err := j.requester.GetProperties(lastSync)
	if err != nil {
		return err
	}

	if err = syncSingle(tx, j, cs, j.config.Mappings.Properties, j.GetPropertiesTable()); err != nil {
		return err
	}

	if err := config.SaveWith("last_sync", time.Now()); err != nil {
		j.L.Error().Err(err).Msg("falha ao salvar data de sincronização no arquivo de configuração")
	}

	if err := j.SyncActiveProperties(tx); err != nil {
		return err
	}

	return nil
}

func (j JSync) SyncBrokers(tx *sql.Tx) error {
	j.L.Info().Msg("iniciando sincronização de corretores")
	bs, err := j.requester.GetBrokers()
	if err != nil {
		return err
	}

	return syncSingle(tx, j, bs, j.config.Mappings.Brokers, j.GetBrokersTable())
}

func (j JSync) SyncBanners(tx *sql.Tx) error {
	j.L.Info().Msg("iniciando sincronização de banners")
	vs, err := j.requester.GetBanners()
	if err != nil {
		return err
	}

	return syncSingle(tx, j, vs, j.config.Mappings.Banners, j.GetBannersTable())
}

func (j JSync) SyncCondominiums(tx *sql.Tx) error {
	j.L.Info().Msg("iniciando sincronização de condomínios")
	cs, err := j.requester.GetCondominiums()
	if err != nil {
		return err
	}

	e := syncSingle(tx, j, cs, j.config.Mappings.Condominiums, j.GetCondominiumsTable())
	return e
}

func (j JSync) SyncActiveProperties(tx *sql.Tx) error {
	j.L.Info().Msg("iniciando sincronização de imóveis ativos")
	ids, err := j.requester.GetActiveProperties()
	if err != nil {
		return err
	}

	if tx != nil {
		return j.MarkPropertiesAsActive(tx, ids)
	}

	return j.Db().ExecInTx(func(tx *sql.Tx) error {
		return j.MarkPropertiesAsActive(tx, ids)
	})
}

func (j *JSync) ForEachTenant(f func() error) error {
	for _, t := range j.GetTenants() {
		j.SetCurrentTenant(t)
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}

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

package config

import (
	"time"
)

const (
	DefaultPropertiesTable   = "properties"
	DefaultCondominiumsTable = "condominiums"
	DefaultBannersTable      = "banners"
	DefaultBrokersTable      = "brokers"
)

type DB struct {
	ConnectionString string `mapstructure:"connection_string"`
}

type TenantMapping struct {
	Identifier    string `mapstructure:"identifier"`
	WebserviceKey string `mapstructure:"webservice_key"`
}

type Mappings struct {
	CondominiumsTable *string        `mapstructure:"condominiums_table"`
	PropertiesTable   *string        `mapstructure:"properties_table"`
	BrokersTable      *string        `mapstructure:"brokers_table"`
	BannersTable      *string        `mapstructure:"banners_table"`
	Condominiums      map[string]any `mapstructure:"condominiums"`
	Properties        map[string]any `mapstructure:"properties"`
	Brokers           map[string]any `mapstructure:"brokers"`
	Banners           map[string]any `mapstructure:"banners"`
}

type JetimobCfg struct {
	DB                        DB              `mapstructure:"db"`
	WebserviceKey             *string         `mapstructure:"webservice_key"`
	LastSync                  *time.Time      `mapstructure:"last_sync"` // última vez que houve uma sincronização com a Jetimob
	TenantDiscriminatorColumn *string         `mapstructure:"tenant_column"`
	TenantMapping             []TenantMapping `mapstructure:"tenant_mapping"`
	TruncateAll               bool            `mapstructure:"truncate_all"`
	Mappings                  Mappings        `mapstructure:"mappings"`
	CmdCfg                    CmdCfg
}

type CmdCfg struct {
	TenantId       string
	IgnoreLastSync bool
	MaxPages       int
}

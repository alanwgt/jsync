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

package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/alanwgt/jsync/log"
	_ "github.com/lib/pq"
)

type Db struct {
	driver           string
	connectionString string
	connection       *sql.DB
}

func New(connectionString string) *Db {
	return &Db{
		driver:           "postgres",
		connectionString: connectionString,
	}
}

func (db *Db) Connect() (err error) {
	db.connection, err = sql.Open(db.driver, db.connectionString)
	if err != nil {
		return
	}

	return nil
}

func (db Db) Exec(query string, args ...any) (sql.Result, error) {
	return db.connection.Exec(query, args...)
}

func (db Db) Query(query string, args ...any) (*sql.Rows, error) {
	return db.connection.Query(query, args...)
}

func (db Db) Connection() *sql.DB {
	return db.connection
}

func (db Db) Truncate(table string) (sql.Result, error) {
	return db.Exec(fmt.Sprintf("TRUNCATE %s", table))
}

func (db Db) ExecInTx(f func(*sql.Tx) error) error {
	tx, err := db.Connection().BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	err = f(tx)

	if err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			log.Error().Err(err).Msg("falha ao realizar rollback da transação")
			return txErr
		}
		return err
	} else if err = tx.Commit(); err != nil {
		log.Error()
		return err
	}

	return nil
}

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

package json

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/doug-martin/goqu/v9/exp"
	"gopkg.in/guregu/null.v4"
	"strings"
	"time"
)

type StrSlice []string
type CommaStrSlice []string

func (ss *CommaStrSlice) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil || len(strings.TrimSpace(s)) == 0 {
		*ss = []string{}
		return nil
	}

	*ss = strings.Split(s, ",")
	return nil
}

func (ss CommaStrSlice) Value() (driver.Value, error) {
	if len(ss) == 0 {
		return exp.NewLiteralExpression("array[]::text[]"), nil
	}

	s := "array["
	addStr := func(addStr string) {
		s += fmt.Sprintf("'%s'", strings.TrimSpace(addStr))
	}

	addStr(ss[0])
	for i := 1; i < len(ss); i++ {
		s += ","
		addStr(ss[i])
	}

	s += "]"
	return exp.NewLiteralExpression(s), nil
}

func (ss StrSlice) Value() (driver.Value, error) {
	return CommaStrSlice(ss).Value()
}

type IntBool bool

func (ib *IntBool) UnmarshalJSON(data []byte) error {
	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	*ib = i == 1
	return nil
}

type JTime struct {
	time.Time
}

func (jt *JTime) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}

	jt.Time, err = time.Parse("2006-01-02 15:04:05", s)
	return err
}

func (jt JTime) Value() (driver.Value, error) {
	return json.Marshal(jt.Time)
}

type NullEmptyString null.String

func (ns *NullEmptyString) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return err
	}

	if len(strings.TrimSpace(s)) == 0 {
		*ns = (NullEmptyString)(null.NewString("", false))
		return nil
	}

	*ns = (NullEmptyString)(null.NewString(s, true))
	return nil
}

func (ns NullEmptyString) Value() (driver.Value, error) {
	return ns.NullString.Value()
}

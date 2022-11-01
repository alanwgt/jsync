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

package model

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/alanwgt/jsync/internal/encoding"
	"gopkg.in/guregu/null.v4"
)

type Model interface {
	Identifier() int
	IdentityColumn() string
}

type Media struct {
	Url   string      `json:"link"`
	Title null.String `json:"titulo"`
}

type MappedMedia struct {
	Url   string      `json:"url"`
	Title null.String `json:"title"`
}

func (m Media) MarshalJSON() ([]byte, error) {
	return json.Marshal(&MappedMedia{
		Url:   m.Url,
		Title: m.Title,
	})
}

func (m Media) Value() (driver.Value, error) {
	return json.Marshal(m)
}

type MediaArray []Media

func (ma MediaArray) Value() (driver.Value, error) {
	return encoding.JsonbArray(ma)
}

type Video struct {
	Url   string      `json:"href"`
	Title null.String `json:"title"`
}

type MappedVideo struct {
	Url   string      `json:"url"`
	Title null.String `json:"title"`
}

func (v Video) MarshalJSON() ([]byte, error) {
	return json.Marshal(&MappedVideo{
		Url:   v.Url,
		Title: v.Title,
	})
}

func (v Video) Value() (driver.Value, error) {
	return json.Marshal(v)
}

type VideoArray []Video

func (v VideoArray) Value() (driver.Value, error) {
	return encoding.JsonbArray(v)
}

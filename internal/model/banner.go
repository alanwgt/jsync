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
	"github.com/alanwgt/jsync/internal/json"
	"gopkg.in/guregu/null.v4"
)

type Banner struct {
	JetimobBannerId int          `json:"id_banner"`
	JetimobImageId  int          `json:"id_imagem"`
	URL             null.String  `json:"link"`
	Order           int          `json:"ordem"`
	Title           null.String  `json:"titulo"`
	Description     null.String  `json:"descricao"`
	HrefTarget      string       `json:"abrir_em"`
	IsVideo         json.IntBool `json:"is_video"`
	VideoUrl        string       `json:"video"`
	ImageUrl        string       `json:"imagem"`
}

func (b Banner) Identifier() int {
	return b.JetimobBannerId
}

func (b Banner) IdentityColumn() string {
	return "jetimob_banner_id"
}

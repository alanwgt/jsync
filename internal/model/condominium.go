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

type Condominium struct {
	Id                       int                `json:"id_condominio"`
	Type                     string             `json:"tipo"`
	JetimobCoverImageId      null.Int           `json:"id_imagem"`
	Name                     string             `json:"nome"`
	Featured                 bool               `json:"destaque"`
	Launch                   bool               `json:"lancamento"`
	Gated                    bool               `json:"fechado"`
	ConstructionPercentage   null.Int           `json:"acabamentos"`
	BrickPercentage          null.Int           `json:"alvenaria"`
	StructurePercentage      null.Int           `json:"estruturas"`
	FoundationPercentage     null.Int           `json:"fundacoes"`
	InstallationsPercentage  null.Int           `json:"intalacoes"`
	LandscapingPercentage    null.Int           `json:"paisagismo"`
	ProjectPercentage        null.Int           `json:"projetos"`
	GroundLevelingPercentage null.Int           `json:"terraplanagem"`
	Latitude                 float32            `json:"latitude"`
	Longitude                float32            `json:"longitue"`
	Notes                    string             `json:"observacoes"`
	IncorporationRecord      null.String        `json:"registro_incorporacao"`
	JetimobNeighborhoodId    int                `json:"id_bairro"`
	JetimobCityId            int                `json:"id_cidade"`
	JetimobStateId           int                `json:"id_estado"`
	AddressZipcode           string             `json:"endereco_cep"`
	AddressStreet            string             `json:"endereco_logradouro"`
	AddressNeighborhood      string             `json:"endereco_bairro"`
	AddressNumber            string             `json:"endereco_numero"`
	CityName                 string             `json:"endereco_cidade"`
	StateName                string             `json:"endereco_estado"`
	Status                   string             `json:"situacao"`
	DeliveryMonth            null.Int           `json:"entrega_mes"`
	DeliveryYear             null.Int           `json:"entrega_ano"`
	AdministeringCompanyName null.String        `json:"administradora"`
	BuildingCompanyName      null.String        `json:"construtora"`
	RealStateDeveloperName   null.String        `json:"incorporadora"`
	ArchitectName            null.String        `json:"projeto_arquitetonico"`
	LandscaperName           null.String        `json:"projeto_paisagismo"`
	DecoratorName            null.String        `json:"projeto_decoracao"`
	CoverImage               string             `json:"logotipo"`
	AvailableProperties      int                `json:"total_imoveis_disponiveis"`
	Infrastructures          json.CommaStrSlice `json:"infraestruturas"`
	Labels                   json.CommaStrSlice `json:"etiquetas"`
	Videos                   VideoArray         `json:"videos"`
	Images                   MediaArray         `json:"imagens"`
	Blueprints               MediaArray         `json:"plantas"`
	ArTour                   json.StrSlice      `json:"tour360"`
	CreatedAt                json.JTime         `json:"data_cadastro"`
	UpdatedAt                json.JTime         `json:"data_update"`
}

func (c Condominium) Identifier() int {
	return c.Id
}

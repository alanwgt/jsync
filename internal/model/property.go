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
	j "encoding/json"
	"github.com/alanwgt/jsync/internal/encoding"
	"github.com/alanwgt/jsync/internal/json"
	"gopkg.in/guregu/null.v4"
)

type Rural struct {
	Activities   json.CommaStrSlice `json:"atividade_rural"`
	Headquarters null.Int           `json:"rural_sedes"`
	ArableArea   null.Float         `json:"rural_area_aravel"`
}

type SeasonCalendarArray []SeasonCalendar
type SeasonCalendar struct {
	Name               null.String `json:"nome"`
	StartDate          json.JTime  `json:"inicio"`
	EndDate            json.JTime  `json:"fim"`
	DailyRate          null.Float  `json:"valor_diaria"`
	MinimumDailyRental null.Int    `json:"minimo_diarias"`
}

type Property struct {
	Id                         int                 `json:"id_imovel"`
	Type                       string              `json:"tipo"`
	CondominiumId              null.Int            `json:"id_condominio"`
	BrokerId                   int                 `json:"id_corretor"`
	StateId                    int                 `json:"id_estado"`
	CityId                     int                 `json:"id_cidade"`
	NeighborhoodId             int                 `json:"id_bairro"`
	IdentifierCode             string              `json:"codigo"`
	Contracts                  json.CommaStrSlice  `json:"contrato"`
	Subtype                    string              `json:"subtipo"`
	Notes                      string              `json:"observacoes"`
	BuildingType               null.String         `json:"tipo_construcao"`
	DeliveryYear               null.Int            `json:"entrega_ano"`
	DeliveryMonth              null.Int            `json:"entrega_mes"`
	Furnished                  int                 `json:"mobiliado"`
	Suites                     int                 `json:"suites"`
	Bathrooms                  int                 `json:"banheiros"`
	Bedrooms                   int                 `json:"dormitorios"`
	Garages                    int                 `json:"garagens"`
	Financeable                int                 `json:"financiavel"`
	HasExclusivity             bool                `json:"exclusividade"`
	TotalArea                  null.Float          `json:"area_total"`
	PrivateArea                null.Float          `json:"area_privativa"`
	UsefulArea                 null.Float          `json:"area_util"`
	MeasurementType            string              `json:"medida"`
	FloorTypes                 json.CommaStrSlice  `json:"tipo_piso"`
	TerrainWidthFront          null.Float          `json:"terreno_frente"`
	TerrainWidthBack           null.Float          `json:"terreno_fundos"`
	TerrainLengthLeft          null.Float          `json:"terreno_esquerdo"`
	TerrainLengthRight         null.Float          `json:"terreno_direita"`
	TerrainArea                null.Float          `json:"terreno_total"`
	CreatedAt                  json.JTime          `json:"data_cadastro"`
	BuildingStatus             string              `json:"status"`
	ShowCondominiumValue       bool                `json:"valor_condominio_visivel"`
	CondominiumValue           null.Float          `json:"valor_condominio"`
	ShowSaleValue              bool                `json:"valor_venda_visivel"`
	SaleValue                  null.Float          `json:"valor_venda"`
	ShowRentalValue            bool                `json:"valor_locacao_visivel"`
	RentalValue                null.Float          `json:"valor_locacao"`
	ShowSeasonalValue          bool                `json:"valor_temporada_visivel"`
	SeasonalValue              null.Float          `json:"valor_temporada"`
	IptuFrequency              null.String         `json:"periodicidade_iptu"`
	IptuValueExempt            string              `json:"valor_iptu_isento"`
	ShowIptuValue              bool                `json:"valor_iptu_visivel"`
	IptuValue                  null.Float          `json:"valor_iptu"`
	SeasonCalendar             SeasonCalendarArray `json:"calendario_temporada"`
	RuralActivities            json.CommaStrSlice  `json:"rural.atividade_rural"`
	RuralHeadquarters          null.Int            `json:"rural.rural_sedes"`
	ArableArea                 null.Float          `json:"rural.rural_area_aravel"`
	AllowedGuests              null.Int            `json:"numero_pessoas"`
	MetaTitle                  string              `json:"meta_title"`
	MetaDescription            string              `json:"meta_description"`
	FireInsuranceValue         null.Float          `json:"valor_seguro_incendio"`
	CleaningFeeValue           null.Float          `json:"valor_taxa_limpeza"`
	Position                   null.String         `json:"posicao"`
	SolarPositions             json.CommaStrSlice  `json:"posicao_solar"`
	SeaDistance                null.Int            `json:"distancia_mar"`
	AcceptExchange             bool                `json:"permuta"`
	Latitude                   null.Float          `json:"latitude"`
	Longitude                  null.Float          `json:"longitude"`
	OccupancyStatus            string              `json:"situacao"`
	Featured                   featuredStr         `json:"destaque"`
	FeatureUntil               null.Time           `json:"destaque_fim"`
	CondominiumType            null.String         `json:"condominio_tipo"`
	CondominiumName            null.String         `json:"condominio_nome"`
	GatedCondominium           null.Bool           `json:"condominio_fechado"`
	ShowFullAddress            bool                `json:"endereco_completamente_visivel"`
	ShowAddressState           bool                `json:"endereco_estado_visivel"`
	ShowAddressCity            bool                `json:"endereco_cidade_visivel"`
	ShowAddressNeighborhood    bool                `json:"endereco_bairro_visivel"`
	ShowAddressStreet          bool                `json:"endereco_logradouro_visivel"`
	ShowAddressReference       bool                `json:"endereco_referencia_visivel"`
	ShowAddressNumber          bool                `json:"endereco_numero_visivel"`
	ShowAddressFloor           bool                `json:"andar_visivel"`
	AddressState               string              `json:"endereco_estado"`
	AddressCity                string              `json:"endereco_cidade"`
	AddressNeighborhood        string              `json:"endereco_bairro"`
	AddressStreet              string              `json:"endereco_logradouro"`
	AddressZipcode             null.String         `json:"endereco_cep"`
	AddressReference           string              `json:"endereco_referencia"`
	AddressNumber              string              `json:"endereco_numero"`
	AddressFloor               null.Int            `json:"andar"`
	GeopositionVisibility      int                 `json:"geoposicionamento_visivel"`
	AdTitle                    string              `json:"titulo_anuncio"`
	AdDescription              string              `json:"descricao_anuncio"`
	Labels                     json.CommaStrSlice  `json:"tags"`
	SuretyInsurance            null.Bool           `json:"seguro_fianca"`
	PropertyInfrastructures    json.CommaStrSlice  `json:"imovel_comodidades"`
	CondominiumInfrastructures json.CommaStrSlice  `json:"condominio_comodidades"`
	UpdatedAt                  json.JTime          `json:"updated_at"`
	ValidatedAt                json.JTime          `json:"data_atualizacao"`
	Videos                     VideoArray          `json:"videos"`
	Blueprints                 MediaArray          `json:"plantas"`
	Images                     MediaArray          `json:"imagens"`
	ArTour                     json.StrSlice       `json:"tour360"`
}

func (p *Property) UnmarshalJSON(data []byte) error {
	type localPropertyType Property
	var property localPropertyType
	var m map[string]any
	if err := j.Unmarshal(data, &property); err != nil {
		return err
	}

	if err := j.Unmarshal(data, &m); err != nil {
		return err
	}

	ruralData, ok := m["rural"]
	if !ok {
		return nil
	}

	bs, err := j.Marshal(ruralData)
	if err != nil {
		return err
	}

	var rural Rural
	if err := j.Unmarshal(bs, &rural); err != nil {
		return err
	}

	property.RuralActivities = rural.Activities
	property.RuralHeadquarters = rural.Headquarters
	property.ArableArea = rural.ArableArea
	pCopy := Property(property)
	*p = pCopy

	return nil
}

func (sc SeasonCalendarArray) Value() (driver.Value, error) {
	return encoding.JsonbArray(sc)
}

func (p Property) Identifier() int {
	return p.Id
}

type featuredStr bool

func (fs *featuredStr) UnmarshalJSON(data []byte) error {
	var s string
	if err := j.Unmarshal(data, &s); err != nil {
		return err
	}

	*fs = s == "Destaque"
	return nil
}

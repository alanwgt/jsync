package jsync

import (
	"github.com/alanwgt/jsync/internal/json"
	"strings"
)

func (j JSync) remapPropertyRow(row map[any]any) map[any]any {
	contractsMapping := j.config.Mappings.Contracts
	if contractsMapping == nil {
		return row
	}

	contracts, ok := row[j.config.Mappings.Properties["contrato"]]
	if !ok {
		j.L.Warn().Msg("mapeamento de contratos especificado, mas o mapeamento para a coluna de contratos não foi encontrado")
		return row
	}

	contractSlice, ok := contracts.(json.CommaStrSlice)
	if !ok {
		contractSlice, ok = contracts.([]string)
		if !ok {
			j.L.Warn().Interface("value", contracts).Msg("falha ao converter contratos para vetor")
			return row
		}
	}

	for i, contract := range contractSlice {
		if remap, ok := contractsMapping[strings.ToLower(contract)]; ok {
			contractSlice[i] = remap
		}
	}

	return row
}

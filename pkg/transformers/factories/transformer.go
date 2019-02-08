// VulcanizeDB
// Copyright © 2018 Vulcanize

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package factories

import (
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared/constants"
)

type Transformer struct {
	Config     shared.TransformerConfig
	Converter  Converter
	Repository Repository
}

func (transformer Transformer) NewTransformer(db *postgres.DB) shared.Transformer {
	transformer.Repository.SetDB(db)
	return transformer
}

func (transformer Transformer) Execute(logs []types.Log, header core.Header, recheckHeaders constants.TransformerExecution) error {
	transformerName := transformer.Config.TransformerName
	config := transformer.Config

	if len(logs) < 1 {
		err := transformer.Repository.MarkHeaderChecked(header.Id)
		if err != nil {
			log.Printf("Error marking header as checked in %v: %v", transformerName, err)
			return err
		}
		return nil
	}

	entities, err := transformer.Converter.ToEntities(config.ContractAbi, logs)
	if err != nil {
		log.Printf("Error converting logs to entities in %v: %v", transformerName, err)
		return err
	}

	models, err := transformer.Converter.ToModels(entities)
	if err != nil {
		log.Printf("Error converting entities to models in %v: %v", transformerName, err)
		return err
	}

	err = transformer.Repository.Create(header.Id, models)
	if err != nil {
		log.Printf("Error persisting %v record: %v", transformerName, err)
		return err
	}

	return nil
}

func (transformer Transformer) GetName() string {
	return transformer.Config.TransformerName
}

func (transformer Transformer) GetConfig() shared.TransformerConfig {
	return transformer.Config
}

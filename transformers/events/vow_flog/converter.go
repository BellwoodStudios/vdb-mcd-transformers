// VulcanizeDB
// Copyright © 2019 Vulcanize

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

package vow_flog

import (
	"errors"
	"github.com/vulcanize/vulcanizedb/pkg/core"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type VowFlogConverter struct{}

func (VowFlogConverter) ToModels(logs []core.HeaderSyncLog) ([]shared.InsertionModel, error) {
	var models []shared.InsertionModel
	for _, log := range logs {
		err := verifyLog(log.Log)
		if err != nil {
			return nil, err
		}

		era := log.Log.Topics[2].Big()

		model := shared.InsertionModel{
			TableName: "vow_flog",
			OrderedColumns: []string{
				"header_id", "era", "log_id",
			},
			ColumnValues: shared.ColumnValues{
				"era":       era.String(),
				"header_id": log.HeaderID,
				"log_id":    log.ID,
			},
			ForeignKeyValues: shared.ForeignKeyValues{},
		}
		models = append(models, model)
	}
	return models, nil
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 3 {
		return errors.New("log missing topics")
	}
	if len(log.Data) < constants.DataItemLength {
		return errors.New("log missing data")
	}
	return nil
}

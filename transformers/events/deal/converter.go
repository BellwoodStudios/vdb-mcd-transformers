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

package deal

import (
	"encoding/json"
	"errors"
	"github.com/vulcanize/mcd_transformers/transformers/shared"

	"github.com/ethereum/go-ethereum/core/types"
)

type DealConverter struct{}

func (DealConverter) ToModels(ethLogs []types.Log) (result []shared.InsertionModel, err error) {
	for _, log := range ethLogs {
		err := validateLog(log)
		if err != nil {
			return nil, err
		}

		bidId := log.Topics[2].Big()
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}

		model := shared.InsertionModel{
			TableName: "deal",
			OrderedColumns: []string{
				"header_id", "bid_id", "contract_address", "log_idx", "tx_idx", "raw_log",
			},
			ColumnToValue: map[string]interface{}{
				"bid_id":           bidId.String(),
				"contract_address": log.Address.String(),
				"log_idx":          log.Index,
				"tx_idx":           log.TxIndex,
				"raw_log":          raw,
			},
			ForeignKeyToValue: map[string]string{},
		}
		result = append(result, model)
	}

	return result, nil
}

func validateLog(ethLog types.Log) error {
	if len(ethLog.Topics) < 3 {
		return errors.New("deal log does not contain expected topics")
	}
	return nil
}

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

package vat_slip

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type VatSlipConverter struct{}

func (VatSlipConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	var models []interface{}
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}
		ilk := shared.GetHexWithoutPrefix(ethLog.Topics[1].Bytes())
		usr := shared.GetHexWithoutPrefix(ethLog.Topics[2].Bytes())
		rad := ethLog.Topics[3].Big()

		raw, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}
		model := VatSlipModel{
			Ilk:              ilk,
			Usr:              usr,
			Rad:              rad.String(),
			TransactionIndex: ethLog.TxIndex,
			LogIndex:         ethLog.Index,
			Raw:              raw,
		}
		models = append(models, model)
	}
	return models, nil
}

func verifyLog(log types.Log) error {
	numTopicInValidLog := 4
	if len(log.Topics) < numTopicInValidLog {
		return errors.New("log missing topics")
	}
	return nil
}

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

package vat_flux

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type VatFluxConverter struct{}

func (VatFluxConverter) ToModels(ethLogs []types.Log) ([]shared.InsertionModel, error) {
	var models []shared.InsertionModel
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}

		ilk := ethLog.Topics[1].Hex()
		src := common.BytesToAddress(ethLog.Topics[2].Bytes()).String()
		dst := common.BytesToAddress(ethLog.Topics[3].Bytes()).String()
		wadBytes, wadErr := shared.GetLogNoteArgumentAtIndex(3, ethLog.Data)
		if wadErr != nil {
			return nil, wadErr
		}
		wad := shared.ConvertUint256HexToBigInt(hexutil.Encode(wadBytes))

		rawLogJson, jsonErr := json.Marshal(ethLog)
		if jsonErr != nil {
			return nil, jsonErr
		}

		model := shared.InsertionModel{
			TableName: "vat_flux",
			OrderedColumns: []string{
				"header_id", "ilk_id", "src", "dst", "wad", "tx_idx", "log_idx", "raw_log",
			},
			ColumnToValue: map[string]interface{}{
				"src":     src,
				"dst":     dst,
				"wad":     wad.String(),
				"tx_idx":  ethLog.TxIndex,
				"log_idx": ethLog.Index,
				"raw_log": rawLogJson,
			},
			ForeignKeyToValue: map[string]string{
				"ilk_id": ilk,
			},
		}
		models = append(models, model)
	}

	return models, nil
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 4 {
		return errors.New("log missing topics")
	}
	return nil
}

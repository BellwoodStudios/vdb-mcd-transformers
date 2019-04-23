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

package ilk

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type JugFileIlkConverter struct{}

func (JugFileIlkConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	var models []interface{}
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}

		ilk := shared.GetHexWithoutPrefix(ethLog.Topics[2].Bytes())
		what := string(bytes.Trim(ethLog.Topics[3].Bytes(), "\x00"))
		dataBytes := shared.GetDSNoteThirdArgument(ethLog.Data)
		data := shared.ConvertUint256HexToBigInt(hexutil.Encode(dataBytes))
		raw, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}

		model := JugFileIlkModel{
			Ilk:              ilk,
			What:             what,
			Data:             data.String(),
			LogIndex:         ethLog.Index,
			TransactionIndex: ethLog.TxIndex,
			Raw:              raw,
		}
		models = append(models, model)
	}
	return models, nil
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 4 {
		return errors.New("log missing topics")
	}
	if len(log.Data) < constants.DataItemLength {
		return errors.New("log missing data")
	}
	return nil
}

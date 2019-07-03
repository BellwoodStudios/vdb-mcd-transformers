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

package flip_kick

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/geth"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type FlipKickConverter struct{}

func (FlipKickConverter) ToEntities(contractAbi string, ethLogs []types.Log) ([]interface{}, error) {
	var entities []interface{}
	for _, ethLog := range ethLogs {
		entity := &FlipKickEntity{}
		address := ethLog.Address
		abi, err := geth.ParseAbi(contractAbi)
		if err != nil {
			return nil, err
		}

		contract := bind.NewBoundContract(address, abi, nil, nil, nil)

		err = contract.UnpackLog(entity, "Kick", ethLog)
		if err != nil {
			return nil, err
		}
		entity.ContractAddress = address
		entity.Raw = ethLog
		entity.TransactionIndex = ethLog.TxIndex
		entity.LogIndex = ethLog.Index
		entities = append(entities, *entity)
	}

	return entities, nil
}

func (FlipKickConverter) ToModels(entities []interface{}) ([]interface{}, error) {
	var models []interface{}
	for _, entity := range entities {
		flipKickEntity, ok := entity.(FlipKickEntity)
		if !ok {
			return nil, fmt.Errorf("entity of type %T, not %T", entity, FlipKickEntity{})
		}
		if flipKickEntity.Id == nil {
			return nil, errors.New("FlipKick log ID cannot be nil.")
		}

		id := flipKickEntity.Id.String()
		lot := shared.BigIntToString(flipKickEntity.Lot)
		bid := shared.BigIntToString(flipKickEntity.Bid)
		tab := shared.BigIntToString(flipKickEntity.Tab)
		usr := flipKickEntity.Usr.String()
		gal := flipKickEntity.Gal.String()
		contractAddress := flipKickEntity.ContractAddress.String()
		rawLog, err := json.Marshal(flipKickEntity.Raw)
		if err != nil {
			return nil, err
		}

		model := FlipKickModel{
			BidId:            id,
			Lot:              lot,
			Bid:              bid,
			Tab:              tab,
			Usr:              usr,
			Gal:              gal,
			ContractAddress:  contractAddress,
			TransactionIndex: flipKickEntity.TransactionIndex,
			LogIndex:         flipKickEntity.LogIndex,
			Raw:              rawLog,
		}
		models = append(models, model)
	}
	return models, nil
}

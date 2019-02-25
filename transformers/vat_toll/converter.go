package vat_toll

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type VatTollConverter struct{}

func (VatTollConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	var models []interface{}
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}
		ilk := shared.GetHexWithoutPrefix(ethLog.Topics[1].Bytes())
		urn := shared.GetHexWithoutPrefix(ethLog.Topics[2].Bytes())
		take := ethLog.Topics[3].Big()

		raw, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}
		model := VatTollModel{
			Ilk:              ilk,
			Urn:              urn,
			Take:             take.String(),
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

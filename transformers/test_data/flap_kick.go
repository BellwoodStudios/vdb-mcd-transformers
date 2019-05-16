package test_data

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/events/flap_kick"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var EthFlapKickLog = types.Log{
	Address: common.HexToAddress(constants.FlapperContractAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.FlapKickSignature()),
		common.HexToHash("0x00000000000000000000000000000000000000000000000000000000069f6bc7"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000003ade68b100000000000000000000000000000000000000000000000000000000075bcd150000000000000000000000007d7bee5fcfd8028cf7b00876c5b1421c800561a6000000000000000000000000000000000000000000000000000000005be469c5"),
	BlockNumber: 65,
	TxHash:      common.HexToHash("0xee7930b76b6e93974bd3f37824644ae42a89a3887a1131a7bcb3267ab4dc0169"),
	TxIndex:     66,
	BlockHash:   fakes.FakeHash,
	Index:       67,
	Removed:     false,
}

var FlapKickEntity = flap_kick.FlapKickEntity{
	Id:               big.NewInt(111111111),
	Lot:              big.NewInt(987654321),
	Bid:              big.NewInt(123456789),
	Gal:              common.HexToAddress("0x7d7bEe5fCfD8028cf7b00876C5b1421c800561A6"),
	End:              big.NewInt(1541695941),
	Raw:              EthFlapKickLog,
	TransactionIndex: EthFlapKickLog.TxIndex,
	LogIndex:         EthFlapKickLog.Index,
}

var rawFlapKickLog, _ = json.Marshal(EthFlapKickLog)
var FlapKickModel = flap_kick.FlapKickModel{
	BidId:            FlapKickEntity.Id.String(),
	Lot:              FlapKickEntity.Lot.String(),
	Bid:              FlapKickEntity.Bid.String(),
	Gal:              FlapKickEntity.Gal.String(),
	End:              time.Unix(1541695941, 0),
	Raw:              rawFlapKickLog,
	TransactionIndex: EthFlapKickLog.TxIndex,
	LogIndex:         EthFlapKickLog.Index,
}

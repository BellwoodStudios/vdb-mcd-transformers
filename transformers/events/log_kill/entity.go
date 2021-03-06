package log_kill

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LogKillEntity struct {
	Id              *big.Int
	Pair            common.Hash
	Maker           common.Address
	PayGem          common.Address
	BuyGem          common.Address
	PayAmt          *big.Int
	BuyAmt          *big.Int
	Timestamp       uint64
	HeaderID        int64
	LogID           int64
	ContractAddress common.Address
}

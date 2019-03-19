package test_data

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_grab"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var EthVatGrabLog = types.Log{
	Address: common.HexToAddress(KovanVatContractAddress),
	Topics: []common.Hash{
		common.HexToHash(KovanVatGrabSignature),
		common.HexToHash("0x5245500000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000006a3ae20c315e845b2e398e68effe39139ec6060c"),
		common.HexToHash("0x0000000000000000000000002f34f22a00ee4b7a8f8bbc4eaee1658774c624e0"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000c43690ae4c52455000000000000000000000000000000000000000000000000000000000000000000000000000000000006a3ae20c315e845b2e398e68effe39139ec6060c0000000000000000000000002f34f22a00ee4b7a8f8bbc4eaee1658774c624e00000000000000000000000003728e9777b2a0a611ee0f89e00e01044ce4736d1ffffffffffffffffffffffffffffffffffffffffffffffffe43e9298b1380000ffffffffffffffffffffffffffffffffffffffffffffffff24fea01c7f8f8000"),
	BlockNumber: 23,
	TxHash:      common.HexToHash("0x7cb84c750ce4985f7811abf641d52ffcb35306d943081475226484cf1470c6fa"),
	TxIndex:     4,
	BlockHash:   fakes.FakeHash,
	Index:       5,
	Removed:     false,
}

var rawVatGrabLog, _ = json.Marshal(EthVatGrabLog)
var VatGrabModel = vat_grab.VatGrabModel{
	Ilk:              "5245500000000000000000000000000000000000000000000000000000000000",
	Urn:              "0000000000000000000000006a3ae20c315e845b2e398e68effe39139ec6060c",
	V:                "0000000000000000000000002f34f22a00ee4b7a8f8bbc4eaee1658774c624e0",
	W:                "0000000000000000000000003728e9777b2a0a611ee0f89e00e01044ce4736d1",
	Dink:             "115792089237316195423570985008687907853269984665640564039455584007913129639936",
	Dart:             "115792089237316195423570985008687907853269984665640564039441803007913129639936",
	LogIndex:         EthVatGrabLog.Index,
	TransactionIndex: EthVatGrabLog.TxIndex,
	Raw:              rawVatGrabLog,
}

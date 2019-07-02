package test_data

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/mcd_transformers/transformers/events/spot_file/pip"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

var EthSpotFileMatLog = types.Log{
	Address: common.HexToAddress(constants.SpotContractAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.SpotFileMatSignature()),
		common.HexToHash("0x00000000000000000000000071ce79fcfec71760d51f6b3589c0d9ec0ccb64a8"),
		common.HexToHash("0x4554482d41000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x6d61740000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e01a0b287e4554482d410000000000000000000000000000000000000000000000000000006d61740000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004d8c55aefb8c05b5c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 11257385,
	TxHash:      common.HexToHash("0xb4d19aaf5be5077db65aeeb16906a5498cfa94836952191258cc45966e1d7198"),
	TxIndex:     2,
	BlockHash:   common.HexToHash("0x968cd16acb356de42e9f3ab17583988b49173c0339af5afa3f79cecdbc111d69"),
	Index:       3,
	Removed:     false,
}

var rawSpotFileMatLog, _ = json.Marshal(EthSpotFileMatLog)
var SpotFileMatModel = shared.InsertionModel{
	TableName: "spot_file_mat",
	OrderedColumns: []string{
		"header_id", "ilk_id", "what", "data", "log_idx", "tx_idx", "raw_log",
	},
	ColumnToValue: map[string]interface{}{
		"what":    "mat",
		"data":    "1500000000000000000000000000",
		"log_idx": EthSpotFileMatLog.Index,
		"tx_idx":  EthSpotFileMatLog.TxIndex,
		"raw_log": rawSpotFileMatLog,
	},
	ForeignKeyToValue: map[string]string{
		"ilk_id": "0x4554482d41000000000000000000000000000000000000000000000000000000",
	},
}

var EthSpotFilePipLog = types.Log{
	Address: common.HexToAddress(constants.SpotContractAddress()),
	Topics: []common.Hash{
		common.HexToHash(constants.SpotFilePipSignature()),
		common.HexToHash("0x0000000000000000000000004ba936a9338ae211300ea47899fbd111fd5dca31"),
		common.HexToHash("0x4554482d41000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000008c73ec0fbcdec6b8c060bc224d94740fd41f3774"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000e0d4e8be834554482d410000000000000000000000000000000000000000000000000000000000000000000000000000008c73ec0fbcdec6b8c060bc224d94740fd41f3774000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	BlockNumber: 11257191,
	TxHash:      common.HexToHash("0xaae9e8bce346f86a01c5a3af137bc1f9bc7c0c767804a2b9b6356849aee0d7dd"),
	TxIndex:     1,
	BlockHash:   common.HexToHash("0xfa28e186578238fdd6b971add2ebe62a26dddf5ff971d50ee476c86b45362da1"),
	Index:       2,
	Removed:     false,
}

var rawSpotFilePipLog, _ = json.Marshal(EthSpotFilePipLog)
var SpotFilePipModel = pip.SpotFilePipModel{
	Ilk:              "0x4554482d41000000000000000000000000000000000000000000000000000000",
	Pip:              "0x8C73Ec0fBCdEC6b8C060BC224D94740FD41f3774",
	LogIndex:         EthSpotFilePipLog.Index,
	TransactionIndex: EthSpotFilePipLog.TxIndex,
	Raw:              rawSpotFilePipLog,
}

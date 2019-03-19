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

package test_data

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/chop_lump"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/flip"
	"github.com/vulcanize/mcd_transformers/transformers/events/cat_file/pit_vow"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var EthCatFileChopLog = types.Log{
	Address: common.HexToAddress(KovanCatContractAddress),
	Topics: []common.Hash{
		common.HexToHash(KovanCatFileChopLumpSignature),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x63686f7000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000641a0b287e66616b6520696c6b00000000000000000000000000000000000000000000000063686f700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000018EE90FF6C373E0EE4E3F0AD2"),
	BlockNumber: 110,
	TxHash:      common.HexToHash("0xe32dfe6afd7ea28475569756fc30f0eea6ad4cfd32f67436ff1d1c805e4382df"),
	TxIndex:     13,
	BlockHash:   fakes.FakeHash,
	Index:       1,
	Removed:     false,
}
var rawCatFileChopLog, _ = json.Marshal(EthCatFileChopLog)
var CatFileChopModel = chop_lump.CatFileChopLumpModel{
	Ilk:              "66616b6520696c6b000000000000000000000000000000000000000000000000",
	What:             "chop",
	Data:             "123.456789012345680589533003513",
	TransactionIndex: EthCatFileChopLog.TxIndex,
	LogIndex:         EthCatFileChopLog.Index,
	Raw:              rawCatFileChopLog,
}

var EthCatFileLumpLog = types.Log{
	Address: common.HexToAddress(KovanCatContractAddress),
	Topics: []common.Hash{
		common.HexToHash(KovanCatFileChopLumpSignature),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x6c756d7000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000641a0b287e66616b6520696c6b00000000000000000000000000000000000000000000000063686f700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000029D42B64E76714244CB"),
	BlockNumber: 110,
	TxHash:      common.HexToHash("0xe32dfe6afd7ea28475569756fc30f0eea6ad4cfd32f67436ff1d1c805e4382df"),
	TxIndex:     15,
	BlockHash:   fakes.FakeHash,
	Index:       3,
	Removed:     false,
}
var rawCatFileLumpLog, _ = json.Marshal(EthCatFileLumpLog)
var CatFileLumpModel = chop_lump.CatFileChopLumpModel{
	Ilk:              "66616b6520696c6b000000000000000000000000000000000000000000000000",
	What:             "lump",
	Data:             "12345.678901234567092615",
	TransactionIndex: EthCatFileLumpLog.TxIndex,
	LogIndex:         EthCatFileLumpLog.Index,
	Raw:              rawCatFileLumpLog,
}

var EthCatFileFlipLog = types.Log{
	Address: common.HexToAddress(KovanCatContractAddress),
	Topics: []common.Hash{
		common.HexToHash(KovanCatFileFlipSignature),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x66616b6520696c6b000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x666c697000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000064ebecb39d66616b6520696c6b000000000000000000000000000000000000000000000000666c69700000000000000000000000000000000000000000000000000000000000000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977"),
	BlockNumber: 88,
	TxHash:      common.HexToHash("0xc71ef3e9999595913d31e89446cab35319bd4289520e55611a1b42fc2a8463b6"),
	TxIndex:     12,
	BlockHash:   fakes.FakeHash,
	Index:       1,
	Removed:     false,
}

var rawCatFileFlipLog, _ = json.Marshal(EthCatFileFlipLog)
var CatFileFlipModel = flip.CatFileFlipModel{
	Ilk:              "66616b6520696c6b000000000000000000000000000000000000000000000000",
	What:             "flip",
	Flip:             "0x07Fa9eF6609cA7921112231F8f195138ebbA2977",
	TransactionIndex: EthCatFileFlipLog.TxIndex,
	LogIndex:         EthCatFileFlipLog.Index,
	Raw:              rawCatFileFlipLog,
}

var EthCatFilePitVowLog = types.Log{
	Address: common.HexToAddress(KovanCatContractAddress),
	Topics: []common.Hash{
		common.HexToHash(KovanCatFilePitVowSignature),
		common.HexToHash("0x00000000000000000000000064d922894153be9eef7b7218dc565d1d0ce2a092"),
		common.HexToHash("0x7069740000000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000008e84a1e068d77059cbe263c43ad0cdc130863313"),
	},
	Data:        hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000044d4e8be8370697400000000000000000000000000000000000000000000000000000000000000000000000000000000008e84a1e068d77059cbe263c43ad0cdc130863313"),
	BlockNumber: 87,
	TxHash:      common.HexToHash("0x6515c7dfe53f0ad83ce1173fa99032c24a07cfd8b5d5a1c1f80486c99dd52800"),
	TxIndex:     11,
	BlockHash:   fakes.FakeHash,
	Index:       2,
	Removed:     false,
}

var rawCatFilePitVowLog, _ = json.Marshal(EthCatFilePitVowLog)
var CatFilePitVowModel = pit_vow.CatFilePitVowModel{
	What:             "pit",
	Data:             "0x8E84a1e068d77059Cbe263C43AD0cDc130863313",
	TransactionIndex: EthCatFilePitVowLog.TxIndex,
	LogIndex:         EthCatFilePitVowLog.Index,
	Raw:              rawCatFilePitVowLog,
}

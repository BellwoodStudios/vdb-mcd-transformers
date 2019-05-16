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
	"github.com/vulcanize/mcd_transformers/transformers/events/vow_fess"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

var EthVowFessLog = types.Log{
	Address: common.HexToAddress(constants.VowContractAddress()),
	Topics: []common.Hash{
		common.HexToHash("0x697efb7800000000000000000000000000000000000000000000000000000000"),
		common.HexToHash("0x0000000000000000000000002f34f22a00ee4b7a8f8bbc4eaee1658774c624e0"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000539"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
	},
	Data:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000002544faa778090e00000"),
	BlockNumber: 9377319,
	TxHash:      common.HexToHash("0x71f86b6154333d88811d151a9afecd00b39d6a326ef308ac97f66ca61264d7cb"),
	TxIndex:     4,
	BlockHash:   fakes.FakeHash,
	Index:       3,
	Removed:     false,
}

var rawVowFessLog, _ = json.Marshal(EthVowFessLog)
var VowFessModel = vow_fess.VowFessModel{
	Tab:              "1337",
	LogIndex:         EthVowFessLog.Index,
	TransactionIndex: EthVowFessLog.TxIndex,
	Raw:              rawVowFessLog,
}

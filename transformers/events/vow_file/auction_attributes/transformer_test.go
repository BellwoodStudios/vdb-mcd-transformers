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

package auction_attributes_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/vow_file/auction_attributes"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Vow file auction attributes transformer", func() {
	var (
		transformer = auction_attributes.Transformer{}
		db          = test_config.NewTestDB(test_config.NewTestNode())
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("returns err if log missing topics", func() {
		badLog := core.EventLog{
			Log: types.Log{
				Topics: []common.Hash{{}},
				Data:   []byte{1, 1, 1, 1, 1},
			}}

		_, err := transformer.ToModels(constants.VowABI(), []core.EventLog{badLog}, nil)
		Expect(err).To(HaveOccurred())
	})

	It("converts a log to a model", func() {
		models, toModelsErr := transformer.ToModels(constants.VowABI(), []core.EventLog{test_data.VowFileAuctionAttributesEventLog}, db)
		Expect(toModelsErr).NotTo(HaveOccurred())

		var msgSenderAddressID int64
		msgSenderAddressErr := db.Get(&msgSenderAddressID, `SELECT id from addresses WHERE address = $1`,
			common.HexToAddress(test_data.VowFileAuctionAttributesEventLog.Log.Topics[1].Hex()).Hex())
		Expect(msgSenderAddressErr).NotTo(HaveOccurred())

		expectedModel := test_data.VowFileAuctionAttributesModel()
		expectedModel.ColumnValues[constants.MsgSenderColumn] = msgSenderAddressID

		Expect(models).To(ConsistOf(expectedModel))
	})
})

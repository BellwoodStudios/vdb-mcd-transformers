// Copyright 2018 Vulcanize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package price_feeds

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/price_feeds"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
)

type MockPriceFeedConverter struct {
	converterErr   error
	PassedLogs     []types.Log
	PassedHeaderID int64
}

func (converter *MockPriceFeedConverter) ToModels(logs []types.Log, headerID int64) ([]price_feeds.PriceFeedModel, error) {
	converter.PassedLogs = logs
	converter.PassedHeaderID = headerID
	return []price_feeds.PriceFeedModel{test_data.PriceFeedModel}, converter.converterErr
}

func (converter *MockPriceFeedConverter) SetConverterErr(e error) {
	converter.converterErr = e
}

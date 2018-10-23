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

package mocks

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/vulcanize/vulcanizedb/pkg/core"
)

type MockLogFetcher struct {
	FetchedContractAddresses [][]string
	FetchedTopics            [][]common.Hash
	FetchedBlocks            []int64
	fetcherError             error
	FetchedLogs              []types.Log
	SetBcCalled              bool
}

func (mlf *MockLogFetcher) FetchLogs(contractAddresses []string, topics [][]common.Hash, blockNumber int64) ([]types.Log, error) {
	mlf.FetchedContractAddresses = append(mlf.FetchedContractAddresses, contractAddresses)
	mlf.FetchedTopics = topics
	mlf.FetchedBlocks = append(mlf.FetchedBlocks, blockNumber)

	return mlf.FetchedLogs, mlf.fetcherError
}

func (mlf *MockLogFetcher) SetBC(bc core.BlockChain) {
	mlf.SetBcCalled = true
}

func (mlf *MockLogFetcher) SetFetcherError(err error) {
	mlf.fetcherError = err
}

func (mlf *MockLogFetcher) SetFetchedLogs(logs []types.Log) {
	mlf.FetchedLogs = logs
}

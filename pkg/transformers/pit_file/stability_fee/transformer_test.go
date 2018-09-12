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

package stability_fee_test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/pit_file/stability_fee"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/shared"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks"
	stability_fee_mocks "github.com/vulcanize/vulcanizedb/pkg/transformers/test_data/mocks/pit_file/stability_fee"
)

var _ = Describe("", func() {
	It("gets missing headers for block numbers specified in config", func() {
		repository := &stability_fee_mocks.MockPitFileStabilityFeeRepository{}
		transformer := stability_fee.PitFileStabilityFeeTransformer{
			Config:     pit_file.PitFileConfig,
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &stability_fee_mocks.MockPitFileStabilityFeeConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedStartingBlockNumber).To(Equal(pit_file.PitFileConfig.StartingBlockNumber))
		Expect(repository.PassedEndingBlockNumber).To(Equal(pit_file.PitFileConfig.EndingBlockNumber))
	})

	It("returns error if repository returns error for missing headers", func() {
		repository := &stability_fee_mocks.MockPitFileStabilityFeeRepository{}
		repository.SetMissingHeadersErr(fakes.FakeError)
		transformer := stability_fee.PitFileStabilityFeeTransformer{
			Fetcher:    &mocks.MockLogFetcher{},
			Converter:  &stability_fee_mocks.MockPitFileStabilityFeeConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("fetches logs for missing headers", func() {
		fetcher := &mocks.MockLogFetcher{}
		repository := &stability_fee_mocks.MockPitFileStabilityFeeRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}, {BlockNumber: 2}})
		transformer := stability_fee.PitFileStabilityFeeTransformer{
			Fetcher:    fetcher,
			Converter:  &stability_fee_mocks.MockPitFileStabilityFeeConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(fetcher.FetchedBlocks).To(Equal([]int64{1, 2}))
		Expect(fetcher.FetchedContractAddress).To(Equal(pit_file.PitFileConfig.ContractAddress))
		Expect(fetcher.FetchedTopics).To(Equal([][]common.Hash{{common.HexToHash(shared.PitFileStabilityFeeSignature)}}))
	})

	It("returns error if fetcher returns error", func() {
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetcherError(fakes.FakeError)
		repository := &stability_fee_mocks.MockPitFileStabilityFeeRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := stability_fee.PitFileStabilityFeeTransformer{
			Fetcher:    fetcher,
			Converter:  &stability_fee_mocks.MockPitFileStabilityFeeConverter{},
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("converts matching logs", func() {
		converter := &stability_fee_mocks.MockPitFileStabilityFeeConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileStabilityFeeLog})
		repository := &stability_fee_mocks.MockPitFileStabilityFeeRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := stability_fee.PitFileStabilityFeeTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(converter.PassedLog).To(Equal(test_data.EthPitFileStabilityFeeLog))
	})

	It("returns error if converter returns error", func() {
		converter := &stability_fee_mocks.MockPitFileStabilityFeeConverter{}
		converter.SetConverterError(fakes.FakeError)
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileStabilityFeeLog})
		repository := &stability_fee_mocks.MockPitFileStabilityFeeRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1}})
		transformer := stability_fee.PitFileStabilityFeeTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})

	It("persists pit file model", func() {
		converter := &stability_fee_mocks.MockPitFileStabilityFeeConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileStabilityFeeLog})
		repository := &stability_fee_mocks.MockPitFileStabilityFeeRepository{}
		fakeHeader := core.Header{BlockNumber: 1, Id: 2}
		repository.SetMissingHeaders([]core.Header{fakeHeader})
		transformer := stability_fee.PitFileStabilityFeeTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).NotTo(HaveOccurred())
		Expect(repository.PassedHeaderID).To(Equal(fakeHeader.Id))
		Expect(repository.PassedModel).To(Equal(test_data.PitFileStabilityFeeModel))
	})

	It("returns error if repository returns error for create", func() {
		converter := &stability_fee_mocks.MockPitFileStabilityFeeConverter{}
		fetcher := &mocks.MockLogFetcher{}
		fetcher.SetFetchedLogs([]types.Log{test_data.EthPitFileStabilityFeeLog})
		repository := &stability_fee_mocks.MockPitFileStabilityFeeRepository{}
		repository.SetMissingHeaders([]core.Header{{BlockNumber: 1, Id: 2}})
		repository.SetCreateError(fakes.FakeError)
		transformer := stability_fee.PitFileStabilityFeeTransformer{
			Fetcher:    fetcher,
			Converter:  converter,
			Repository: repository,
		}

		err := transformer.Execute()

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(fakes.FakeError))
	})
})

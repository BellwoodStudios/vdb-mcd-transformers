package test_helpers

import (
	. "github.com/onsi/gomega"
)

type BlockMetadata struct {
	BlockNumber int    `db:"block_number"`
	BlockHash   string `db:"block_hash"`
}

type VariableRes struct {
	BlockMetadata
	Value string
}

type AuctionVariableRes struct {
	VariableRes
	ContractAddress string `db:"contract_address"`
}

type MappingRes struct {
	BlockMetadata
	Key   string
	Value string
}

type AuctionMappingRes struct {
	MappingRes
	ContractAddress string `db:"contract_address"`
}

type DoubleMappingRes struct {
	BlockMetadata
	KeyOne string `db:"key_one"`
	KeyTwo string `db:"key_two"`
	Value  string
}

func AssertVariable(res VariableRes, blockNumber int, blockHash, value string) {
	Expect(res.BlockNumber).To(Equal(blockNumber))
	Expect(res.BlockHash).To(Equal(blockHash))
	Expect(res.Value).To(Equal(value))
}

func AssertAuctionVariable(res AuctionVariableRes, blockNumber int, blockHash, contractAddress, value string) {
	Expect(res.BlockNumber).To(Equal(blockNumber))
	Expect(res.BlockHash).To(Equal(blockHash))
	Expect(res.ContractAddress).To(Equal(contractAddress))
	Expect(res.Value).To(Equal(value))
}

func AssertMapping(res MappingRes, blockNumber int, blockHash, key, value string) {
	Expect(res.BlockNumber).To(Equal(blockNumber))
	Expect(res.BlockHash).To(Equal(blockHash))
	Expect(res.Key).To(Equal(key))
	Expect(res.Value).To(Equal(value))
}

func AssertDoubleMapping(res DoubleMappingRes, blockNumber int, blockHash, keyOne, keyTwo, value string) {
	Expect(res.BlockNumber).To(Equal(blockNumber))
	Expect(res.BlockHash).To(Equal(blockHash))
	Expect(res.KeyOne).To(Equal(keyOne))
	Expect(res.KeyTwo).To(Equal(keyTwo))
	Expect(res.Value).To(Equal(value))
}

func AssertAuctionMapping(res AuctionMappingRes, blockNumber int, blockHash, contractAddress, key, value string) {
	Expect(res.BlockNumber).To(Equal(blockNumber))
	Expect(res.BlockHash).To(Equal(blockHash))
	Expect(res.ContractAddress).To(Equal(contractAddress))
	Expect(res.Key).To(Equal(key))
	Expect(res.Value).To(Equal(value))
}

package test_data

import "github.com/vulcanize/mcd_transformers/transformers/shared/constants"

// This file contains "shortcuts" to some configuration values useful for testing

func CatAddress() string  { return constants.GetContractAddress("MCD_CAT") }
func FlapAddress() string { return constants.GetContractAddress("MCD_FLAP") }
func FlipAddresses() []string {
	return constants.GetContractAddresses([]string{
		"ETH_FLIP_A", "ETH_FLIP_B", "ETH_FLIP_C",
		"MCD_FLIP_REP_A", "MCD_FLIP_ZRX_A", "MCD_FLIP_OMG_A", "MCD_FLIP_BAT_A", "MCD_FLIP_DGD_A", "MCD_FLIP_GNT_A",
	})
}
func EthFlipAddress() string { return constants.GetContractAddress("ETH_FLIP_A") }
func FlopAddress() string    { return constants.GetContractAddress("MCD_FLOP") }
func JugAddress() string     { return constants.GetContractAddress("MCD_JUG") }
func SpotAddress() string    { return constants.GetContractAddress("MCD_SPOT") }
func VatAddress() string     { return constants.GetContractAddress("MCD_VAT") }
func VowAddress() string     { return constants.GetContractAddress("MCD_VOW") }

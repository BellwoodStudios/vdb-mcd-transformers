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

package pit

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/maker"
	"github.com/vulcanize/vulcanizedb/pkg/transformers/storage_diffs/shared"
)

const (
	IlkLine = "line"
	IlkSpot = "spot"
	PitDrip = "drip"
	PitLine = "Line"
	PitLive = "live"
	PitVat  = "vat"
)

var (
	// storage key and value metadata for "drip" on the Pit contract
	DripKey      = common.HexToHash(storage_diffs.IndexFive)
	DripMetadata = shared.StorageValueMetadata{
		Name: PitDrip,
		Keys: nil,
		Type: shared.Address,
	}

	IlkSpotIndex = storage_diffs.IndexOne

	// storage key and value metadata for "Spot" on the Pit contract
	LineKey      = common.HexToHash(storage_diffs.IndexThree)
	LineMetadata = shared.StorageValueMetadata{
		Name: PitLine,
		Keys: nil,
		Type: shared.Uint256,
	}

	// storage key and value metadata for "live" on the Pit contract
	LiveKey      = common.HexToHash(storage_diffs.IndexTwo)
	LiveMetadata = shared.StorageValueMetadata{
		Name: PitLive,
		Keys: nil,
		Type: shared.Uint256,
	}

	// storage key and value metadata for "vat" on the Pit contract
	VatKey      = common.HexToHash(storage_diffs.IndexFour)
	VatMetadata = shared.StorageValueMetadata{
		Name: PitVat,
		Keys: nil,
		Type: shared.Address,
	}
)

type PitMappings struct {
	StorageRepository maker.IMakerStorageRepository
	mappings          map[common.Hash]shared.StorageValueMetadata
}

func (mappings *PitMappings) SetDB(db *postgres.DB) {
	mappings.StorageRepository.SetDB(db)
}

func (mappings *PitMappings) Lookup(key common.Hash) (shared.StorageValueMetadata, error) {
	metadata, ok := mappings.mappings[key]
	if !ok {
		err := mappings.loadMappings()
		if err != nil {
			return metadata, err
		}
		metadata, ok = mappings.mappings[key]
		if !ok {
			return metadata, shared.ErrStorageKeyNotFound{Key: key.Hex()}
		}
	}
	return metadata, nil
}

func (mappings *PitMappings) loadMappings() error {
	mappings.mappings = getStaticMappings()
	ilks, err := mappings.StorageRepository.GetIlks()
	if err != nil {
		return err
	}
	for _, ilk := range ilks {
		mappings.mappings[getSpotKey(ilk)] = getSpotMetadata(ilk)
		mappings.mappings[getLineKey(ilk)] = getLineMetadata(ilk)
	}
	return nil
}

func getStaticMappings() map[common.Hash]shared.StorageValueMetadata {
	mappings := make(map[common.Hash]shared.StorageValueMetadata)
	mappings[DripKey] = DripMetadata
	mappings[LineKey] = LineMetadata
	mappings[LiveKey] = LiveMetadata
	mappings[VatKey] = VatMetadata
	return mappings
}

func getSpotKey(ilk string) common.Hash {
	return storage_diffs.GetMapping(IlkSpotIndex, ilk)
}

func getSpotMetadata(ilk string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Ilk: ilk}
	return shared.GetStorageValueMetadata(IlkSpot, keys, shared.Uint256)
}

func getLineKey(ilk string) common.Hash {
	return storage_diffs.GetIncrementedKey(getSpotKey(ilk), 1)
}

func getLineMetadata(ilk string) shared.StorageValueMetadata {
	keys := map[shared.Key]string{shared.Ilk: ilk}
	return shared.GetStorageValueMetadata(IlkLine, keys, shared.Uint256)
}

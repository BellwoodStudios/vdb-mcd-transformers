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

package spot

import (
	"fmt"

	"github.com/makerdao/vdb-mcd-transformers/transformers/shared"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/utils"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

const (
	InsertSpotIlkPipQuery = `INSERT INTO maker.spot_ilk_pip (header_id, ilk_id, pip) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	InsertSpotIlkMatQuery = `INSERT INTO maker.spot_ilk_mat (header_id, ilk_id, mat) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	insertSpotVatQuery    = `INSERT INTO maker.spot_vat (header_id, vat) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	insertSpotParQuery    = `INSERT INTO maker.spot_par (header_id, par) VALUES ($1, $2) ON CONFLICT DO NOTHING`
)

type SpotStorageRepository struct {
	db *postgres.DB
}

func (repository *SpotStorageRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

func (repository SpotStorageRepository) Create(headerID int64, metadata utils.StorageValueMetadata, value interface{}) error {
	switch metadata.Name {
	case IlkPip:
		return repository.insertIlkPip(headerID, metadata, value.(string))
	case IlkMat:
		return repository.insertIlkMat(headerID, metadata, value.(string))
	case Vat:
		return repository.insertSpotVat(headerID, value.(string))
	case Par:
		return repository.insertSpotPar(headerID, value.(string))

	default:
		panic(fmt.Sprintf("unrecognized spot contract storage name: %s", metadata.Name))
	}
}

func (repository SpotStorageRepository) insertIlkPip(headerID int64, metadata utils.StorageValueMetadata, pip string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}

	return repository.insertFieldWithIlk(headerID, ilk, IlkPip, InsertSpotIlkPipQuery, pip)
}

func (repository SpotStorageRepository) insertIlkMat(headerID int64, metadata utils.StorageValueMetadata, mat string) error {
	ilk, err := getIlk(metadata.Keys)
	if err != nil {
		return err
	}
	return repository.insertFieldWithIlk(headerID, ilk, IlkMat, InsertSpotIlkMatQuery, mat)
}

func (repository SpotStorageRepository) insertSpotVat(headerID int64, vat string) error {
	_, err := repository.db.Exec(insertSpotVatQuery, headerID, vat)
	return err
}

func (repository SpotStorageRepository) insertSpotPar(headerID int64, par string) error {
	_, err := repository.db.Exec(insertSpotParQuery, headerID, par)
	return err
}

func (repository *SpotStorageRepository) insertFieldWithIlk(headerID int64, ilk, variableName, query, value string) error {
	tx, txErr := repository.db.Beginx()
	if txErr != nil {
		return txErr
	}
	ilkID, ilkErr := shared.GetOrCreateIlkInTransaction(ilk, tx)
	if ilkErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError("ilk", ilkErr.Error())
		}
		return ilkErr
	}
	_, writeErr := tx.Exec(query, headerID, ilkID, value)

	if writeErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return shared.FormatRollbackError(variableName, writeErr.Error())
		}
		return writeErr
	}
	return tx.Commit()
}

func getIlk(keys map[utils.Key]string) (string, error) {
	ilk, ok := keys[constants.Ilk]
	if !ok {
		return "", utils.ErrMetadataMalformed{MissingData: constants.Ilk}
	}
	return ilk, nil
}

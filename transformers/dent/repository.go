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

package dent

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	repo "github.com/vulcanize/vulcanizedb/libraries/shared/repository"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
)

type DentRepository struct {
	db *postgres.DB
}

func (repository DentRepository) Create(headerID int64, models []interface{}) error {
	tx, dBaseErr := repository.db.Begin()
	if dBaseErr != nil {
		return dBaseErr
	}

	tic, getTicErr := repo.GetTicInTx(headerID, tx)
	if getTicErr != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			log.Error("failed to rollback ", rollbackErr)
		}
		return getTicErr
	}

	for _, model := range models {
		dent, ok := model.(DentModel)
		if !ok {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Error("failed to rollback ", rollbackErr)
			}
			return fmt.Errorf("model of type %T, not %T", model, DentModel{})
		}

		_, execErr := tx.Exec(
			`INSERT into maker.dent (header_id, bid_id, lot, bid, guy, tic, log_idx, tx_idx, raw_log)
			VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (header_id, tx_idx, log_idx) DO UPDATE SET bid_Id = $2, lot = $3, bid = $4, guy = $5, tic = $6, raw_log = $9;`,
			headerID, dent.BidId, dent.Lot, dent.Bid, dent.Guy, tic, dent.LogIndex, dent.TransactionIndex, dent.Raw,
		)
		if execErr != nil {
			tx.Rollback()
			return execErr
		}
	}

	err := repo.MarkHeaderCheckedInTransaction(headerID, tx, constants.DentChecked)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (repository DentRepository) MarkHeaderChecked(headerId int64) error {
	return repo.MarkHeaderChecked(headerId, repository.db, constants.DentChecked)
}

func (repository DentRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return repo.MissingHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.DentChecked)
}

func (repository DentRepository) RecheckHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	return repo.RecheckHeaders(startingBlockNumber, endingBlockNumber, repository.db, constants.DentChecked)
}

func (repository *DentRepository) SetDB(db *postgres.DB) {
	repository.db = db
}

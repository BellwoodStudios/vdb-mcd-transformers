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

package tend

import (
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type Repository interface {
	Create(headerId int64, tends []TendModel) error
	MarkHeaderChecked(headerId int64) error
	MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error)
}

type TendRepository struct {
	DB *postgres.DB
}

func NewTendRepository(db *postgres.DB) TendRepository {
	return TendRepository{DB: db}
}

func (r TendRepository) Create(headerId int64, tends []TendModel) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	for _, tend := range tends {
		_, err = tx.Exec(
			`INSERT into maker.tend (header_id, bid_id, lot, bid, guy, tic, tx_idx, raw_log)
        VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
			headerId, tend.BidId, tend.Lot, tend.Bid, tend.Guy, tend.Tic, tend.TransactionIndex, tend.Raw,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	_, err = tx.Exec(`INSERT INTO public.checked_headers (header_id, tend_checked)
			VALUES ($1, $2)
		ON CONFLICT (header_id) DO
			UPDATE SET tend_checked = $2`, headerId, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r TendRepository) MarkHeaderChecked(headerId int64) error {
	_, err := r.DB.Exec(`INSERT INTO public.checked_headers (header_id, tend_checked)
		VALUES ($1, $2)
	ON CONFLICT (header_id) DO
		UPDATE SET tend_checked = $2`, headerId, true)
	return err
}

func (r TendRepository) MissingHeaders(startingBlockNumber, endingBlockNumber int64) ([]core.Header, error) {
	var result []core.Header
	err := r.DB.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               LEFT JOIN checked_headers on headers.id = header_id
               WHERE (header_id ISNULL OR tend_checked IS FALSE)
               AND headers.block_number >= $1
               AND headers.block_number <= $2
               AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		r.DB.Node.ID,
	)

	return result, err
}

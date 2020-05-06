package backfill

import (
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
)

type Frob struct {
	HeaderID int `db:"header_id"`
	Dink     string
	Dart     string
}

type EventsRepository interface {
	GetFrobs(urnID, startingBlock int) ([]Frob, error)
	GetHeaderByID(id int) (core.Header, error)
}

type eventsRepository struct {
	db *postgres.DB
}

func NewEventsRepository(db *postgres.DB) EventsRepository {
	return eventsRepository{db: db}
}

func (e eventsRepository) GetFrobs(urnID, startingBlock int) ([]Frob, error) {
	var frobs []Frob
	err := e.db.Select(&frobs, `SELECT header_id, dink, dart
		FROM maker.vat_frob
		JOIN public.headers ON vat_frob.header_id = headers.id
		WHERE urn_id = $1 AND headers.block_number >= $2
		ORDER BY headers.block_number ASC`, urnID, startingBlock)
	return frobs, err
}

func (e eventsRepository) GetHeaderByID(id int) (core.Header, error) {
	var header core.Header
	headerErr := e.db.Get(&header, `SELECT id, block_number, hash, raw, block_timestamp FROM headers WHERE id = $1`, id)
	return header, headerErr
}

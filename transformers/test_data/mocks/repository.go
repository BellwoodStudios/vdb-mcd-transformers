package mocks

import (
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockRepository struct {
	createError                      error
	markHeaderCheckedError           error
	MarkHeaderCheckedPassedHeaderIDs []int64
	CreatedHeaderIds                 []int64
	allHeaders                       []core.Header
	PassedHeaderID                   int64
	PassedModels                     []interface{}
	SetDbCalled                      bool
	CreateCalledCounter              int
}

func (repository *MockRepository) Create(headerID int64, models []interface{}) error {
	repository.PassedHeaderID = headerID
	repository.PassedModels = models
	repository.CreatedHeaderIds = append(repository.CreatedHeaderIds, headerID)
	repository.CreateCalledCounter++

	return repository.createError
}

func (repository *MockRepository) MarkHeaderChecked(headerID int64) error {
	repository.MarkHeaderCheckedPassedHeaderIDs = append(repository.MarkHeaderCheckedPassedHeaderIDs, headerID)
	return repository.markHeaderCheckedError
}

func (repository *MockRepository) SetDB(db *postgres.DB) {
	repository.SetDbCalled = true
}

func (repository *MockRepository) SetMarkHeaderCheckedError(e error) {
	repository.markHeaderCheckedError = e
}

func (repository *MockRepository) SetCreateError(e error) {
	repository.createError = e
}

func (repository *MockRepository) AssertMarkHeaderCheckedCalledWith(i int64) {
	Expect(repository.MarkHeaderCheckedPassedHeaderIDs).To(ContainElement(i))
}

func (repository *MockRepository) AssertMarkHeaderCheckedNotCalled() {
	Expect(len(repository.MarkHeaderCheckedPassedHeaderIDs)).To(Equal(0))
}

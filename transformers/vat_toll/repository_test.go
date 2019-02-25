package vat_toll_test

import (
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/mcd_transformers/transformers/test_data/shared_behaviors"
	"github.com/vulcanize/mcd_transformers/transformers/vat_toll"
)

var _ = Describe("Vat toll repository", func() {
	var (
		db         *postgres.DB
		repository vat_toll.VatTollRepository
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		repository = vat_toll.VatTollRepository{}
		repository.SetDB(db)
	})

	Describe("Create", func() {
		modelWithDifferentLogIdx := test_data.VatTollModel
		modelWithDifferentLogIdx.LogIndex++
		inputs := shared_behaviors.CreateBehaviorInputs{
			CheckedHeaderColumnName:  constants.VatTollChecked,
			LogEventTableName:        "maker.vat_toll",
			TestModel:                test_data.VatTollModel,
			ModelWithDifferentLogIdx: modelWithDifferentLogIdx,
			Repository:               &repository,
		}

		shared_behaviors.SharedRepositoryCreateBehaviors(&inputs)

		It("adds a vat toll event", func() {
			headerRepository := repositories.NewHeaderRepository(db)
			headerID, err := headerRepository.CreateOrUpdateHeader(fakes.FakeHeader)
			Expect(err).NotTo(HaveOccurred())
			err = repository.Create(headerID, []interface{}{test_data.VatTollModel})
			Expect(err).NotTo(HaveOccurred())

			var dbVatToll vat_toll.VatTollModel
			err = db.Get(&dbVatToll, `SELECT ilk, urn, take, tx_idx, log_idx, raw_log FROM maker.vat_toll WHERE header_id = $1`, headerID)
			Expect(err).NotTo(HaveOccurred())
			ilkID, err := shared.GetOrCreateIlk(test_data.VatTollModel.Ilk, db)
			Expect(err).NotTo(HaveOccurred())
			Expect(dbVatToll.Ilk).To(Equal(strconv.Itoa(ilkID)))
			Expect(dbVatToll.Urn).To(Equal(test_data.VatTollModel.Urn))
			Expect(dbVatToll.Take).To(Equal(test_data.VatTollModel.Take))
			Expect(dbVatToll.TransactionIndex).To(Equal(test_data.VatTollModel.TransactionIndex))
			Expect(dbVatToll.LogIndex).To(Equal(test_data.VatTollModel.LogIndex))
			Expect(dbVatToll.Raw).To(MatchJSON(test_data.VatTollModel.Raw))
		})
	})

	Describe("MarkHeaderChecked", func() {
		inputs := shared_behaviors.MarkedHeaderCheckedBehaviorInputs{
			CheckedHeaderColumnName: constants.VatTollChecked,
			Repository:              &repository,
		}

		shared_behaviors.SharedRepositoryMarkHeaderCheckedBehaviors(&inputs)
	})
})

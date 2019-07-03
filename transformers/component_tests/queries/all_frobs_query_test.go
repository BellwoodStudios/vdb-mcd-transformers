package queries

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/transformers/events/vat_frob"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
)

var _ = Describe("Frobs query", func() {
	var (
		db         *postgres.DB
		frobRepo   vat_frob.VatFrobRepository
		headerRepo repositories.HeaderRepository
		fakeUrn    = test_data.RandomString(5)
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		frobRepo = vat_frob.VatFrobRepository{}
		frobRepo.SetDB(db)
	})

	AfterEach(func() {
		closeErr := db.Close()
		Expect(closeErr).NotTo(HaveOccurred())
	})

	Describe("urn_frobs", func() {
		It("returns frobs for relevant ilk/urn", func() {
			headerOne := fakes.GetFakeHeaderWithTimestamp(int64(111111111), 1)

			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			frobBlockOne := test_helpers.CopyModel(test_data.VatFrobModelWithPositiveDart)
			frobBlockOne.ForeignKeyToValue["ilk_id"] = test_helpers.FakeIlk.Hex
			frobBlockOne.ForeignKeyToValue["urn_id"] = fakeUrn
			frobBlockOne.ColumnToValue["dink"] = strconv.Itoa(rand.Int())
			frobBlockOne.ColumnToValue["dart"] = strconv.Itoa(rand.Int())

			irrelevantFrob := test_helpers.CopyModel(test_data.VatFrobModelWithPositiveDart)
			irrelevantFrob.ForeignKeyToValue["ilk_id"] = test_helpers.AnotherFakeIlk.Hex
			irrelevantFrob.ForeignKeyToValue["urn_id"] = fakeUrn
			irrelevantFrob.ColumnToValue["dink"] = strconv.Itoa(rand.Int())
			irrelevantFrob.ColumnToValue["dart"] = strconv.Itoa(rand.Int())
			irrelevantFrob.ColumnToValue["tx_idx"] = frobBlockOne.ColumnToValue["tx_idx"].(uint) + 1

			err = frobRepo.Create(headerOneId, []shared.InsertionModel{frobBlockOne, irrelevantFrob})
			Expect(err).NotTo(HaveOccurred())

			// New block
			headerTwo := fakes.GetFakeHeaderWithTimestamp(int64(222222222), 2)
			headerTwo.Hash = "anotherHash"
			headerTwoId, err := headerRepo.CreateOrUpdateHeader(headerTwo)
			Expect(err).NotTo(HaveOccurred())

			frobBlockTwo := test_helpers.CopyModel(test_data.VatFrobModelWithPositiveDart)
			frobBlockTwo.ForeignKeyToValue["ilk_id"] = test_helpers.FakeIlk.Hex
			frobBlockTwo.ForeignKeyToValue["urn_id"] = fakeUrn
			frobBlockTwo.ColumnToValue["dink"] = strconv.Itoa(rand.Int())
			frobBlockTwo.ColumnToValue["dart"] = strconv.Itoa(rand.Int())

			err = frobRepo.Create(headerTwoId, []shared.InsertionModel{frobBlockTwo})
			Expect(err).NotTo(HaveOccurred())

			var actualFrobs []test_helpers.FrobEvent
			err = db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart FROM api.urn_frobs($1, $2)`, test_helpers.FakeIlk.Identifier, fakeUrn)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualFrobs).To(ConsistOf(
				test_helpers.FrobEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn,
					Dink: frobBlockOne.ColumnToValue["dink"].(string), Dart: frobBlockOne.ColumnToValue["dart"].(string)},
				test_helpers.FrobEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn,
					Dink: frobBlockTwo.ColumnToValue["dink"].(string), Dart: frobBlockTwo.ColumnToValue["dart"].(string)},
			))
		})

		It("fails if no argument is supplied (STRICT)", func() {
			_, err := db.Exec(`SELECT * FROM api.urn_frobs()`)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("function api.urn_frobs() does not exist"))
		})

		It("fails if only one argument is supplied (STRICT)", func() {
			_, err := db.Exec(`SELECT * FROM api.urn_frobs($1::text)`, test_helpers.FakeIlk.Identifier)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("function api.urn_frobs(text) does not exist"))
		})
	})

	Describe("all_frobs", func() {
		It("returns all frobs for a whole ilk", func() {
			headerOne := fakes.GetFakeHeader(1)

			headerOneId, err := headerRepo.CreateOrUpdateHeader(headerOne)
			Expect(err).NotTo(HaveOccurred())

			frobOne := test_helpers.CopyModel(test_data.VatFrobModelWithPositiveDart)
			frobOne.ForeignKeyToValue["ilk_id"] = test_helpers.FakeIlk.Hex
			frobOne.ForeignKeyToValue["urn_id"] = fakeUrn
			frobOne.ColumnToValue["dink"] = strconv.Itoa(rand.Int())
			frobOne.ColumnToValue["dart"] = strconv.Itoa(rand.Int())

			anotherUrn := "anotherUrn"
			frobTwo := test_helpers.CopyModel(test_data.VatFrobModelWithPositiveDart)
			frobTwo.ForeignKeyToValue["ilk_id"] = test_helpers.FakeIlk.Hex
			frobTwo.ForeignKeyToValue["urn_id"] = anotherUrn
			frobTwo.ColumnToValue["dink"] = strconv.Itoa(rand.Int())
			frobTwo.ColumnToValue["dart"] = strconv.Itoa(rand.Int())
			frobTwo.ColumnToValue["tx_idx"] = frobOne.ColumnToValue["tx_idx"].(uint) + 1

			err = frobRepo.Create(headerOneId, []shared.InsertionModel{frobOne, frobTwo})
			Expect(err).NotTo(HaveOccurred())

			var actualFrobs []test_helpers.FrobEvent
			err = db.Select(&actualFrobs, `SELECT ilk_identifier, urn_identifier, dink, dart FROM api.all_frobs($1)`, test_helpers.FakeIlk.Identifier)
			Expect(err).NotTo(HaveOccurred())

			Expect(actualFrobs).To(ConsistOf(
				test_helpers.FrobEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: fakeUrn,
					Dink: frobOne.ColumnToValue["dink"].(string), Dart: frobOne.ColumnToValue["dart"].(string)},
				test_helpers.FrobEvent{IlkIdentifier: test_helpers.FakeIlk.Identifier, UrnIdentifier: anotherUrn,
					Dink: frobTwo.ColumnToValue["dink"].(string), Dart: frobTwo.ColumnToValue["dart"].(string)},
			))
		})

		It("fails if no argument is supplied (STRICT)", func() {
			_, err := db.Exec(`SELECT * FROM api.all_frobs()`)
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("function api.all_frobs() does not exist"))
		})
	})
})

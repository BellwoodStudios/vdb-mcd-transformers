package queries_test

import (
	"database/sql"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/vulcanize/mcd_transformers/pkg/queries/test_helpers"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/storage/pit"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strconv"
	"time"
)

var _ = Describe("Single urn view", func() {
	var (
		db         *postgres.DB
		vatRepo    vat.VatStorageRepository
		pitRepo    pit.PitStorageRepository
		headerRepo repositories.HeaderRepository
		urnOne     string
		urnTwo     string
		ilkOne     string
		ilkTwo     string
		err        error
	)

	BeforeEach(func() {
		db = test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		vatRepo = vat.VatStorageRepository{}
		vatRepo.SetDB(db)
		pitRepo = pit.PitStorageRepository{}
		pitRepo.SetDB(db)
		headerRepo = repositories.NewHeaderRepository(db)
		rand.Seed(time.Now().UnixNano())

		urnOne = test_data.RandomString(5)
		urnTwo = test_data.RandomString(5)
		ilkOne = test_data.RandomString(5)
		ilkTwo = test_data.RandomString(5)
	})

	It("gets only the specified urn", func() {
		blockOne := rand.Int()
		timestampOne := rand.Int()
		blockTwo := blockOne + 1
		timestampTwo := timestampOne + 1

		urnOneMetadata := GetUrnMetadata(ilkOne, urnOne)
		urnOneSetupData := GetUrnSetupData(blockOne, timestampOne)
		CreateUrn(urnOneSetupData, urnOneMetadata, vatRepo, pitRepo, headerRepo)

		urnTwoMetadata := GetUrnMetadata(ilkTwo, urnTwo)
		urnTwoSetupData := GetUrnSetupData(blockTwo, timestampTwo)
		CreateUrn(urnTwoSetupData, urnTwoMetadata, vatRepo, pitRepo, headerRepo)

		var result UrnState
		err = db.Get(&result, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
			FROM get_urn_state_at_block($1, $2, $3)`, ilkOne, urnOne, blockTwo)
		Expect(err).NotTo(HaveOccurred())

		expectedRatio := GetExpectedRatio(urnOneSetupData.Ink, urnOneSetupData.Spot, urnOneSetupData.Art, urnOneSetupData.Rate)

		expectedUrn := UrnState{
			UrnId:   urnOne,
			IlkId:   ilkOne,
			Ink:     strconv.Itoa(urnOneSetupData.Ink),
			Art:     strconv.Itoa(urnOneSetupData.Art),
			Ratio:   sql.NullString{String: "", Valid: true}, // Checked separately, floating point arithmetic errors
			Safe:    expectedRatio >= 1,
			Created: sql.NullString{String: strconv.Itoa(timestampOne), Valid: true},
			Updated: sql.NullString{String: strconv.Itoa(timestampOne), Valid: true},
		}

		actualRatio, err := strconv.ParseFloat(result.Ratio.String, 64)
		Expect(err).NotTo(HaveOccurred())
		result.Ratio.String = ""

		Expect(result).To(Equal(expectedUrn))
		Expect(actualRatio).Should(BeNumerically("~", expectedRatio))
	})

	It("returns urn state without timestamps if corresponding headers aren't synced", func() {
		block := rand.Int()
		timestamp := rand.Int()
		metadata := GetUrnMetadata(ilkOne, urnOne)
		setupData := GetUrnSetupData(block, timestamp)

		CreateUrn(setupData, metadata, vatRepo, pitRepo, headerRepo)
		_, err = db.Exec(`DELETE FROM headers`)
		Expect(err).NotTo(HaveOccurred())

		var result UrnState
		err = db.Get(&result, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
			FROM get_urn_state_at_block($1, $2, $3)`, ilkOne, urnOne, block)

		Expect(err).NotTo(HaveOccurred())
		Expect(result.Created.String).To(BeEmpty())
		Expect(result.Updated.String).To(BeEmpty())
	})

	Describe("it includes diffs only up to given block height", func() {
		var (
			result       UrnState
			blockOne     int
			timestampOne int
			setupDataOne UrnSetupData
			metadata     UrnMetadata
		)

		BeforeEach(func() {
			blockOne = rand.Int()
			timestampOne = rand.Int()
			setupDataOne = GetUrnSetupData(blockOne, timestampOne)
			metadata = GetUrnMetadata(ilkOne, urnOne)
			CreateUrn(setupDataOne, metadata, vatRepo, pitRepo, headerRepo)
		})

		It("gets urn state as of block one", func() {
			blockTwo := blockOne + 1
			hashTwo := test_data.RandomString(5)
			updatedInk := rand.Int()

			err = vatRepo.Create(blockTwo, hashTwo, metadata.UrnInk, strconv.Itoa(updatedInk))
			Expect(err).NotTo(HaveOccurred())

			err = db.Get(&result, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
				FROM get_urn_state_at_block($1, $2, $3)`, ilkOne, urnOne, blockOne)
			Expect(err).NotTo(HaveOccurred())

			expectedRatio := GetExpectedRatio(setupDataOne.Ink, setupDataOne.Spot, setupDataOne.Art, setupDataOne.Rate)

			expectedUrn := UrnState{
				UrnId:   urnOne,
				IlkId:   ilkOne,
				Ink:     strconv.Itoa(setupDataOne.Ink),
				Art:     strconv.Itoa(setupDataOne.Art),
				Ratio:   sql.NullString{String: "", Valid: true}, // Checked separately, floating point arithmetic errors
				Safe:    expectedRatio >= 1,
				Created: sql.NullString{String: strconv.Itoa(timestampOne), Valid: true},
				Updated: sql.NullString{String: strconv.Itoa(timestampOne), Valid: true},
			}

			actualRatio, err := strconv.ParseFloat(result.Ratio.String, 64)
			Expect(err).NotTo(HaveOccurred())
			result.Ratio.String = ""

			Expect(result).To(Equal(expectedUrn))
			Expect(actualRatio).Should(BeNumerically("~", expectedRatio))
		})

		It("gets urn state with updated values", func() {
			blockTwo := blockOne + 1
			timestampTwo := timestampOne + 1
			hashTwo := test_data.RandomString(5)
			updatedInk := rand.Int()

			err = vatRepo.Create(blockTwo, hashTwo, metadata.UrnInk, strconv.Itoa(updatedInk))
			Expect(err).NotTo(HaveOccurred())

			fakeHeaderTwo := fakes.GetFakeHeader(int64(blockTwo))
			fakeHeaderTwo.Timestamp = strconv.Itoa(timestampTwo)
			fakeHeaderTwo.Hash = hashTwo

			_, err = headerRepo.CreateOrUpdateHeader(fakeHeaderTwo)
			Expect(err).NotTo(HaveOccurred())

			err = db.Get(&result, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
				FROM get_urn_state_at_block($1, $2, $3)`, ilkOne, urnOne, blockTwo)
			Expect(err).NotTo(HaveOccurred())

			expectedRatio := GetExpectedRatio(updatedInk, setupDataOne.Spot, setupDataOne.Art, setupDataOne.Rate)

			expectedUrn := UrnState{
				UrnId:   urnOne,
				IlkId:   ilkOne,
				Ink:     strconv.Itoa(updatedInk),
				Art:     strconv.Itoa(setupDataOne.Art),          // Not changed
				Ratio:   sql.NullString{String: "", Valid: true}, // Checked separately, floating point arithmetic errors
				Safe:    expectedRatio >= 1,
				Created: sql.NullString{String: strconv.Itoa(timestampOne), Valid: true},
				Updated: sql.NullString{String: strconv.Itoa(timestampTwo), Valid: true},
			}

			actualRatio, err := strconv.ParseFloat(result.Ratio.String, 64)
			Expect(err).NotTo(HaveOccurred())
			result.Ratio.String = ""

			Expect(result).To(Equal(expectedUrn))
			Expect(actualRatio).Should(BeNumerically("~", expectedRatio))
		})
	})

	It("returns null ratio and urn being safe if there is no debt", func() {
		block := rand.Int()
		setupData := GetUrnSetupData(block, 1)
		setupData.Art = 0
		metadata := GetUrnMetadata(ilkOne, urnOne)
		CreateUrn(setupData, metadata, vatRepo, pitRepo, headerRepo)

		fakeHeader := fakes.GetFakeHeader(int64(block))
		_, err = headerRepo.CreateOrUpdateHeader(fakeHeader)
		Expect(err).NotTo(HaveOccurred())

		var result UrnState
		err = db.Get(&result, `SELECT urnId, ilkId, ink, art, ratio, safe, created, updated
			FROM get_urn_state_at_block($1, $2, $3)`, ilkOne, urnOne, block)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Ratio.String).To(BeEmpty())
		Expect(result.Safe).To(BeTrue())
	})
})
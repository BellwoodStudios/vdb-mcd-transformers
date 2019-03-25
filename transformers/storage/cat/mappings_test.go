package cat_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/vulcanize/vulcanizedb/libraries/shared/storage"
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"

	"github.com/vulcanize/mcd_transformers/transformers/shared/constants"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/test_helpers"
)

var _ = Describe("Cat storage mappings", func() {
	const (
		fakeIlk  = "fakeIlk"
		fakeFlip = "2"
	)

	var (
		storageRepository *test_helpers.MockMakerStorageRepository
		mappings          cat.CatMappings
	)

	BeforeEach(func() {
		storageRepository = &test_helpers.MockMakerStorageRepository{}
		mappings = cat.CatMappings{StorageRepository: storageRepository}
	})

	Describe("looking up static keys", func() {
		It("returns value metadata if key exists", func() {
			Expect(mappings.Lookup(cat.NFlipKey)).To(Equal(cat.NFlipMetadata))
			Expect(mappings.Lookup(cat.LiveKey)).To(Equal(cat.LiveMetadata))
			Expect(mappings.Lookup(cat.VatKey)).To(Equal(cat.VatMetadata))
			Expect(mappings.Lookup(cat.VowKey)).To(Equal(cat.VowMetadata))
		})

		It("returns error if key does not exist", func() {
			_, err := mappings.Lookup(common.HexToHash(fakes.FakeHash.Hex()))

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(utils.ErrStorageKeyNotFound{Key: fakes.FakeHash.Hex()}))
		})
	})

	Describe("looking up dynamic keys", func() {
		It("refreshes mappings from repository if key not found", func() {
			_, _ = mappings.Lookup(fakes.FakeHash)

			Expect(storageRepository.GetIlksCalled).To(BeTrue())
			Expect(storageRepository.GetMaxFlipCalled).To(BeTrue())
		})

		It("returns error if ilks lookup fails", func() {
			storageRepository.GetIlksError = fakes.FakeError

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("returns error if max flip lookup fails", func() {
			storageRepository.GetMaxFlipError = fakes.FakeError

			_, err := mappings.Lookup(fakes.FakeHash)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(fakes.FakeError))
		})

		It("interpolates flips up to max", func() {
			storageRepository.MaxFlip = big.NewInt(1)

			_, err := mappings.Lookup(storage.GetMapping(storage.IndexTwo, "0"))
			Expect(err).NotTo(HaveOccurred())

			_, err = mappings.Lookup(storage.GetMapping(storage.IndexTwo, "1"))
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("ilk", func() {
			var ilkFlipKey = common.BytesToHash(crypto.Keccak256(common.FromHex(fakeIlk + cat.IlksMappingIndex)))

			BeforeEach(func() {
				storageRepository.Ilks = []string{fakeIlk}
			})

			It("returns value metadata for ilk flip", func() {
				expectedMetadata := utils.StorageValueMetadata{
					Name: cat.IlkFlip,
					Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
					Type: utils.Address,
				}
				Expect(mappings.Lookup(ilkFlipKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk chop", func() {
				ilkChopKey := storage.GetIncrementedKey(ilkFlipKey, 1)
				expectedMetadata := utils.StorageValueMetadata{
					Name: cat.IlkChop,
					Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
					Type: utils.Uint256,
				}
				Expect(mappings.Lookup(ilkChopKey)).To(Equal(expectedMetadata))
			})

			It("returns value metadata for ilk lump", func() {
				ilkLumpKey := storage.GetIncrementedKey(ilkFlipKey, 2)
				expectedMetadata := utils.StorageValueMetadata{
					Name: cat.IlkLump,
					Keys: map[utils.Key]string{constants.Ilk: fakeIlk},
					Type: utils.Uint256,
				}
				Expect(mappings.Lookup(ilkLumpKey)).To(Equal(expectedMetadata))
			})
		})

		Describe("flip", func() {
			var flipIlkKey = common.BytesToHash(crypto.Keccak256(common.FromHex(fakeFlip + cat.FlipsMappingIndex)))

			BeforeEach(func() {
				storageRepository.MaxFlip = big.NewInt(2)
			})

			It("returns value metadata for flip ilk", func() {
				expectedMetadata := utils.StorageValueMetadata{
					Name: cat.FlipIlk,
					Keys: map[utils.Key]string{constants.Flip: fakeFlip},
					Type: utils.Bytes32,
				}
				actualMetadata, err := mappings.Lookup(flipIlkKey)
				Expect(err).NotTo(HaveOccurred())
				Expect(actualMetadata).To(Equal(expectedMetadata))
			})

			It("returns value metadata for flip urn", func() {
				flipUrnKey := storage.GetIncrementedKey(flipIlkKey, 1)
				expectedMetadata := utils.StorageValueMetadata{
					Name: cat.FlipUrn,
					Keys: map[utils.Key]string{constants.Flip: fakeFlip},
					Type: utils.Bytes32,
				}
				actualMetadata, err := mappings.Lookup(flipUrnKey)
				Expect(err).NotTo(HaveOccurred())
				Expect(actualMetadata).To(Equal(expectedMetadata))
			})

			It("returns value metadata for flip ink", func() {
				flipInkKey := storage.GetIncrementedKey(flipIlkKey, 2)
				expectedMetadata := utils.StorageValueMetadata{
					Name: cat.FlipInk,
					Keys: map[utils.Key]string{constants.Flip: fakeFlip},
					Type: utils.Uint256,
				}
				actualMetadata, err := mappings.Lookup(flipInkKey)
				Expect(err).NotTo(HaveOccurred())
				Expect(actualMetadata).To(Equal(expectedMetadata))
			})

			It("returns value metadata for flip tab", func() {
				flipTabKey := storage.GetIncrementedKey(flipIlkKey, 3)
				expectedMetadata := utils.StorageValueMetadata{
					Name: cat.FlipTab,
					Keys: map[utils.Key]string{constants.Flip: fakeFlip},
					Type: utils.Uint256,
				}
				actualMetadata, err := mappings.Lookup(flipTabKey)
				Expect(err).NotTo(HaveOccurred())
				Expect(actualMetadata).To(Equal(expectedMetadata))
			})
		})
	})
})

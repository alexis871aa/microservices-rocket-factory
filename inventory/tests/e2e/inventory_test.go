package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
)

var _ = Describe("Inventory Service", func() {
	var testCtx context.Context

	BeforeEach(func() {
		// создаем контекст для каждого теста
		testCtx = suiteCtx

		// очищаем коллекцию перед каждым тестом
		err := env.ClearInventoryCollection(testCtx)
		Expect(err).NotTo(HaveOccurred())

		// вставляем базовые тестовые данные
		err = env.InsertTestPart(testCtx, RocketEngine)
		Expect(err).NotTo(HaveOccurred())

		err = env.InsertTestPart(testCtx, FuelTank)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("GetPart 🚀", func() {
		It("должен успешно вернуть деталь по существующему UUID", func() {
			// выполняем запрос
			resp, err := inventoryClient.GetPart(testCtx, &inventoryV1.GetPartRequest{
				Uuid: RocketEngine.UUID,
			})

			// проверяем результат
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.Part).NotTo(BeNil())

			// проверяем данные детали
			part := resp.Part
			Expect(part.Uuid).To(Equal(RocketEngine.UUID))
			Expect(part.Name).To(Equal(RocketEngine.Name))
			Expect(part.Description).To(Equal(RocketEngine.Description))
			Expect(part.Price).To(Equal(RocketEngine.Price))
			Expect(part.StockQuantity).To(Equal(RocketEngine.StockQuantity))
			Expect(part.Category).To(Equal(RocketEngine.Category))
		})

		It("должен вернуть ошибку NotFound для несуществующего UUID", func() {
			// выполняем запрос с несуществующим UUID
			resp, err := inventoryClient.GetPart(testCtx, &inventoryV1.GetPartRequest{
				Uuid: "non-existent-uuid",
			})

			// проверяем что вернулась ошибка NotFound
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())

			// проверяем тип ошибки
			grpcStatus, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(grpcStatus.Code()).To(Equal(codes.NotFound))
		})

		It("должен вернуть ошибку для пустого UUID", func() {
			// выполняем запрос с пустым UUID
			resp, err := inventoryClient.GetPart(testCtx, &inventoryV1.GetPartRequest{
				Uuid: "",
			})

			// проверяем результат
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
		})
	})

	Context("ListParts 📦", func() {
		It("должен вернуть все детали без фильтра", func() {
			// выполняем запрос без фильтра
			resp, err := inventoryClient.ListParts(testCtx, &inventoryV1.ListPartsRequest{})

			// проверяем результат
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.Parts).To(HaveLen(2)) // RocketEngine + FuelTank

			// проверяем что вернулись наши тестовые детали
			partUUIDs := make([]string, len(resp.Parts))
			for i, part := range resp.Parts {
				partUUIDs[i] = part.Uuid
			}
			Expect(partUUIDs).To(ContainElements(RocketEngine.UUID, FuelTank.UUID))
		})

		It("должен фильтровать детали по категории ENGINE", func() {
			// выполняем запрос с фильтром по категории ENGINE
			resp, err := inventoryClient.ListParts(testCtx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Categories: []inventoryV1.Category{inventoryV1.Category_ENGINE},
				},
			})

			// проверяем результат
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.Parts).To(HaveLen(1))

			// проверяем что вернулся только двигатель
			part := resp.Parts[0]
			Expect(part.Uuid).To(Equal(RocketEngine.UUID))
			Expect(part.Category).To(Equal(inventoryV1.Category_ENGINE))
		})

		It("должен фильтровать детали по списку UUID", func() {
			// выполняем запрос с фильтром по конкретному UUID
			resp, err := inventoryClient.ListParts(testCtx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Uuids: []string{FuelTank.UUID},
				},
			})

			// проверяем результат
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.Parts).To(HaveLen(1))

			// проверяем что вернулся только топливный бак
			part := resp.Parts[0]
			Expect(part.Uuid).To(Equal(FuelTank.UUID))
			Expect(part.Category).To(Equal(inventoryV1.Category_FUEL))
		})

		It("должен вернуть пустой список для несуществующей категории", func() {
			// выполняем запрос с фильтром по несуществующей категории
			resp, err := inventoryClient.ListParts(testCtx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Categories: []inventoryV1.Category{inventoryV1.Category_WING},
				},
			})

			// проверяем результат
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.Parts).To(HaveLen(0))
		})
	})

	Context("Интеграция с MongoDB 🗄️", func() {
		It("должен корректно работать с базой данных", func() {
			// проверяем что детали действительно есть в БД
			exists, err := env.PartExists(testCtx, RocketEngine.UUID)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())

			// проверяем общее количество
			count, err := env.GetPartsCount(testCtx)
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(int64(2)))

			// проверяем количество по категории
			engineCount, err := env.GetPartsByCategory(testCtx, inventoryV1.Category_ENGINE)
			Expect(err).NotTo(HaveOccurred())
			Expect(engineCount).To(Equal(int64(1)))
		})
	})
})

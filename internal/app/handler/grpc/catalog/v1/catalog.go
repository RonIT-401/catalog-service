package ghcatalogv1

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/RonIT-401/catalog-service/internal/app/entity"
	"github.com/RonIT-401/catalog-service/internal/app/mapper"
	"github.com/RonIT-401/catalog-service/internal/app/mapper/catalog/v1"
	"github.com/RonIT-401/catalog-service/internal/app/service"
	catalogv1 "github.com/RonIT-401/catalog-service/internal/pkg/grpc/gen/catalog/v1"
)

type handler struct {
	catalogv1.UnimplementedCatalogServiceServer
	srv service.Product
}

func NewHandler(srv service.Product) catalogv1.CatalogServiceServer {
	return &handler{srv: srv}
}

func (h *handler) GetProduct(ctx context.Context, req *catalogv1.GetProductRequest) (*catalogv1.GetProductResponse, error) {
	rawGUID := req.GetGuid()

	guid, err := uuid.FromString(rawGUID)
	if err != nil {
		return nil, mapper.ErrorToGRPC(entity.ErrIncorrectParameters)
	}

	products, err := h.srv.GetByGUIDs(ctx, []uuid.UUID{guid})
	if err != nil {
		return nil, mapper.ErrorToGRPC(err)
	}

	if len(products) == 0 {
		return nil, mapper.ErrorToGRPC(entity.ErrNotFound)
	}

	return &catalogv1.GetProductResponse{
		Product: mcatv1.ProductToProto(products[0]),
	}, nil
}

func (h *handler) GetProducts(ctx context.Context, req *catalogv1.GetProductsRequest) (*catalogv1.GetProductsResponse, error) {
	rawGUIDs := req.GetGuids()

	if len(rawGUIDs) == 0 {
		return &catalogv1.GetProductsResponse{}, nil
	}

	guids := make([]uuid.UUID, 0, len(rawGUIDs))
	for _, rawGUID := range rawGUIDs {
		guid, err := uuid.FromString(rawGUID)
		if err != nil {
			return nil, mapper.ErrorToGRPC(entity.ErrIncorrectParameters)
		}

		guids = append(guids, guid)
	}

	products, err := h.srv.GetByGUIDs(ctx, guids)
	if err != nil {
		return nil, mapper.ErrorToGRPC(err)
	}

	result := &catalogv1.GetProductsResponse{}

	for _, product := range products {
		result.Products = append(result.Products, mcatv1.ProductToProto(product))
	}

	found := make(map[uuid.UUID]struct{}, len(products))

	for _, p := range products {
		found[p.GUID] = struct{}{}
	}

	for i, guid := range guids {
		if _, ok := found[guid]; !ok {
			result.MissingGuids = append(
				result.MissingGuids,
				rawGUIDs[i],
			)
		}
	}
	return result, nil
}

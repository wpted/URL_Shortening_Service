package gRPC

import (
	"KeyGenerationService/internal/controller"
	"KeyGenerationService/internal/handler/gRPC/gen"
	"context"
)

// Handler implements the generated gRPC server and is responsible for accepting all incoming GetKeyMetadata requests.
type Handler struct {
	gen.UnimplementedKeyGenerationServiceServer
	controller *controller.KGS
}

// New creates a new handler instance.
func New(ctrl *controller.KGS) *Handler {
	return &Handler{controller: ctrl}
}

// GetKeyMetadata accepts all incoming gen.GetKeyMetadataRequest and fetches keys from the database.
func (h *Handler) GetKeyMetadata(ctx context.Context, req *gen.GetKeyMetadataRequest) (*gen.GetKeyMetadataResponse, error) {
	keys, err := h.controller.GetKeys(int(req.RequiredKeys))
	if err != nil {
		// TODO: Handle GetKeys error
		return &gen.GetKeyMetadataResponse{Success: false}, err
	}
	return &gen.GetKeyMetadataResponse{Keys: keys, Success: true}, nil
}

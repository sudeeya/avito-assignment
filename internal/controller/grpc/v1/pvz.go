package v1

import (
	context "context"

	"github.com/sudeeya/avito-assignment/internal/model"
	"github.com/sudeeya/avito-assignment/internal/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type pvzServiceServerImplementation struct {
	UnimplementedPVZServiceServer

	services *service.Services
}

func NewPVZServiceServerImplementation(services *service.Services) *pvzServiceServerImplementation {
	return &pvzServiceServerImplementation{
		services: services,
	}
}

func (p *pvzServiceServerImplementation) GetPVZList(ctx context.Context, req *GetPVZListRequest) (*GetPVZListResponse, error) {
	var res GetPVZListResponse

	pvzs, err := p.services.PVZ.GetPVZList(ctx)
	if err != nil {
		return nil, err
	}

	for _, pvz := range pvzs {
		res.Pvzs = append(res.Pvzs, toProto(pvz))
	}

	return &res, nil
}

func toProto(pvz model.PVZ) *PVZ {
	return &PVZ{
		Id:               pvz.ID.String(),
		RegistrationDate: timestamppb.New(pvz.RegistrationDate),
		City:             pvz.City,
	}
}

package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/alexis871aa/microservices-rocket-factory/notification/internal/model"
	eventsV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/events/v1"
)

type shipAssembledDecoder struct{}

func NewShipAssembledDecoder() *shipAssembledDecoder {
	return &shipAssembledDecoder{}
}

func (d *shipAssembledDecoder) Decode(data []byte) (model.ShipAssembled, error) {
	var pb eventsV1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembled{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.ShipAssembled{
		EventUUID:    pb.EventUuid,
		OrderUUID:    pb.OrderUuid,
		UserUUID:     pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}

package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	pb "ride-sharing/shared/proto/trip"
)

type RideFareModel struct {
	ID                primitive.ObjectID
	UserID            string
	PackageSlug       string // eg. van, luxury, sedan
	TotalPriceInCents float64
}

func (r *RideFareModel) ToProto() *pb.RideFare {
	return &pb.RideFare{
		Id:                r.ID.Hex(),
		UserID:            r.UserID,
		PackageSlug:       r.PackageSlug,
		TotalPriceInCents: r.TotalPriceInCents,
	}
}

func ToRideFaresProto(fares []*RideFareModel) []*pb.RideFare {
	var protoFares []*pb.RideFare

	for _, fare := range fares {
		protoFares = append(protoFares, fare.ToProto())
	}

	return protoFares
}

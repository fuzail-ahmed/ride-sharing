package main

import "ride-sharing/shared/types"

type previewTripRequest struct {
	UserID      string `json:"userID"`
	PickUp      types.Coordinate
	Destination types.Coordinate
}

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	t := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
		Driver:   &trip.TripDriver{},
	}

	return s.repo.CreateTrip(ctx, t)
}

func (s *service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripTypes.OsrmApiResponse, error) {
	url := fmt.Sprintf(
		"http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude,
	)

	log.Printf("Fetching from OSRM API: URL: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route from OSRM API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %v", err)
	}

	log.Printf("GOT RESPONSE FROM API %s", string(body))

	var routeResp tripTypes.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &routeResp, nil
}

func (s *service) EstimatePackagesPriceWithRoute(route *tripTypes.OsrmApiResponse) []*domain.RideFareModel {
	baseFare := getBaseFares()

	estimatedFares := make([]*domain.RideFareModel, len(baseFare))

	for i, fare := range estimatedFares {
		estimatedFares[i] = estimatedFareRoute(fare, route)
	}

	return estimatedFares
}

func (s *service) GenerateTripFares(ctx context.Context, rideFares []*domain.RideFareModel, userID string, route *tripTypes.OsrmApiResponse) ([]*domain.RideFareModel, error) {
	fares := make([]*domain.RideFareModel, len(rideFares))

	for i, fare := range rideFares {
		id := primitive.NewObjectID()

		newFare := &domain.RideFareModel{
			ID:                id,
			UserID:            userID,
			TotalPriceInCents: fare.TotalPriceInCents,
			PackageSlug:       fare.PackageSlug,
			Route:             route,
		}

		if err := s.repo.SaveRideFare(ctx, newFare); err != nil {
			return nil, fmt.Errorf("failed to save trip fare: %w", err)
		}

		fares[i] = fare
	}

	return fares, nil
}

func (s *service) GetAndValidateFare(ctx context.Context, fareID, userID string) (*domain.RideFareModel, error) {
	fare, err := s.repo.GetRideFareByID(ctx, fareID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trip fare: %w", err)
	}

	if fare == nil {
		return nil, fmt.Errorf("fare does not exist")
	}

	// User fare validation (user is owner of this fare?)
	if userID != fare.UserID {
		return nil, fmt.Errorf("fare does not belong to the user")
	}

	return fare, nil
}

// Required functions
func estimatedFareRoute(fare *domain.RideFareModel, route *tripTypes.OsrmApiResponse) *domain.RideFareModel {
	pricingConfig := tripTypes.DefaultPricingConfig()
	carPackagePrice := fare.TotalPriceInCents

	distanceInKm := route.Routes[0].Distance
	durationInMinutes := route.Routes[0].Duration

	distanceFare := distanceInKm * pricingConfig.PricePerUnitOfDistance
	timeFare := durationInMinutes * pricingConfig.PricingPerMinute
	totalPrice := distanceFare + timeFare + carPackagePrice

	return &domain.RideFareModel{
		TotalPriceInCents: totalPrice,
		PackageSlug:       fare.PackageSlug,
	}
}

func getBaseFares() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{
			PackageSlug:       "suv",
			TotalPriceInCents: 200,
		},
		{
			PackageSlug:       "sedan",
			TotalPriceInCents: 350,
		},
		{
			PackageSlug:       "van",
			TotalPriceInCents: 400,
		},
		{
			PackageSlug:       "luxury",
			TotalPriceInCents: 1000,
		},
	}
}

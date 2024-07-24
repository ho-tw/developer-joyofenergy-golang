package readings

import (
	"joi-energy-golang/domain"
	"joi-energy-golang/repository"
	"time"
)

type Service interface {
	StoreReadings(smartMeterId string, reading []domain.ElectricityReading)
	GetReadings(smartMeterId string) []domain.ElectricityReading
	CalculateLastWeekUsageCost(smartMeterId string, tariff float64) (float64, error)
}

type service struct {
	meterReadings *repository.MeterReadings
}

func NewService(
	meterReadings *repository.MeterReadings,
) Service {
	return &service{
		meterReadings: meterReadings,
	}
}

func (s *service) StoreReadings(smartMeterId string, reading []domain.ElectricityReading) {
	s.meterReadings.StoreReadings(smartMeterId, reading)
}

func (s *service) GetReadings(smartMeterId string) []domain.ElectricityReading {
	return s.meterReadings.GetReadings(smartMeterId)
}

func (s *service) CalculateLastWeekUsageCost(smartMeterId string, tariff float64) (float64, error) {
	readings := s.GetReadings(smartMeterId)

	if len(readings) == 0 {
		return 0, nil // No cost if there are no readings
	}

	var sumReadings float64
	lastWeekReadings := []domain.ElectricityReading{}
	for _, reading := range readings {
		if reading.Time.Before(time.Now().AddDate(0, 0, -7)) {
			continue
		}
		sumReadings += reading.Reading
		lastWeekReadings = append(lastWeekReadings, reading)
	}

	averageReading := sumReadings / float64(len(lastWeekReadings))
	usageTimeHours := 168.0 // Hours in a week
	energyConsumed := averageReading * usageTimeHours
	cost := tariff * energyConsumed

	return cost, nil
}

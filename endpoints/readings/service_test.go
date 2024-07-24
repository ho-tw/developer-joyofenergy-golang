package readings

import (
	"testing"
	"time"

	"joi-energy-golang/domain"
	"joi-energy-golang/repository"
)

func TestStoreReadings(t *testing.T) {
	meterReadings := repository.NewMeterReadings(
		map[string][]domain.ElectricityReading{},
	)
	service := NewService(
		&meterReadings,
	)
	service.StoreReadings("1", []domain.ElectricityReading{})
}

func CalculateLastWeekUsageCost(t *testing.T) {
	expectedSmartMeterId := "1"
	tests := []struct {
		name         string
		smartMeterId string
		readings     []domain.ElectricityReading
		want         float64
	}{
		{
			name:         "no readings",
			smartMeterId: expectedSmartMeterId,
			readings:     []domain.ElectricityReading{},
			want:         0,
		},
		{
			name:         "single reading",
			smartMeterId: expectedSmartMeterId,
			readings: []domain.ElectricityReading{
				{Time: time.Now(), Reading: 1},
			},
			want: 168,
		},
		{
			name:         "multiple readings",
			smartMeterId: expectedSmartMeterId,
			readings: []domain.ElectricityReading{
				{Time: time.Now(), Reading: 5},
				{Time: time.Now().AddDate(0, 0, -14), Reading: 10},
				{Time: time.Now(), Reading: 10},
				{Time: time.Now(), Reading: 15},
			},
			want: 1680,
		},
		{
			name:         "smart meter not found",
			smartMeterId: "2",
			readings: []domain.ElectricityReading{
				{Time: time.Now(), Reading: 5},
				{Time: time.Now(), Reading: 10},
				{Time: time.Now(), Reading: 15},
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meterReadings := repository.NewMeterReadings(
				map[string][]domain.ElectricityReading{
					expectedSmartMeterId: tt.readings,
				},
			)
			service := NewService(
				&meterReadings,
			)

			got, err := service.CalculateLastWeekUsageCost(tt.smartMeterId, 1)
			if err != nil {
				t.Errorf("CalculateAverageReading() error = %v", err)
				return

			}
			if got != tt.want {
				t.Errorf("CalculateAverageReading() = %v, want %v", got, tt.want)
			}
		})
	}
}

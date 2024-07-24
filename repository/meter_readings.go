package repository

import "joi-energy-golang/domain"

type MeterReadings struct {
	meterAssociatedReadings map[string][]domain.ElectricityReading
}

func NewMeterReadings(meterAssociatedReadings map[string][]domain.ElectricityReading) MeterReadings {
	return MeterReadings{meterAssociatedReadings: meterAssociatedReadings}
}

func (m *MeterReadings) GetReadings(smartMeterId string) []domain.ElectricityReading {
	v, ok := m.meterAssociatedReadings[smartMeterId]
	if !ok {
		return nil
	}
	return v
}

func (m *MeterReadings) StoreReadings(smartMeterId string, electricityReadings []domain.ElectricityReading) {
	m.meterAssociatedReadings[smartMeterId] = append(m.meterAssociatedReadings[smartMeterId], electricityReadings...)
}

func (m *MeterReadings) CalculateLastWeekUsageCost(smartMeterId string) (float64, error) {
	// Placeholder implementation: calculate the cost based on last week's readings.
	// You'll need to define how to determine which readings are from the last week and how to calculate the cost.
	// This is a simplified example.
	readings := m.GetReadings(smartMeterId)
	var cost float64
	for _, reading := range readings {
		// Assuming each reading has a 'Timestamp' and 'Value' you can use to calculate the cost.
		// This is just a placeholder calculation.
		// Replace it with your actual logic for calculating the cost based on readings.
		cost += reading.Reading // Simplified; likely involves more complex calculations.
	}
	return cost, nil
}

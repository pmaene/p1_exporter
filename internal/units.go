package internal

import (
	"fmt"
	"strconv"

	"github.com/skoef/gop1"
)

type ElectricCurrent float64

func ParseElectricCurrent(v gop1.TelegramValue) (ElectricCurrent, error) {
	if v.Unit != "A" {
		return 0, fmt.Errorf("unknown energy unit %v", v.Unit)
	}

	u, err := strconv.ParseFloat(v.Value, 64)
	if err != nil {
		return 0, err
	}

	return ElectricCurrent(u), nil
}

type Energy float64

func ParseEnergy(v gop1.TelegramValue) (Energy, error) {
	if v.Unit != "kWh" {
		return 0, fmt.Errorf("unknown energy unit %v", v.Unit)
	}

	u, err := strconv.ParseFloat(v.Value, 64)
	if err != nil {
		return 0, err
	}

	return Energy(1000 * u), nil
}

type Power float64

func ParsePower(v gop1.TelegramValue) (Power, error) {
	if v.Unit != "kW" {
		return 0, fmt.Errorf("unknown power unit %v", v.Unit)
	}

	u, err := strconv.ParseFloat(v.Value, 64)
	if err != nil {
		return 0, err
	}

	return Power(1000 * u), nil
}

type Voltage float64

func ParseVoltage(v gop1.TelegramValue) (Voltage, error) {
	if v.Unit != "V" {
		return 0, fmt.Errorf("unknown voltage unit %v", v.Unit)
	}

	u, err := strconv.ParseFloat(v.Value, 64)
	if err != nil {
		return 0, err
	}

	return Voltage(u), nil
}

type Volume float64

func ParseVolume(v gop1.TelegramValue) (Volume, error) {
	if v.Unit != "m3" {
		return 0, fmt.Errorf("unknown volume unit %v", v.Unit)
	}

	u, err := strconv.ParseFloat(v.Value, 64)
	if err != nil {
		return 0, err
	}

	return Volume(u), nil
}

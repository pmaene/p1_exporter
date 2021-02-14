package internal

import (
	"errors"
	"strconv"
	"time"

	"github.com/skoef/gop1"
)

var (
	ErrInvalidTimestampSeason = errors.New("invalid timestamp season")
	ErrUnknownBreakerState    = errors.New("unknown breaker state")
	ErrUnknownGasValveState   = errors.New("unknown gas valve state")
)

type BreakerState int

const (
	BreakerStateDisconnected         = 0
	BreakerStateConnected            = 1
	BreakerStateReadyForReconnection = 2
)

func ParseBreakerState(v gop1.TelegramValue) (BreakerState, error) {
	u, err := strconv.Atoi(v.Value)
	if err != nil {
		return 0, err
	}

	switch u {
	case 0:
		return BreakerStateDisconnected, nil
	case 1:
		return BreakerStateConnected, nil
	case 2:
		return BreakerStateReadyForReconnection, nil
	default:
		return 0, ErrUnknownBreakerState
	}
}

type GasValveState int

const (
	GasValveStateDisconnected         = 0
	GasValveStateConnected            = 1
	GasValveStateReadyForReconnection = 2
)

func ParseGasValveState(v gop1.TelegramValue) (GasValveState, error) {
	u, err := strconv.Atoi(v.Value)
	if err != nil {
		return 0, err
	}

	switch u {
	case 0:
		return GasValveStateDisconnected, nil
	case 1:
		return GasValveStateConnected, nil
	case 2:
		return GasValveStateReadyForReconnection, nil
	default:
		return 0, ErrUnknownGasValveState
	}
}

func ParseElectricityTariffIndicator(v gop1.TelegramValue) (int, error) {
	return strconv.Atoi(v.Value)
}

func ParseTimestamp(v gop1.TelegramValue) (time.Time, error) {
	t := v.Value[:len(v.Value)-1]
	s := v.Value[len(v.Value)-1:]

	switch s {
	case "S":
		t += "CEST"
	case "W":
		t += "CET"
	default:
		return time.Time{}, ErrInvalidTimestampSeason
	}

	return time.Parse("060102150405MST", t)
}

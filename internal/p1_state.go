package internal

import (
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/common/log"
	"github.com/skoef/gop1"
)

type P1State struct {
	Logger log.Logger
	P1     *gop1.P1

	mutex                       sync.RWMutex
	timestamp                   time.Time
	timestampDifference         time.Duration
	version                     int
	equipmentIdentifier         string
	gasEquipmentIdentifier      string
	electricPowerDelivered      Power
	totalElectricityDelivered   map[int]Energy
	electricPowerInjected       Power
	totalElectricityInjected    map[int]Energy
	electricCurrent             map[string]ElectricCurrent
	voltage                     map[string]Voltage
	electricityTariffIndicator  int
	breakerState                BreakerState
	electricityLimiterThreshold Power
	fuseThreshold               map[string]ElectricCurrent
	totalGasDeliveredTimestamp  time.Time
	totalGasDelivered           Volume
	gasValveState               GasValveState
}

func (s *P1State) Timestamp() time.Time {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.timestamp
}

func (s *P1State) Version() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.version
}

func (s *P1State) EquipmentIndentifier() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.equipmentIdentifier
}

func (s *P1State) GasEquipmentIndentifier() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.gasEquipmentIdentifier
}

func (s *P1State) ElectricPowerDelivered() Power {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.electricPowerDelivered
}

func (s *P1State) TotalElectricityDelivered() map[int]Energy {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.totalElectricityDelivered
}

func (s *P1State) ElectricPowerInjected() Power {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.electricPowerInjected
}

func (s *P1State) TotalElectricityInjected() map[int]Energy {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.totalElectricityInjected
}

func (s *P1State) ElectricCurrent() map[string]ElectricCurrent {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.electricCurrent
}

func (s *P1State) Voltage() map[string]Voltage {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.voltage
}

func (s *P1State) ElectricityTariffIndicator() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.electricityTariffIndicator
}

func (s *P1State) BreakerState() BreakerState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.breakerState
}

func (s *P1State) ElectricityLimiterThreshold() Power {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.electricityLimiterThreshold
}

func (s *P1State) FuseThreshold() map[string]ElectricCurrent {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.fuseThreshold
}

func (s *P1State) TotalGasDeliveredTimestamp() time.Time {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.totalGasDeliveredTimestamp
}

func (s *P1State) TotalGasDelivered() Volume {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.totalGasDelivered
}

func (s *P1State) GasValveState() GasValveState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.gasValveState
}

func (s *P1State) Start() {
	s.P1.Start()
	for {
		if err := s.handleTelegram(<-s.P1.Incoming); err != nil {
			s.Logger.Errorln(err)
		}
	}
}

func (s *P1State) handleTelegram(t *gop1.Telegram) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, o := range t.Objects {
		switch o.Type {
		case gop1.OBISTypeVersionInformation:
			v, err := strconv.Atoi(o.Values[0].Value)
			if err != nil {
				return err
			}

			s.version = v

		case gop1.OBISTypeDateTimestamp:
			v, err := ParseTimestamp(o.Values[0])
			if err != nil {
				return err
			}

			d := time.Since(v)
			if s.timestampDifference > 0 {
				s.timestampDifference = (s.timestampDifference + d) / 2
			} else {
				s.timestampDifference = d
			}

			s.timestamp = v.Add(s.timestampDifference)

		case gop1.OBISTypeEquipmentIdentifier:
			s.equipmentIdentifier = o.Values[0].Value

		case gop1.OBISTypeGasEquipmentIdentifier:
			s.gasEquipmentIdentifier = o.Values[0].Value

		case gop1.OBISTypeElectricityDeliveredTariff1:
			v, err := ParseEnergy(o.Values[0])
			if err != nil {
				return err
			}

			s.totalElectricityDelivered[1] = v

		case gop1.OBISTypeElectricityDeliveredTariff2:
			v, err := ParseEnergy(o.Values[0])
			if err != nil {
				return err
			}

			s.totalElectricityDelivered[2] = v

		case gop1.OBISTypeElectricityGeneratedTariff1:
			v, err := ParseEnergy(o.Values[0])
			if err != nil {
				return err
			}

			s.totalElectricityInjected[1] = v

		case gop1.OBISTypeElectricityGeneratedTariff2:
			v, err := ParseEnergy(o.Values[0])
			if err != nil {
				return err
			}

			s.totalElectricityInjected[2] = v

		case gop1.OBISTypeElectricityTariffIndicator:
			v, err := ParseElectricityTariffIndicator(o.Values[0])
			if err != nil {
				return err
			}

			s.electricityTariffIndicator = v

		case gop1.OBISTypeElectricityDelivered:
			v, err := ParsePower(o.Values[0])
			if err != nil {
				return err
			}

			s.electricPowerDelivered = v

		case gop1.OBISTypeElectricityGenerated:
			v, err := ParsePower(o.Values[0])
			if err != nil {
				return err
			}

			s.electricPowerInjected = v

		case gop1.OBISTypeInstantaneousVoltageL1:
			v, err := ParseVoltage(o.Values[0])
			if err != nil {
				return err
			}

			s.voltage["l1"] = v

		case gop1.OBISTypeInstantaneousVoltageL2:
			v, err := ParseVoltage(o.Values[0])
			if err != nil {
				return err
			}

			s.voltage["l2"] = v

		case gop1.OBISTypeInstantaneousVoltageL3:
			v, err := ParseVoltage(o.Values[0])
			if err != nil {
				return err
			}

			s.voltage["l3"] = v

		case gop1.OBISTypeInstantaneousCurrentL1:
			v, err := ParseElectricCurrent(o.Values[0])
			if err != nil {
				return err
			}

			s.electricCurrent["l1"] = v

		case gop1.OBISTypeInstantaneousCurrentL2:
			v, err := ParseElectricCurrent(o.Values[0])
			if err != nil {
				return err
			}

			s.electricCurrent["l2"] = v

		case gop1.OBISTypeInstantaneousCurrentL3:
			v, err := ParseElectricCurrent(o.Values[0])
			if err != nil {
				return err
			}

			s.electricCurrent["l3"] = v

		case gop1.OBISTypeGasDelivered:
			{
				v, err := ParseTimestamp(o.Values[0])
				if err != nil {
					return err
				}

				s.totalGasDeliveredTimestamp = v
			}

			{
				v, err := ParseVolume(o.Values[1])
				if err != nil {
					return err
				}

				s.totalGasDelivered = v
			}

		case gop1.OBISTypeBreakerState:
			v, err := ParseBreakerState(o.Values[0])
			if err != nil {
				return err
			}

			s.breakerState = v

		case gop1.OBISTypeLimiterThreshold:
			v, err := ParsePower(o.Values[0])
			if err != nil {
				return err
			}

			s.electricityLimiterThreshold = v

		case gop1.OBISTypeFuseThresholdL1:
			v, err := ParseElectricCurrent(o.Values[0])
			if err != nil {
				return err
			}

			s.fuseThreshold["l1"] = v

		case gop1.OBISTypeGasValveState:
			v, err := ParseGasValveState(o.Values[0])
			if err != nil {
				return err
			}

			s.gasValveState = v
		}
	}

	return nil
}

func NewP1State(l log.Logger, p1 *gop1.P1) *P1State {
	return &P1State{
		Logger: l,
		P1:     p1,

		totalElectricityDelivered: make(map[int]Energy),
		totalElectricityInjected:  make(map[int]Energy),
		electricCurrent:           make(map[string]ElectricCurrent),
		voltage:                   make(map[string]Voltage),
		fuseThreshold:             make(map[string]ElectricCurrent),
	}
}

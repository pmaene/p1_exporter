package internal

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "p1"
	upTimeout = time.Minute
)

var (
	upDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Whether collecting smart meter metrics was successful.",
		nil,
		nil,
	)

	versionDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "version"),
		"Version information for the P1 output.",
		nil,
		nil,
	)

	electricPowerDeliveredDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "electricity", "power_delivered"),
		"Electricity being delivered to the premises.",
		[]string{"equipment_id"},
		nil,
	)

	totalElectricityDeliveredDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "electricity", "delivered_total"),
		"Total electricity delivered to the premises.",
		[]string{"equipment_id", "tariff"},
		nil,
	)

	electricPowerInjectedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "electricity", "power_injected"),
		"Electricity being injected by the premises.",
		[]string{"equipment_id"},
		nil,
	)

	totalElectricityInjectedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "electricity", "injected_total"),
		"Total electricity injected by the premises.",
		[]string{"equipment_id", "tariff"},
		nil,
	)

	electricCurrentDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "electricity", "current"),
		"Instantaneous current measured by the smart meter.",
		[]string{"equipment_id", "phase"},
		nil,
	)

	voltageDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "electricity", "voltage"),
		"Instantaneous voltage measured by the smart meter.",
		[]string{"equipment_id", "phase"},
		nil,
	)

	electricityTariffIndicatorDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "electricity", "tariff_indicator"),
		"Electricity tariff that is currently active.",
		[]string{"equipment_id"},
		nil,
	)

	breakerStateDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "electricity", "breaker_state"),
		"State of the smart meter's breaker.",
		[]string{"equipment_id"},
		nil,
	)

	electricityLimiterThresholdDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "electricity", "limiter_threshold"),
		"Threshold for the electricity limiter.",
		[]string{"equipment_id"},
		nil,
	)

	fuseThresholdDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "electricity", "fuse_threshold"),
		"Threshold for the smart meter's fuse.",
		[]string{"equipment_id", "phase"},
		nil,
	)

	totalGasDeliveredDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "gas", "delivered_total"),
		"Total gas volume delivered to the premises.",
		[]string{"equipment_id"},
		nil,
	)

	gasValveStateDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "gas", "valve_state"),
		"State of the gas valve.",
		[]string{"equipment_id"},
		nil,
	)
)

type Collector struct {
	P1State *P1State
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upDesc
	ch <- versionDesc
	ch <- electricPowerDeliveredDesc
	ch <- totalElectricityDeliveredDesc
	ch <- electricPowerInjectedDesc
	ch <- totalElectricityInjectedDesc
	ch <- electricCurrentDesc
	ch <- voltageDesc
	ch <- electricityTariffIndicatorDesc
	ch <- breakerStateDesc
	ch <- electricityLimiterThresholdDesc
	ch <- fuseThresholdDesc
	ch <- totalGasDeliveredDesc
	ch <- gasValveStateDesc
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.NewMetricWithTimestamp(
		c.P1State.Timestamp(),
		prometheus.MustNewConstMetric(
			upDesc,
			prometheus.GaugeValue,
			c.up(),
		),
	)

	ch <- prometheus.NewMetricWithTimestamp(
		c.P1State.Timestamp(),
		prometheus.MustNewConstMetric(
			versionDesc,
			prometheus.CounterValue,
			float64(c.P1State.Version()),
		),
	)

	ch <- prometheus.NewMetricWithTimestamp(
		c.P1State.Timestamp(),
		prometheus.MustNewConstMetric(
			electricPowerDeliveredDesc,
			prometheus.GaugeValue,
			float64(c.P1State.ElectricPowerDelivered()),
			c.P1State.EquipmentIndentifier(),
		),
	)

	for k, v := range c.P1State.TotalElectricityDelivered() {
		ch <- prometheus.NewMetricWithTimestamp(
			c.P1State.Timestamp(),
			prometheus.MustNewConstMetric(
				totalElectricityDeliveredDesc,
				prometheus.CounterValue,
				float64(v),
				c.P1State.EquipmentIndentifier(),
				strconv.Itoa(k),
			),
		)
	}

	ch <- prometheus.NewMetricWithTimestamp(
		c.P1State.Timestamp(),
		prometheus.MustNewConstMetric(
			electricPowerInjectedDesc,
			prometheus.GaugeValue,
			float64(c.P1State.ElectricPowerInjected()),
			c.P1State.EquipmentIndentifier(),
		),
	)

	for k, v := range c.P1State.TotalElectricityInjected() {
		ch <- prometheus.NewMetricWithTimestamp(
			c.P1State.Timestamp(),
			prometheus.MustNewConstMetric(
				totalElectricityInjectedDesc,
				prometheus.CounterValue,
				float64(v),
				c.P1State.EquipmentIndentifier(),
				strconv.Itoa(k),
			),
		)
	}

	for k, v := range c.P1State.ElectricCurrent() {
		ch <- prometheus.NewMetricWithTimestamp(
			c.P1State.Timestamp(),
			prometheus.MustNewConstMetric(
				electricCurrentDesc,
				prometheus.GaugeValue,
				float64(v),
				c.P1State.EquipmentIndentifier(),
				k,
			),
		)
	}

	for k, v := range c.P1State.Voltage() {
		ch <- prometheus.NewMetricWithTimestamp(
			c.P1State.Timestamp(),
			prometheus.MustNewConstMetric(
				voltageDesc,
				prometheus.GaugeValue,
				float64(v),
				c.P1State.EquipmentIndentifier(),
				k,
			),
		)
	}

	ch <- prometheus.NewMetricWithTimestamp(
		c.P1State.Timestamp(),
		prometheus.MustNewConstMetric(
			electricityTariffIndicatorDesc,
			prometheus.GaugeValue,
			float64(c.P1State.ElectricityTariffIndicator()),
			c.P1State.EquipmentIndentifier(),
		),
	)

	ch <- prometheus.NewMetricWithTimestamp(
		c.P1State.Timestamp(),
		prometheus.MustNewConstMetric(
			breakerStateDesc,
			prometheus.GaugeValue,
			float64(c.P1State.BreakerState()),
			c.P1State.EquipmentIndentifier(),
		),
	)

	ch <- prometheus.NewMetricWithTimestamp(
		c.P1State.Timestamp(),
		prometheus.MustNewConstMetric(
			electricityLimiterThresholdDesc,
			prometheus.GaugeValue,
			float64(c.P1State.ElectricityLimiterThreshold()),
			c.P1State.EquipmentIndentifier(),
		),
	)

	for k, v := range c.P1State.FuseThreshold() {
		ch <- prometheus.NewMetricWithTimestamp(
			c.P1State.Timestamp(),
			prometheus.MustNewConstMetric(
				fuseThresholdDesc,
				prometheus.GaugeValue,
				float64(v),
				c.P1State.EquipmentIndentifier(),
				k,
			),
		)
	}

	ch <- prometheus.NewMetricWithTimestamp(
		c.P1State.TotalGasDeliveredTimestamp(),
		prometheus.MustNewConstMetric(
			totalGasDeliveredDesc,
			prometheus.CounterValue,
			float64(c.P1State.TotalGasDelivered()),
			c.P1State.GasEquipmentIndentifier(),
		),
	)

	ch <- prometheus.NewMetricWithTimestamp(
		c.P1State.Timestamp(),
		prometheus.MustNewConstMetric(
			gasValveStateDesc,
			prometheus.GaugeValue,
			float64(c.P1State.GasValveState()),
			c.P1State.GasEquipmentIndentifier(),
		),
	)
}

func (c *Collector) up() float64 {
	if time.Since(c.P1State.Timestamp()) > upTimeout {
		return 0
	}

	return 1
}

func NewCollector(s *P1State) *Collector {
	return &Collector{
		P1State: s,
	}
}

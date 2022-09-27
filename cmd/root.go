package cmd

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pmaene/p1_exporter/internal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/skoef/gop1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version string
	commit  string
)

var (
	listenAddress     string
	metricsPath       string
	readHeaderTimeout time.Duration
	p1USBDevice       string
	p1Baudrate        int
	p1Timeout         int

	rootCmd = &cobra.Command{
		Use:          "p1_exporter",
		Short:        "P1 Exporter",
		SilenceUsage: true,
		Run:          runRoot,
	}
)

func GetVersion() string {
	if version == "" {
		return "(devel)"
	}

	return version
}

func SetVersion(v string) {
	version = v
}

func GetCommit() string {
	return commit
}

func SetCommit(c string) {
	commit = c
}

func Execute() {
	rootCmd.Version = fmt.Sprintf(
		"%s (%s)",
		GetVersion(),
		GetCommit(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(
		&listenAddress,
		"web.listen-address",
		":9786",
		"address on which to expose metrics and web interface",
	)

	rootCmd.Flags().StringVar(
		&metricsPath,
		"web.telemetry-path",
		"/metrics",
		"path under which to expose metrics",
	)

	rootCmd.Flags().DurationVar(
		&readHeaderTimeout,
		"web.read-header-timeout",
		5*time.Second,
		"timeout for reading request headers",
	)

	rootCmd.Flags().StringVar(
		&p1USBDevice,
		"p1.usb-device",
		"/dev/ttyUSB0",
		"path to the smart meter's serial device",
	)

	rootCmd.Flags().IntVar(
		&p1Baudrate,
		"p1.baudrate",
		115200,
		"baud rate of the smart meter's serial connection",
	)

	rootCmd.Flags().IntVar(
		&p1Timeout,
		"p1.timeout",
		500,
		"smart meter read timeout in milliseconds",
	)

	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("P1_EXPORTER")
	viper.SetEnvKeyReplacer(
		strings.NewReplacer(".", "_", "-", "_"),
	)
}

func runRoot(cmd *cobra.Command, args []string) {
	log.Infoln("starting", cmd.Name(), cmd.Version)

	p1, err := gop1.New(
		gop1.P1Config{
			USBDevice: viper.GetString("p1.usb-device"),
			Baudrate:  viper.GetInt("p1.baudrate"),
			Timeout:   viper.GetInt("p1.timeout"),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	s := internal.NewP1State(log.Base(), p1)
	go s.Start()

	c := internal.NewCollector(s)
	if err := prometheus.Register(c); err != nil {
		log.Fatal(err)
	}

	http.Handle(
		viper.GetString("web.telemetry-path"),
		promhttp.Handler(),
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(
			[]byte(
				`<html>
				<head><title>P1 Exporter</title></head>
				<body>
				<h1>P1 Exporter</h1>
				<p><a href='/metrics'>Metrics</a></p>
				</body>
				</html>`,
			),
		)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	srv := http.Server{
		ReadHeaderTimeout: viper.GetDuration("web.read-header-timeout"),
	}

	lst, err := net.Listen(
		"tcp",
		viper.GetString("web.listen-address"),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Infoln("listening on", viper.GetString("web.listen-address"))
	if err := srv.Serve(lst); err != nil {
		log.Fatal(err)
	}
}

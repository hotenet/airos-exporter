package main

import (
	log "log/slog"
	"net/http"
	"os"
	"strings"
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

var (
	configFile = kingpin.Flag("config", "Configuration file path").Required().File()
)


func parseLevel(lvl string) log.Level {
	switch (strings.ToUpper(lvl)) {
	case log.LevelDebug.String():
		return log.LevelDebug
	case log.LevelInfo.String():
		return log.LevelInfo
	case log.LevelError.String():
		return log.LevelError
	case log.LevelWarn.String():
		return log.LevelError
	default:
		return log.LevelError
	}
}

func main() {
	kingpin.Version(version.Print("s3rw"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := log.New(log.NewJSONHandler(os.Stdout, nil))
	log.SetDefault(logger)
	config := NewConfig(*configFile)
	logger = log.New(log.NewJSONHandler(os.Stdout, &log.HandlerOptions{Level: parseLevel(config.Log.Level)}))
	log.SetDefault(logger)

	handler := NewHandler(config.Ubnt)
	mux := http.NewServeMux()
	mux.Handle(config.Exporter.Path, http.HandlerFunc(func(w http.ResponseWriter,  req *http.Request) {
		target := NewTarget(req)
		set := NewMetricSet(config.Exporter.Namespace)
		reg := set.Registry()
		handler.Update(set, target)
		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, req)
	}))

	log.Info("listening http", "listen", config.Exporter.Listen)
	if err := http.ListenAndServe(config.Exporter.Listen, mux); err != nil {
		log.Error("unable to start exporter server", "error", err)
		os.Exit(1)
	}
}

type target struct {
	Scheme string
	Host string
	Port string
	InsecureSkipVerify bool
}

func NewTarget(req *http.Request) target {
	target := target{
		Scheme: "https",
		InsecureSkipVerify: false,
	}
	values := req.URL.Query()
	address := values.Get("address")
	parts := strings.Split(address, ":")
	host := parts[0]
	port := "443"

	if len(parts) > 1 {
		port = parts[1]
	}
	target.Port = port
	target.Host = host
	if values.Has("scheme") {
		target.Scheme = values.Get("scheme")
	}
	if values.Has("skip-ssl") {
		target.InsecureSkipVerify = true
	}
	return target
}

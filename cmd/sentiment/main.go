package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"github.com/jukeizu/sentiment/pkg/treediagram"
	"github.com/machinebox/sdk-go/textbox"
	"github.com/oklog/run"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

var Version = ""

var (
	flagVersion = false

	httpPort           = "10002"
	flagTextBoxAddress = "http://localhost:8080"
)

func parseConfig() {
	flag.StringVar(&httpPort, "http.port", httpPort, "http port for handler")
	flag.StringVar(&flagTextBoxAddress, "textbox.addr", flagTextBoxAddress, "textbox address")
	flag.BoolVar(&flagVersion, "v", false, "version")

	flag.Parse()
}

func main() {
	parseConfig()

	if flagVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().
		Str("instance", xid.New().String()).
		Str("component", "sentiment").
		Str("version", Version).
		Logger()

	httpAddr := ":" + httpPort

	textboxClient := textbox.New(flagTextBoxAddress)
	handler := treediagram.NewHandler(logger, httpAddr, textboxClient)

	g := run.Group{}

	g.Add(func() error {
		return handler.Start()
	}, func(error) {
		err := handler.Stop()
		if err != nil {
			logger.Error().Err(err).Caller().Msg("couldn't stop handler")
		}
	})

	cancel := make(chan struct{})
	g.Add(func() error {
		return interrupt(cancel)
	}, func(error) {
		close(cancel)
	})

	logger.Info().Err(g.Run()).Msg("stopped")
}

func interrupt(cancel <-chan struct{}) error {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-cancel:
		return errors.New("stopping")
	case sig := <-c:
		return fmt.Errorf("%s", sig)
	}
}

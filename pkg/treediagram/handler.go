package treediagram

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jukeizu/contract"
	"github.com/machinebox/sdk-go/textbox"
	"github.com/rs/zerolog"
)

const AppId = "intent.endpoint.sentiment"

type Handler struct {
	logger        zerolog.Logger
	httpServer    *http.Server
	textboxClient *textbox.Client
}

func NewHandler(logger zerolog.Logger, addr string, textboxClient *textbox.Client) Handler {
	logger = logger.With().Str("component", AppId).Logger()

	httpServer := http.Server{
		Addr: addr,
	}

	return Handler{logger, &httpServer, textboxClient}
}

func (h Handler) Sentiment(request contract.Request) (*contract.Response, error) {
	analysis, err := h.textboxClient.Check(strings.NewReader(request.Content))
	if err != nil {
		return nil, fmt.Errorf("machinebox: %s", err.Error())
	}

	reaction := FormatSentimentReaction(request, analysis)
	if reaction == nil {
		return nil, nil
	}

	return &contract.Response{Reactions: []*contract.Reaction{reaction}}, nil
}

func (h Handler) Start() error {
	h.logger.Info().Msg("starting")

	mux := http.NewServeMux()
	mux.HandleFunc("/sentiment", h.makeLoggingHttpHandlerFunc("sentiment", h.Sentiment))

	h.httpServer.Handler = mux

	return h.httpServer.ListenAndServe()
}

func (h Handler) Stop() error {
	h.logger.Info().Msg("stopping")

	return h.httpServer.Shutdown(context.Background())
}

func (h Handler) makeLoggingHttpHandlerFunc(name string, f func(contract.Request) (*contract.Response, error)) http.HandlerFunc {
	contractHandlerFunc := contract.MakeRequestHttpHandlerFunc(f)

	return func(w http.ResponseWriter, r *http.Request) {
		defer func(begin time.Time) {
			h.logger.Info().
				Str("intent", name).
				Str("took", time.Since(begin).String()).
				Msg("called")
		}(time.Now())

		contractHandlerFunc.ServeHTTP(w, r)
	}
}

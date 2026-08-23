package feeds

import (
	"net/http"

	"github.com/aweigor/nts-feed-generator/config"
)

type FeedsHandler struct{}

func NewFeedsHandler(router *http.ServeMux, conf *config.Config) {
	handler := &FeedsHandler{}
	router.HandleFunc("/feeds/latest", handler.HandleLatest())
	router.HandleFunc("/feeds/show/:showId", handler.HandleShow())
}

func (handler *FeedsHandler) HandleLatest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (handler *FeedsHandler) HandleShow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

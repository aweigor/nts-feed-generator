package feeds

import (
	"net/http"

	"github.com/aweigor/nts-feed-generator/config"
	"github.com/aweigor/nts-feed-generator/pkg/res"
)

type FeedsHandler struct{}

func NewFeedsHandler(router *http.ServeMux, conf *config.Config) {
	handler := &FeedsHandler{}
	router.HandleFunc("/feeds/latest", handler.HandleLatest())
	router.HandleFunc("/feeds/show/:showId", handler.HandleShow())
}

func (handler *FeedsHandler) HandleLatest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParser := NewFeedsQueryParser()
		params, err := queryParser.parseQuery(r.URL.Query())
		if err != nil {
			res.Error(w, http.StatusBadRequest, err.Error())
		}
		switch format := params.format; format {
		case "xml":
			w.Header().Set("Content-Type", "application/xml")
		case "json":
			w.Header().Set("Content-Type", "application/json")
		}
	}
}

func (handler *FeedsHandler) HandleShow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

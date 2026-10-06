package feeds

import (
	"net/http"

	"github.com/aweigor/nts-feed-generator/pkg/ntsclient"
	"github.com/aweigor/nts-feed-generator/pkg/res"
)

type FeedsHandler struct {
	*ntsclient.NTSClient
}

type FeedsHandlerDeps struct {
	*ntsclient.NTSClient
}

func NewFeedsHandler(router *http.ServeMux, deps FeedsHandlerDeps) {
	handler := &FeedsHandler{
		NTSClient: deps.NTSClient,
	}
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

		searchParams := ntsclient.EpisodesSearchParams{
			Page: ntsclient.PageQuery{
				Offset: params.offset,
				Limit:  params.limit,
			},
		}

		res, err := handler.NTSClient.FetchLatest(r.Context(), searchParams)

		channelItems := make([]Item, params.limit)

		for i := range params.limit {
			sourceData := res.Results[i]

			channelItems[i] = Item{
				Title:       sourceData.Title,
				Link:        sourceData.AudioSources[0].URL,
				Description: sourceData.Description,
				PubDate:     sourceData.LocalDate,
				Enclosure:   "",
				Duration:    "",
				GUID:        "",
				AirDate:     sourceData.LocalDate,
				Duration:    "",
				Tracklist:   []Track{},
			}
		}

		data := RSSResponse{
			Channel{
				Title: "Latest episodes",
				Items: []Item{},
			},
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

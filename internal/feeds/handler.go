package feeds

import (
	"net/http"

	"github.com/aweigor/nts-feed-generator/config"
	"github.com/aweigor/nts-feed-generator/pkg/ntsclient"
	"github.com/aweigor/nts-feed-generator/pkg/res"
)

type FeedsHandler struct {
	*ntsclient.NTSClient
	*config.FeedsConfig
}

type FeedsHandlerDeps struct {
	*ntsclient.NTSClient
	*config.FeedsConfig
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

		episodes, err := handler.NTSClient.FetchLatest(r.Context(), searchParams)

		channelItems := make([]Item, params.limit)

		for i := range len(episodes.Results) {
			sourceData := episodes.Results[i]

			channelItems[i] = Item{
				Title:       sourceData.Title,
				Link:        handler.buildEpisodeLink(sourceData.Article.Path),
				Description: "",
				PubDate:     sourceData.LocalDate,
				Enclosure: Enclosure{
					URL:    sourceData.AudioSources[0].URL,
					Type:   sourceData.AudioSources[0].Source,
					Length: nil,
				},
				Duration:  nil,
				GUID:      "",
				AirDate:   sourceData.LocalDate,
				Tracklist: []Track{},
			}
		}

		data := RSSResponse{
			Channel{
				Title: "Latest episodes",
				Items: channelItems,
			},
		}

		switch format := params.format; format {
		case "xml":
			res.Xml(w, data, 200)
		case "json":
			res.Json(w, data, 200)
		}
	}
}

func (handler *FeedsHandler) HandleShow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (handler *FeedsHandler) buildEpisodeLink(articlePath string) string {
	return handler.Nts.APIV2Url + articlePath
}

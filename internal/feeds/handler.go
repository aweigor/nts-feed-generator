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
		NTSClient:   deps.NTSClient,
		FeedsConfig: deps.FeedsConfig,
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

		channelItems := make([]EpisodeItem, len(episodes.Results))

		for i := range len(episodes.Results) {
			sourceData := episodes.Results[i]

			channelItems[i] = handler.buildEpisodeItem(sourceData)
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

func (handler *FeedsHandler) buildEpisodeItem(episodeInfo ntsclient.EpisodeInfo) EpisodeItem {
	return EpisodeItem{
		Title:       episodeInfo.Title,
		Link:        handler.buildEpisodeLink(episodeInfo.Article.Path),
		Description: "",
		PubDate:     episodeInfo.LocalDate,
		Enclosure: Enclosure{
			URL:    episodeInfo.AudioSources[0].URL,
			Type:   episodeInfo.AudioSources[0].Source,
			Length: nil,
		},
		Duration:  nil,
		GUID:      "",
		AirDate:   episodeInfo.LocalDate,
		Tracklist: []Track{},
	}
}

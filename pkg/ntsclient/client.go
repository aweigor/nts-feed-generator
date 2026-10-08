package ntsclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aweigor/nts-feed-generator/config"
)

type NTSClient struct {
	baseUrl    string
	httpClient *http.Client
}

func NewNTSClient(cfg *config.NtsAPIProperties) (*NTSClient, error) {
	if cfg.APIV2Url == "" {
		return nil, errors.New("ntsclient: api base required")
	}
	return &NTSClient{
		baseUrl:    cfg.APIV2Url,
		httpClient: &http.Client{},
	}, nil
}

func (client *NTSClient) FetchLatest(ctx context.Context, params EpisodesSearchParams) (*ArticlesResponse, error) {
	query := buildEpisodesSearchQuery(params)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.buildLatestEpisodesURL(query), nil)
	if err != nil {
		return nil, fmt.Errorf("ntsclient[FetchLatest]: build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ntsclient[FetchLatest]: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ntsclient[FetchLatest]: bad status %d", resp.StatusCode)
	}

	var out ArticlesResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("ntsclient[FetchLatest]: unmarschall response error: %w", err)
	}

	return &out, err
}

func (client *NTSClient) FetchShowEpisodes(ctx context.Context, showId string, params EpisodesSearchParams) (*EpisodesResponse, error) {
	query := buildEpisodesSearchQuery(params)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, client.buildShowEpisodesURL(showId, query), nil)

	request.Header.Set("Accept", "application/json")

	response, err := client.httpClient.Do(request)

	defer response.Body.Close()

	var out EpisodesResponse
	if err := json.NewDecoder(response.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("ntsclient[FetchShowEpisodes]: unmarschall response error", err)
	}

	return &out, err
}

func (client *NTSClient) buildShowEpisodesURL(showId string, urlQuery string) string {
	urlPath := strings.ReplaceAll(GetShowEpisodesPath, "{show_id}", showId)
	return fmt.Sprintf("%s%s?%s", client.baseUrl, urlPath, urlQuery)
}

func (client *NTSClient) buildLatestEpisodesURL(urlQuery string) string {
	return fmt.Sprintf("%s%s?%s", client.baseUrl, GetLatestPath, urlQuery)
}

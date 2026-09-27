package ntsclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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

func (client *NTSClient) FetchLatest(ctx context.Context, params EpisodesSearchParams) (*EpisodesResponse, error) {
	queryString := "?" + buildEpisodesSearchQuery(params)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.baseUrl+GetEpisodesPath+queryString, nil)
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

	var out EpisodesResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("ntsclient[FetchLatest]: unmarschall response error: %w", err)
	}

	return &out, err
}

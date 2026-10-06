package feeds

import (
	"fmt"
	"net/url"
	"strconv"
)

const (
	DEFAULT_FORMAT = "xml"
	DEFAULT_OFFSET = 0
	DEFAULT_LIMIT  = 20
)

type FeedsQueryParser struct{}

func NewFeedsQueryParser() *FeedsQueryParser {
	return &FeedsQueryParser{}
}

func (parser *FeedsQueryParser) parseQuery(q url.Values) (*FeedsQueryParams, error) {
	format := q.Get("format")

	if format == "" {
		format = DEFAULT_FORMAT
	}

	allowedFormats := map[string]bool{
		"json": true,
		"xml":  true,
		"atom": true,
	}

	offset, err := strconv.Atoi(q.Get("offset"))
	if err != nil {
		offset = DEFAULT_OFFSET
	}

	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil {
		limit = DEFAULT_LIMIT
	}

	if !allowedFormats[format] {
		return nil, fmt.Errorf("[FeedsQueryParser]: format must be 'json' or 'xml' or 'atom', got '%s'", format)
	}

	return &FeedsQueryParams{
		format,
		offset,
		limit,
	}, nil
}

package feeds

import (
	"fmt"
	"net/url"
)

const DEFAULT_FORMAT = "xml"

type FeedsQuery struct {
	format string
}

type FeedsQueryParser struct{}

func NewFeedsQueryParser() *FeedsQueryParser {
	return &FeedsQueryParser{}
}

func (parser *FeedsQueryParser) parseQuery(q url.Values) (*FeedsQuery, error) {
	format := q.Get("format")

	if format == "" {
		format = DEFAULT_FORMAT
	}

	allowedFormats := map[string]bool{
		"json": true,
		"xml":  true,
		"atom": true,
	}

	if !allowedFormats[format] {
		return nil, fmt.Errorf("[FeedsQueryParser]: format must be 'json' or 'xml' or 'atom', got '%s'", format)
	}

	return &FeedsQuery{
		format,
	}, nil
}

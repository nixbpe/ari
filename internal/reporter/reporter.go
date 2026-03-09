package reporter

import (
	"context"

	"github.com/bbik/ari/internal/scorer"
)

// Format represents the output format
type Format int

const (
	FormatJSON Format = iota
	FormatHTML
	FormatText
)

func (f Format) String() string {
	switch f {
	case FormatJSON:
		return "json"
	case FormatHTML:
		return "html"
	case FormatText:
		return "text"
	default:
		return "unknown"
	}
}

// Reporter generates reports from evaluation scores
type Reporter interface {
	Report(ctx context.Context, score *scorer.Score) (string, error)
	Format() Format
}

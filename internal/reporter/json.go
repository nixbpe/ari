package reporter

import (
	"context"
	"encoding/json"
	"io"
)

type JSONReporter struct{}

func (j *JSONReporter) Report(ctx context.Context, report *Report, w io.Writer) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

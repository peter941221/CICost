package billing

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func LoadActualFromCSV(path, repo, period string) (float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true
	line := 0
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}
		line++
		if len(rec) < 3 {
			continue
		}

		repoVal := strings.TrimSpace(rec[0])
		periodVal := strings.TrimSpace(rec[1])
		costVal := strings.TrimSpace(rec[2])
		if line == 1 && strings.EqualFold(repoVal, "repo") {
			continue
		}
		if repoVal != repo || periodVal != period {
			continue
		}
		fv, err := strconv.ParseFloat(costVal, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid cost value at line %d: %w", line, err)
		}
		return fv, nil
	}
	return 0, fmt.Errorf("no billing row matched repo=%s period=%s in %s", repo, period, path)
}

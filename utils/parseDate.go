package utils

import (
	"time"

	"github.com/araddon/dateparse"
)

func ParseDate(dateStr string) (time.Time, error) {
    t, err := dateparse.ParseAny(dateStr)
    if err != nil {
        return time.Time{}, err
    }
    return t, nil
}
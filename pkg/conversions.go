package pkg

// Port of src/clu/conversions.py providing byte size and time formatting
// plus parsing of simple SI size strings.

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// BytesToSI converts a byte count to a human-readable SI string (base 1024) with one decimal.
// Example: 1536 -> "1.5 KB"
func BytesToSI(size float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	val := size
	for _, u := range units {
		if val < 1024.0 || u == units[len(units)-1] {
			return fmt.Sprintf("%.1f %s", val, u)
		}
		val /= 1024.0
	}
	return fmt.Sprintf("%.1f %s", val, units[len(units)-1])
}

// SIToBytes parses a human-readable size string (e.g. "1.5 KB") into bytes (float64).
// Accepts units B,K,KB,M,MB,G,GB,T,TB,P,PB,E,EB (case-insensitive). Commas permitted in number.
func SIToBytes(sizeStr string) (float64, error) {
	units := map[string]int{
		"B": 0,
		"K": 1, "KB": 1,
		"M": 2, "MB": 2,
		"G": 3, "GB": 3,
		"T": 4, "TB": 4,
		"P": 5, "PB": 5,
		"E": 6, "EB": 6,
	}
	s := strings.TrimSpace(strings.ToUpper(sizeStr))
	var numberStr, unitStr string
	for _, ch := range s {
		if (ch >= '0' && ch <= '9') || ch == '.' || ch == ',' {
			numberStr += string(ch)
		} else if ch >= 'A' && ch <= 'Z' {
			unitStr += string(ch)
		}
	}
	if numberStr == "" {
		return 0, errors.New("no numeric portion found")
	}
	num, err := parseFloat(numberStr)
	if err != nil {
		return 0, err
	}
	unitStr = strings.TrimSpace(unitStr)
	exp, ok := units[unitStr]
	if !ok {
		return 0, fmt.Errorf("unknown unit: %s", unitStr)
	}
	return num * math.Pow(1024, float64(exp)), nil
}

// parseFloat handles comma removal then parsing.
func parseFloat(s string) (float64, error) {
	s = strings.ReplaceAll(s, ",", "")
	var num float64
	_, err := fmt.Sscan(s, &num)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", s)
	}
	return num, nil
}

// SecondsToText converts seconds into a human-readable string with months, days, hours, minutes, seconds.
// Follows Python logic (using 30-day months) and pluralization rules.
func SecondsToText(secs int64) string {
	if secs <= 0 {
		return "0 seconds"
	}
	months := secs / 2592000 // 30*24*60*60
	remaining := secs % 2592000

	days := remaining / 86400
	remaining = remaining % 86400

	hours := remaining / 3600
	remaining = remaining % 3600

	minutes := remaining / 60
	seconds := remaining % 60

	parts := make([]string, 0, 5)
	if months > 0 {
		parts = append(parts, pluralize(months, "month"))
	}
	if days > 0 {
		parts = append(parts, pluralize(days, "day"))
	}
	if hours > 0 {
		parts = append(parts, pluralize(hours, "hour"))
	}
	if minutes > 0 {
		parts = append(parts, pluralize(minutes, "minute"))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, pluralize(seconds, "second"))
	}
	return strings.Join(parts, ", ")
}

func pluralize(n int64, word string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, word)
	}
	return fmt.Sprintf("%d %ss", n, word)
}

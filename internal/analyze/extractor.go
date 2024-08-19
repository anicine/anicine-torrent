package analyze

import (
	"math"
	"strconv"
	"strings"
)

func ExtractNum(i1 string) int {
	if i1 == "" {
		return 0
	}

	match := numExp.FindStringSubmatch(i1)
	if len(match) < 2 {
		return 0
	}

	number, err := strconv.Atoi(match[1])
	if err != nil {
		return 0
	}

	return number
}

func ExtractNyaaID(i1 string) int64 {
	if i1 == "" {
		return 0
	}

	match := nyaaExp.FindStringSubmatch(i1)
	if len(match) > 1 {
		id, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return 0
		}

		return id
	}

	return 0
}

func ExtractMiB(i1 string) (float64, error) {
	size, err := strconv.ParseFloat(i1, 64)
	if err != nil {
		return 0, err
	}

	convertedSize := size * 1024

	return convertedSize, nil
}

func ExtractMagnetHash(i1 string) string {
	match := hashExp.FindStringSubmatch(i1)

	if len(match) > 1 {
		return match[1]
	}

	return ""
}

func ExtractBrackets(i1 string) []string {

	matches := bracketsExp.FindAllString(i1, -1)

	return matches
}

func ExtractParentheses(i1 string) []string {

	matches := parenthesesExp.FindAllString(i1, -1)

	return matches
}

func ExtractLanguage(i1 string) string {
	if i1 == "" {
		return ""
	}

	i1 = strings.ToLower(i1)
	for k1, v1 := range languages {
		if strings.Contains(i1, k1) {
			return v1
		}
	}

	return ""
}

func ExtractQuality(i1 string) string {
	if i1 == "" {
		return ""
	}

	i1 = strings.ToLower(i1)
	for k1, v1 := range quality {
		if strings.Contains(i1, k1) {
			return v1
		}
	}

	return ""
}

func ExtractACodec(i1 string) string {
	if i1 == "" {
		return ""
	}

	i1 = strings.ToLower(i1)
	for k1, v1 := range audioCodecs {
		if strings.Contains(i1, k1) {
			return v1
		}
	}

	return ""
}

func ExtractVCodec(i1 string) string {
	if i1 == "" {
		return ""
	}

	i1 = strings.ToLower(i1)
	for k1, v1 := range videoCodecs {
		if strings.Contains(i1, k1) {
			return v1
		}
	}

	return ""
}

func ExtractFloatsWithRanges(i1 string) []float64 {
	var extracted []float64
	encountered := make(map[float64]bool)

	for {
		match := intsExp.FindStringSubmatch(i1)
		if match == nil {
			break
		}

		if match[1] != "" {
			start, _ := strconv.ParseFloat(match[1], 64)
			end, _ := strconv.ParseFloat(match[2], 64)

			for i := start; i <= end+0.99; i++ {
				min := math.Min(i, end)
				if !encountered[min] {
					extracted = append(extracted, min)
					encountered[min] = true
				}
			}
		} else {
			num, _ := strconv.ParseFloat(match[3], 64)

			if !encountered[num] {
				extracted = append(extracted, num)
				encountered[num] = true
			}
		}

		i1 = i1[len(match[0]):]
	}

	return extracted
}

func ExtractSE(i1 string) (int, float64) {
	match := seExp.FindStringSubmatch(i1)
	if match == nil {
		return 0, 0
	}

	season, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, 0
	}

	episode, err := strconv.ParseFloat(match[2], 64)
	if err != nil {
		return 0, 0
	}

	return season, episode
}

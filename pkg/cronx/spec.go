package cronx

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// Parse 解析 crontab 表达式
func Parse(spec string) (*SpecSchedule, error) {
	if len(spec) == 0 {
		return nil, fmt.Errorf("empty spec")
	}

	// 处理特殊描述符
	if strings.HasPrefix(spec, "@") {
		return parseDescriptor(spec)
	}

	fields := strings.Fields(spec)
	if len(fields) != 6 {
		return nil, fmt.Errorf("expected 6 fields, got %d", len(fields))
	}

	second, err := getField(fields[0], seconds)
	if err != nil {
		return nil, err
	}
	minute, err := getField(fields[1], minutes)
	if err != nil {
		return nil, err
	}
	hour, err := getField(fields[2], hours)
	if err != nil {
		return nil, err
	}
	dom, err := getField(fields[3], dom)
	if err != nil {
		return nil, err
	}
	month, err := getField(fields[4], months)
	if err != nil {
		return nil, err
	}
	dow, err := getField(fields[5], dow)
	if err != nil {
		return nil, err
	}

	return &SpecSchedule{
		Second:   second,
		Minute:   minute,
		Hour:     hour,
		Dom:      dom,
		Month:    month,
		Dow:      dow,
		Location: time.Local,
	}, nil
}

func parseDescriptor(spec string) (*SpecSchedule, error) {
	switch spec {
	case "@yearly", "@annually":
		return &SpecSchedule{
			Second:   1,
			Minute:   1,
			Hour:     1,
			Dom:      1,
			Month:    1,
			Dow:      all(dow),
			Location: time.Local,
		}, nil
	case "@monthly":
		return &SpecSchedule{
			Second:   1,
			Minute:   1,
			Hour:     1,
			Dom:      1,
			Month:    all(months),
			Dow:      all(dow),
			Location: time.Local,
		}, nil
	case "@weekly":
		return &SpecSchedule{
			Second:   1,
			Minute:   1,
			Hour:     1,
			Dom:      all(dom),
			Month:    all(months),
			Dow:      1,
			Location: time.Local,
		}, nil
	case "@daily", "@midnight":
		return &SpecSchedule{
			Second:   1,
			Minute:   1,
			Hour:     1,
			Dom:      all(dom),
			Month:    all(months),
			Dow:      all(dow),
			Location: time.Local,
		}, nil
	case "@hourly":
		return &SpecSchedule{
			Second:   1,
			Minute:   1,
			Hour:     all(hours),
			Dom:      all(dom),
			Month:    all(months),
			Dow:      all(dow),
			Location: time.Local,
		}, nil
	}

	const every = "@every "
	if strings.HasPrefix(spec, every) {
		// 简化版不支持 @every，返回错误提示
		return nil, fmt.Errorf("@every descriptor is not supported in simplified version")
	}

	return nil, fmt.Errorf("unrecognized descriptor: %s", spec)
}

var (
	seconds = bounds{0, 59, nil}
	minutes = bounds{0, 59, nil}
	hours   = bounds{0, 23, nil}
	dom     = bounds{1, 31, nil}
	months  = bounds{1, 12, map[string]uint{
		"jan": 1, "feb": 2, "mar": 3, "apr": 4, "may": 5, "jun": 6,
		"jul": 7, "aug": 8, "sep": 9, "oct": 10, "nov": 11, "dec": 12,
	}}
	dow = bounds{0, 6, map[string]uint{
		"sun": 0, "mon": 1, "tue": 2, "wed": 3, "thu": 4, "fri": 5, "sat": 6,
	}}
)

const starBit = 1 << 63

type bounds struct {
	min, max uint
	names    map[string]uint
}

func getField(field string, r bounds) (uint64, error) {
	var bits uint64
	ranges := strings.FieldsFunc(field, func(r rune) bool { return r == ',' })
	for _, expr := range ranges {
		bit, err := getRange(expr, r)
		if err != nil {
			return bits, err
		}
		bits |= bit
	}
	return bits, nil
}

func getRange(expr string, r bounds) (uint64, error) {
	var (
		start, end, step uint
		lowAndHigh       = strings.Split(expr, "-")
		rangeAndStep     = strings.Split(expr, "/")
		singleDigit      = len(lowAndHigh) == 1
		err              error
	)

	var extra uint64
	if lowAndHigh[0] == "*" || lowAndHigh[0] == "?" {
		start = r.min
		end = r.max
		extra = starBit
	} else {
		start, err = parseIntOrName(lowAndHigh[0], r.names)
		if err != nil {
			return 0, err
		}
		switch len(lowAndHigh) {
		case 1:
			end = start
		case 2:
			end, err = parseIntOrName(lowAndHigh[1], r.names)
			if err != nil {
				return 0, err
			}
		default:
			return 0, fmt.Errorf("too many hyphens: %s", expr)
		}
	}

	switch len(rangeAndStep) {
	case 1:
		step = 1
	case 2:
		step, err = mustParseInt(rangeAndStep[1])
		if err != nil {
			return 0, err
		}
		if singleDigit {
			end = r.max
		}
		if step > 1 {
			extra = 0
		}
	default:
		return 0, fmt.Errorf("too many slashes: %s", expr)
	}

	if start < r.min || end > r.max || start > end || step == 0 {
		return 0, fmt.Errorf("invalid range: %s", expr)
	}

	return getBits(start, end, step) | extra, nil
}

func all(r bounds) uint64 {
	return getBits(r.min, r.max, 1) | starBit
}

func getBits(min, max, step uint) uint64 {
	if step == 1 {
		return ^(math.MaxUint64 << (max + 1)) & (math.MaxUint64 << min)
	}
	var bits uint64
	for i := min; i <= max; i += step {
		bits |= 1 << i
	}
	return bits
}

func parseIntOrName(expr string, names map[string]uint) (uint, error) {
	if names != nil {
		if namedInt, ok := names[strings.ToLower(expr)]; ok {
			return namedInt, nil
		}
	}
	return mustParseInt(expr)
}

func mustParseInt(expr string) (uint, error) {
	num, err := strconv.Atoi(expr)
	if err != nil {
		return 0, err
	}
	if num < 0 {
		return 0, fmt.Errorf("negative number not allowed: %s", expr)
	}
	return uint(num), nil
}

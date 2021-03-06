package utils

import (
	"github.com/anywhy/medis/pkg/utils/log"
	"strconv"
)

func Argument(argsMap map[string]interface{}, name string) (string, bool) {
	if argsMap[name] != nil {
		if s, ok := argsMap[name].(string); ok {
			if s != "" {
				return s, true
			}
			log.Panicf("option %s requires an argument", name)
		} else {
			log.Panicf("option %s isn't a valid string", name)
		}
	}

	return "", false
}

func ArgumentMust(argsMap map[string]interface{}, name string) string {
	if v, ok := Argument(argsMap, name); ok {
		return v
	}

	log.Panicf("option %s is required", name)
	return ""
}

func ArgumentInt(argsMap map[string]interface{}, name string) (int, bool) {
	if v, ok := Argument(argsMap, name); ok {
		n, err := strconv.Atoi(v)

		if err != nil {
			log.PanicErrorf(err, "option %s isn't a valid integer", name)
		}

		return n, true
	}

	return 0, false
}

func ArgumentIntMust(argsMap map[string]interface{}, name string) int {
	if v, ok := ArgumentInt(argsMap, name); ok {
		return v
	}

	log.Panicf("option %s is required", name)
	return 0
}

func ArgumentBool(argsMap map[string]interface{}, name string) (bool, bool) {
	if v, ok := Argument(argsMap, name); ok {
		n, err := strconv.ParseBool(v)

		if err != nil {
			log.PanicErrorf(err, "option %s isn't a valid bool", name)
		}

		return n, true
	}

	return false, false
}

func ArgumentBoolMust(argsMap map[string]interface{}, name string) bool {
	if v, ok := ArgumentBool(argsMap, name); ok {
		return v
	}

	log.Panicf("option %s is required", name)
	return false
}

func ArgumentFloat64(argsMap map[string]interface{}, name string) (float64, bool) {
	if v, ok := Argument(argsMap, name); ok {
		n, err := strconv.ParseFloat(v, 64)

		if err != nil {
			log.PanicErrorf(err, "option %s isn't a valid bool", name)
		}

		return n, true
	}

	return 0, false
}

func ArgumentFloat64Must(argsMap map[string]interface{}, name string) float64 {
	if v, ok := ArgumentFloat64(argsMap, name); ok {
		return v
	}

	log.Panicf("option %s is required", name)
	return 0
}

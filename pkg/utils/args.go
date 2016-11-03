package utils

import "github.com/anywhy/medis/pkg/utils/log"

func argsparse(args map[string]interface{}, name string) (string, bool) {
	if args[name] != nil {
		if s, ok := args[name].(string); ok {
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

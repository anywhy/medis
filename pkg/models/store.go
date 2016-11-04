package models

import "path/filepath"

const StateDir  = "/state"

func FrameworkIdPath() string {
	return filepath.Join(StateDir, "frameworkId")
}

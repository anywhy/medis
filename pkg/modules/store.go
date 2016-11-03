package modules

import "path/filepath"

const MedisDir  = "/medis"

func FrameworkIdPath() string {
	return filepath.Join(MedisDir, "frameworkId")
}

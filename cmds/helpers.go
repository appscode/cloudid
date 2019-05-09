package cmds

import (
	"gomodules.xyz/cert"
)

func Filename(cfg cert.Config) string {
	if len(cfg.Organization) == 0 {
		return cfg.CommonName
	}
	return cfg.CommonName + "@" + cfg.Organization[0]
}

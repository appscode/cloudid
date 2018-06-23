package lib

func Scheme(tlsEnabled bool) string {
	if tlsEnabled {
		return "https"
	}
	return "http"
}

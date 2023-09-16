package subspace_experiment

import (
	"net/http"
)

func Init() {
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 1000
}

package subspace_experiment

import (
	"net/http"

	"github.com/xm0onh/subspace_experiment/log"
)

func Init() {
	log.Setup()
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 1000
}

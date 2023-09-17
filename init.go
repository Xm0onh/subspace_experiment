package subspace_experiment

import (
	"net/http"

	"github.com/xm0onh/subspace_experiment/config"
	"github.com/xm0onh/subspace_experiment/log"
)

func Init() {
	log.Setup()
	config.Configuration.Load()
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 1000
}

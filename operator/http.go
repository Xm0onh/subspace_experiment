package operator

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/xm0onh/subspace_experiment/config"
	"github.com/xm0onh/subspace_experiment/log"
)

func (o *Operator) http() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", report)
	ip, err := url.Parse(config.Configuration.HTTPAddrs[o.id])
	if err != nil {
		log.Fatal("http url parse error: ", err)
	}
	port := ":" + ip.Port()
	o.server = &http.Server{
		Addr:    port,
		Handler: mux,
	}
	log.Info("Node ", o.id, " http server starting on ", port)
	log.Fatal(o.server.ListenAndServe())
}

func report(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

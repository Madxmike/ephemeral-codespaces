package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/tierzer0/gateway/deployment"
	"github.com/tierzer0/gateway/helm"
)

var (
	port        string
	chartPath   string
	scriptsPath string
)

func main() {
	flag.StringVar(&port, "port", "3000", "the port the server will run on")
	flag.StringVar(&chartPath, "chartPath", "", "the path to the chart to deploy")
	flag.StringVar(&scriptsPath, "scriptsPath", "", "the path to the runtime scripts")

	flag.Parse()
	handler := deployment.Handler{
		Deployer: helm.CLIDeployer{
			ChartPath: chartPath,
		},
	}
	r := chi.NewRouter()
	r.Post("/", handler.ServeHTTP)

	s := http.Server{
		Handler: r,
		Addr:    "0.0.0.0:" + port,
	}

	log.Println(s.ListenAndServe())
}

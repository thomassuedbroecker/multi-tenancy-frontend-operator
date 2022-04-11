package controllers

import (
	"context"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/log"
)

func getGoobersTotal(PROMETHEUS_URL string, ctx context.Context) {
	log := log.FromContext(ctx)

	goobers_query := "goobers_total"
	prom_api_call := "api/v1/query?query="
	request_url := PROMETHEUS_URL + "/" + prom_api_call + goobers_query

	log.Info("Created request url: " + request_url)

	// 1. Create HTTP request
	req, err := http.NewRequest("GET", request_url, nil)
	if err != nil {
		log.Error(err, "Create HTTP request")
		return err, postgresTableExists
	}

	// 2. Define header
	req.Header.Set("Accept", "application/json")

	// 3. Create client
	client := http.Client{}

	// 4. Invoke HTTP request
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err, "Invoke HTTP request")
		return err, postgresTableExists
	}

}

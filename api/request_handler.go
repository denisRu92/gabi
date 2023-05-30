package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

func (api *API) getStats(writer http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method == http.MethodGet {

		JSON(writer, api.m.GetStates())
	} else {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (api *API) getSimilar(writer http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method == http.MethodGet {
		// Inc request counter
		api.m.IncRequestCounter()

		// Calculate processing time of request
		start := time.Now()
		defer api.m.AddProcessingTiming(start)

		// Get query parameter values
		queryParams := req.URL.Query()
		word := queryParams.Get("word")

		// Check if required query parameters are missing
		if word == "" {
			http.Error(writer, "Missing 'word' query parameter", http.StatusBadRequest)
			return
		}

		JSON(writer, map[string]interface{}{"similar": api.service.GetSimilar(word)})
	} else {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

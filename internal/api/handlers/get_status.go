package handlers

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/Quantum12k/healthcheck-service/internal/api/response"
	"github.com/Quantum12k/healthcheck-service/internal/app_cache"
)

func GetStatus(cache *app_cache.ChecksCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Render(w, r, response.RenderResp(cache.GetMapCopy(), http.StatusOK))
	}
}

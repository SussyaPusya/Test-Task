package clients

import (
	"net/http"
	"test_task/internal/config"
)

type ExternalAPI struct {
	urls       *config.ExternalAPI
	httpClient *http.Client
}

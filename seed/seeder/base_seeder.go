package seeder

import (
	"net/http"
	"os"
	"time"
)

type Seeder interface {
	Seed()
	Name() string
	IsServiceReady() bool
}

type BaseSeeder struct {
	serviceName string
	serviceURL  string
	client      *http.Client
}

func NewBaseSeeder(serviceName, envKey string) *BaseSeeder {
	return &BaseSeeder{
		serviceName: serviceName,
		serviceURL:  os.Getenv(envKey),
		client:      &http.Client{Timeout: 10 * time.Second},
	}
}

func (b *BaseSeeder) Name() string {
	return b.serviceName
}

func (b *BaseSeeder) IsEnabled() bool {
	return b.serviceURL != ""
}

func (b *BaseSeeder) IsServiceReady() bool {
	if !b.IsEnabled() {
		return false
	}

	pingURL := b.serviceURL + "/internal/ping"
	for i := 0; i < 15; i++ {
		resp, err := b.client.Get(pingURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return true
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}
	return false
}

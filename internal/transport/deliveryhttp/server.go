package deliveryhttp

import (
	"github.com/RedWood011/ServiceURL/internal/service"
)

type Server struct {
	repository service.Storage
	Addr       string
}

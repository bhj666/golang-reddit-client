package timeprovider

import "time"

type Provider interface {
	GetCurrentSeconds() int64
}

type ProviderImpl struct {
}

func (ProviderImpl) GetCurrentSeconds() int64 {
	return time.Now().Unix()
}

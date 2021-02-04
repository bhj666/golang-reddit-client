package test

type TimeProviderMock struct {
	Time int64
}

func (m TimeProviderMock) GetCurrentSeconds() int64 {
	return m.Time
}

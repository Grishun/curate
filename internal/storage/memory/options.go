package memory

type Options struct {
	HistoryLimit uint64 `json:"history_limit"`
}

type Option func(options *Options)

func NewOptions() *Options {
	return &Options{
		HistoryLimit: 10,
	}
}

func WithHistoryLimit(limit uint64) Option {
	return func(options *Options) { options.HistoryLimit = limit }
}

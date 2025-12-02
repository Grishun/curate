package memory

type Option func(*Options)

type Options struct {
	historyLimit uint32
}

func WithHistoryLimit(l uint32) Option {
	return func(options *Options) {
		options.historyLimit = l
	}
}

func NewOptions() *Options {
	return &Options{historyLimit: 10}
}

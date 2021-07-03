package mr

const (
	defaultWorkers = 16
	minWorkers     = 1
)

type Option func(opts *mapReduceOptions)

type mapReduceOptions struct {
	workers int
}

func newOptions() *mapReduceOptions {
	return &mapReduceOptions{
		workers: defaultWorkers,
	}
}

func buildOptions(opts ...Option) *mapReduceOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

func WithWorkers(workers int) Option {
	return func(opts *mapReduceOptions) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

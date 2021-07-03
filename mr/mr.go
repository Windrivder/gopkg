package mr

import (
	"errors"
	"fmt"
	"sync"

	"github.com/windrivder/gopkg/container/typex"
	"github.com/windrivder/gopkg/errorx"
	"github.com/windrivder/gopkg/syncx"
	"github.com/windrivder/gopkg/threading"
)

var (
	ErrCancelWithNil  = errors.New("mapreduce cancelled with nil")
	ErrReduceNoOutput = errors.New("reduce not writing value")
)

type (
	GenerateFunc    func(source chan<- typex.GenericType)
	MapFunc         func(item typex.GenericType, writer Writer)
	VoidMapFunc     func(item typex.GenericType)
	MapperFunc      func(item typex.GenericType, writer Writer, cancel func(error))
	ReducerFunc     func(pipe <-chan typex.GenericType, writer Writer, cancel func(error))
	VoidReducerFunc func(pipe <-chan typex.GenericType, cancel func(error))
)

func Finish(fns ...func() error) error {
	workers := len(fns)
	if workers == 0 {
		return nil
	}

	// source
	return MapReduceVoid(func(source chan<- typex.GenericType) {
		for _, fn := range fns {
			source <- fn
		}

		// reducer
	}, func(item typex.GenericType, writer Writer, cancel func(error)) {
		fn := item.(func() error)
		if err := fn(); err != nil {
			cancel(err)
		}

		// reducer
	}, func(pipe <-chan typex.GenericType, cancel func(error)) {
		drain(pipe)

	}, WithWorkers(workers))
}

func FinishVoid(fns ...func()) {
	workers := len(fns)
	if workers == 0 {
		return
	}

	MapVoid(func(source chan<- typex.GenericType) {
		for _, fn := range fns {
			source <- fn
		}

	}, func(item typex.GenericType) {
		fn := item.(func())
		fn()

	}, WithWorkers(workers))
}

func Map(generate GenerateFunc, mapper MapFunc, opts ...Option) chan typex.GenericType {
	options := buildOptions(opts...)

	source := buildSource(generate)
	collector := make(chan typex.GenericType, options.workers)
	done := syncx.NewDoneChan()

	go executeMappers(mapper, source, collector, done.Done(), options.workers)

	return collector
}

func MapVoid(generate GenerateFunc, mapper VoidMapFunc, opts ...Option) {
	drain(Map(generate, func(item typex.GenericType, writer Writer) {
		mapper(item)
	}, opts...))
}

func MapReduce(generate GenerateFunc, mapper MapperFunc, reducer ReducerFunc, opts ...Option) (typex.GenericType, error) {
	source := buildSource(generate)
	return MapReduceWithSource(source, mapper, reducer, opts...)
}

func MapReduceVoid(generate GenerateFunc, mapper MapperFunc, reducer VoidReducerFunc, opts ...Option) error {
	rder := func(input <-chan typex.GenericType, writer Writer, cancel func(error)) {
		reducer(input, cancel)
		drain(input)
		// We need to write a placeholder to let MapReduce to continue on reducer done,
		// otherwise, all goroutines are waiting. The placeholder will be discarded by MapReduce.
		writer.Write(typex.Placeholder)
	}
	_, err := MapReduce(generate, mapper, rder, opts...)

	return err
}

func MapReduceWithSource(source <-chan typex.GenericType, mapper MapperFunc, reducer ReducerFunc, opts ...Option) (typex.GenericType, error) {
	// reduce writer
	output := make(chan typex.GenericType)
	done := syncx.NewDoneChan()
	writer := newGuardedWriter(output, done.Done())

	//
	var closeOnce sync.Once
	finish := func() {
		closeOnce.Do(func() {
			done.Close()
			close(output)
		})
	}

	//
	var retErr errorx.AtomicError
	cancel := once(func(err error) {
		if err != nil {
			retErr.Set(err)
		} else {
			retErr.Set(ErrCancelWithNil)
		}

		drain(source)
		finish()
	})

	options := buildOptions(opts...)
	collector := make(chan typex.GenericType, options.workers)

	// reduce
	go func() {
		defer func() {
			if r := recover(); r != nil {
				cancel(fmt.Errorf("%v", r))
			} else {
				finish()
			}
		}()

		reducer(collector, writer, cancel)
		drain(collector)
	}()

	// map
	go executeMappers(
		func(item typex.GenericType, w Writer) {
			mapper(item, w, cancel)
		},
		source,
		collector,
		done.Done(),
		options.workers,
	)

	// return
	value, ok := <-output
	if err := retErr.Load(); err != nil {
		return nil, err
	} else if ok {
		return value, nil
	} else {
		return nil, ErrReduceNoOutput
	}
}

func buildSource(generate GenerateFunc) chan typex.GenericType {
	source := make(chan typex.GenericType)

	threading.GoSafe(func() {
		defer close(source)
		generate(source)
	})

	return source
}

// drain drains the channel.
func drain(channel <-chan typex.GenericType) {
	// drain the channel
	for range channel {
	}
}

func once(fn func(error)) func(error) {
	once := new(sync.Once)
	return func(err error) {
		once.Do(func() {
			fn(err)
		})
	}
}

func executeMappers(mapper MapFunc, source <-chan typex.GenericType, collector chan<- typex.GenericType, done <-chan typex.PlaceholderType, workers int) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		close(collector)
	}()

	pool := make(chan typex.PlaceholderType, workers)
	writer := newGuardedWriter(collector, done)
	for {
		select {
		case <-done:
			return
		case pool <- typex.Placeholder:
			item, ok := <-source
			if !ok {
				<-pool
				return
			}

			wg.Add(1)
			threading.GoSafe(func() {
				defer func() {
					wg.Done()
					<-pool
				}()

				mapper(item, writer)
			})
		}
	}
}

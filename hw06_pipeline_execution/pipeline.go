package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var out Out
	stopConsumer := true
	for _, stage := range stages {
		if done == nil {
			out = stage(in)
		} else {
			out = stage(outWrap(in, done, stopConsumer))
			stopConsumer = false
		}
		in = out
	}
	return out
}

func outWrap(in In, done In, stopConsume bool) Out {
	out := make(Bi)

	go func() {
		defer close(out)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			case <-done:
				if !stopConsume {
					go func() {
						for {
							_, ok := <-in
							if !ok {
								return
							}
						}
					}()
				}
				return
			}
		}
	}()

	return out
}

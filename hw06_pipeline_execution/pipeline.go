package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = doneWrapper(done, stage(out))
	}
	return out
}

func doneWrapper(done, in In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			for range in { // Draining of wrapped stage
			}
		}()
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			case <-done:
				return
			}
		}
	}()
	return out
}

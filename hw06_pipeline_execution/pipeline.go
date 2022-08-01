package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		stageCh := in
		for _, stage := range stages {
			stageCh = stage(stageCh)
		}
		for {
			select {
			case v, ok := <-stageCh:
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

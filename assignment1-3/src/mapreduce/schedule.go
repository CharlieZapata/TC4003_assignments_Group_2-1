package mapreduce

import (
	"sync"
)

// schedule starts and waits for all tasks in the given phase (Map or Reduce).
func (mr *Master) schedule(phase jobPhase) {
	var ntasks int
	var nios int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mr.files)
		nios = mr.nReduce
	case reducePhase:
		ntasks = mr.nReduce
		nios = len(mr.files)
	}
	debug("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, nios)
	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	var wg sync.WaitGroup         //this variable is used to controll the go routine flow
	for i := 0; i < ntasks; i++ { //Iterate over the total number of tasks that will beocme workers at the first run
		doTaskArgs := new(DoTaskArgs)
		doTaskArgs.JobName = mr.jobName
		doTaskArgs.Phase = phase
		doTaskArgs.TaskNumber = i
		doTaskArgs.NumOtherPhase = nios
		if phase == mapPhase {
			doTaskArgs.File = mr.files[i]
		} //Initialize the workers

		wg.Add(1)   //Add each go routine to the control variable
		go func() { //Definition of the go routine that will be executed for each worker
			for { //This while loop receives until a registered worker is available in the channel
				currentWorker := <-mr.registerChannel
				rpcStatus := call(currentWorker, "Worker.DoTask", &doTaskArgs, new(struct{})) //The RPC call is performed, it returns a boolean based on its result
				if rpcStatus {
					wg.Done()                           //Notify the routine its compketed
					mr.registerChannel <- currentWorker //Return the worker to the channel
					break
				}
			}
		}()
	}
	wg.Wait() //Waits for all the routines to finish
	debug("Schedule: %v phase done\n", phase)
}

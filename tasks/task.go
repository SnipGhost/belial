package tasks

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

type statetype uint8

const (
	sWaiting  statetype = iota // State before starting the worker
	sRunning                   // State after starting the worker
	sReady                     // State after successful worker completion
	sCanceled                  // State after context canceling
	sError                     // State if task has erorr
)

// Task - test struct
type Task struct {
	sync.Mutex
	Name   string             // Task title
	Data   string             // Data (code)
	N      uint64             // n
	K      uint64             // k
	Type   uint64             // Type = {0|1}
	State  statetype          // Ready or not
	Stdout bytes.Buffer       // Buffer to printing
	ctx    context.Context    // Context
	cancel context.CancelFunc // Function to canceling context
}

// NewTask - return initialized Task
func NewTask(data, n, k, t string, timeout uint64) *Task {
	workTime := time.Duration(timeout) * time.Millisecond
	ctx, finish := context.WithTimeout(context.Background(), workTime)
	var b bytes.Buffer
	state := sWaiting
	num, err := strconv.ParseUint(n, 10, 64)
	if err != nil {
		state = sError
		fmt.Fprintf(&b, "Parse error: %s", err)
	}
	kum, err := strconv.ParseUint(k, 10, 64)
	if err != nil {
		state = sError
		fmt.Fprintf(&b, "Parse error: %s", err)
	}
	method, err := strconv.ParseUint(t, 10, 64)
	if err != nil {
		state = sError
		fmt.Fprintf(&b, "Parse error: %s", err)
	}
	return &Task{
		ctx:    ctx,
		cancel: finish,
		Name:   fmt.Sprintf("%s : Ц(%s, %s)", data, n, k),
		Data:   data,
		N:      num,
		K:      kum,
		Type:   method,
		State:  state,
		Stdout: b,
	}
}

// PrintInfo - for template generation
func (t *Task) PrintInfo() string {
	switch t.State {
	case sWaiting:
		return "Scheduling"
	case sRunning:
		return "Running"
	case sReady:
		return "Ready"
	case sCanceled:
		return "Canceled"
	case sError:
		return "Failed"
	default:
		return "Unknown"
	}
}

func (t *Task) check() byte {
	encrypted, err := encryptWrapper(t.Data, t.N, t.K, t.Type)
	if err != nil {
		t.Lock()
		fmt.Fprintln(&t.Stdout, "Error with encrypting:", err)
		t.State = sError
		t.Unlock()
		return 1
	}
	t.Lock()
	fmt.Fprintf(&t.Stdout, "Encrypted: %s\n", bitsToStr(encrypted, t.N))
	t.Unlock()
	var e, rem, i uint64
	var val float64
	var limit uint64 = (1 << t.N) - 1
	success := make([]uint64, t.N+1)
	for e = 0; e <= limit; e++ {
		select {
		case <-t.ctx.Done():
			t.Lock()
			log.Println("Task:", t.Name, "finished by ctx")
			fmt.Fprintln(&t.Stdout, "Timeout error (Finished by ctx)")
			t.State = sCanceled
			t.Unlock()
			return 1
		default:
		}
		_, rem, err = decryptWrapper(encrypted^e, t.N, t.K, t.Type)
		if err != nil {
			t.Lock()
			fmt.Fprintln(&t.Stdout, "Error with decrypting:", err, e)
			t.State = sError
			t.Unlock()
			return 1
		}
		if (e != 0) && (rem != 0) || (e == 0) && (rem == 0) {
			success[weightBits(e)]++
		}
	}
	for i = 0; i <= t.N; i++ {
		select {
		case <-t.ctx.Done():
			t.Lock()
			log.Println("Task:", t.Name, "finished by ctx")
			fmt.Fprintln(&t.Stdout, "Timeout error (Finished by ctx)")
			t.State = sCanceled
			t.Unlock()
			return 1
		default:
		}
		suc := success[i]
		com := countCombinations(t.N, i)
		val = float64(suc) / float64(com) * 100.0
		t.Lock()
		fmt.Fprintf(&t.Stdout, "%d bits in err vector, combinations: %d, finded: %d - %3.2f\n", i, com, suc, val)
		t.Unlock()
	}
	return 0
}

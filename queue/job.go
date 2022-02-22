package queue

import "log"

// Job - holds logic to perform some operations during queue execution.
type Job struct {
	Name    string
	Payload []byte
	Action  func(payload []byte) error
}

func (j Job) GetName() string {
	return j.Name
}

// Run performs job execution.
func (j Job) Run() error {
	log.Printf("Job running: %s", j.GetName())

	err := j.Action(j.Payload)
	if err != nil {
		return err
	}

	return nil
}

package agents

import "fmt"

type Monetary1 struct {
	AgentAbs
}

func (m1 *Monetary1) Execute() {
	fmt.Println("Monetary1 agent is executing...")
}

package main

import (
	"fmt"
	"os"
	"os/signal"
	"proposalsubmitters/agents"
	"proposalsubmitters/utils"
	"runtime"
	"syscall"
	"time"
)

type Server struct {
	quit   chan os.Signal
	finish chan bool
	agents []agents.Agent
}

func NewServer() *Server {
	rpcClient := utils.NewHttpClient()

	f1 := &agents.Fiscal1{}
	f1.ID = 1
	f1.Name = "fiscal agent 1"
	f1.Frequency = 10
	f1.Quit = make(chan bool)
	f1.RPCClient = rpcClient

	m1 := &agents.CascadingAgent{}
	m1.ID = 2
	m1.Name = "cascading agent 1"
	m1.Frequency = 20
	m1.Quit = make(chan bool)
	m1.RPCClient = rpcClient

	agents := []agents.Agent{f1, m1}
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, syscall.SIGTERM)
	signal.Notify(quitChan, syscall.SIGINT)
	return &Server{
		quit:   quitChan,
		finish: make(chan bool, len(agents)),
		agents: agents,
	}
}

func (s *Server) NotifyQuitSignal(agents []agents.Agent) {
	sig := <-s.quit
	fmt.Printf("Caught sig: %+v \n", sig)
	// notify all agents about quit signal
	for _, a := range agents {
		a.GetQuitChan() <- true
	}
}

func (s *Server) Run() {
	agents := s.agents
	go s.NotifyQuitSignal(agents)
	for _, a := range agents {
		go executeAgent(s.finish, a)
	}
}

func executeAgent(
	finish chan bool,
	agent agents.Agent,
) {
	agent.Execute() // execute as soon as starting up
	for {
		select {
		case <-agent.GetQuitChan():
			fmt.Printf("Finishing task for %s ...\n", agent.GetName())
			time.Sleep(time.Second * 1)
			fmt.Printf("Task for %s done! \n", agent.GetName())
			finish <- true
			break
		case <-time.After(time.Duration(agent.GetFrequency()) * time.Second):
			agent.Execute()
		}
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	s := NewServer()
	s.Run()
	for range s.agents {
		<-s.finish
	}
	fmt.Println("Server stopped gracefully!")
}

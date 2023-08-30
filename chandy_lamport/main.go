package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type BankAccount struct {
	AccountID int
	Balance   float64
}

type Process struct {
	ProcessID       int
	Accounts        []BankAccount
	IncomingChan    chan float64
	ChannelStates   map[int]float64
	RandomGenerator *rand.Rand
}

func (p *Process) transaction(processes *[]Process, wg *sync.WaitGroup) {
	defer wg.Done()

	numTransactions := len(p.Accounts) * 15

	for i := 0; i < numTransactions; i++ {
		amount := float64(p.RandomGenerator.Intn(500)) // Random amount between 0 and 500
		recipientProcess := p.RandomGenerator.Intn(len(*processes))

		p.ChannelStates[p.ProcessID] += amount      // Record outgoing amount
		p.ChannelStates[recipientProcess] -= amount // Record incoming amount

		fmt.Printf("Process %d sent %.2f to Process %d\n", p.ProcessID, amount, recipientProcess)

		if (i+1)%15 == 0 {
			p.takeSnapshot()
		}
	}
}

func (p *Process) takeSnapshot() {
	fmt.Printf("Process %d is taking a snapshot.\n", p.ProcessID)

	// Simulate snapshot taking process
	time.Sleep(time.Millisecond * 100)

	// Calculate the total amount in the snapshot
	totalSnapshotAmount := 0.0
	for _, account := range p.Accounts {
		totalSnapshotAmount += account.Balance
	}

	fmt.Printf("Process %d snapshot total amount: %.2f\n", p.ProcessID, totalSnapshotAmount)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var n int
	fmt.Print("Enter the number of processes: ")
	fmt.Scan(&n)

	var totalAmount float64

	var processes []Process
	var waitGroup sync.WaitGroup

	for i := 0; i < n; i++ {
		accounts := []BankAccount{
			{AccountID: 0, Balance: 1000.00},
			{AccountID: 1, Balance: 1000.00},
			{AccountID: 2, Balance: 1000.00},
		}

		incomingChan := make(chan float64)
		channelStates := make(map[int]float64)

		process := Process{
			ProcessID:       i,
			Accounts:        accounts,
			IncomingChan:    incomingChan,
			ChannelStates:   channelStates,
			RandomGenerator: rand.New(rand.NewSource(time.Now().UnixNano())),
		}

		processes = append(processes, process)

		for _, account := range process.Accounts {
			totalAmount += account.Balance
		}
	}

	fmt.Printf("Total amount in the system before starting: %.2f\n", totalAmount)

	waitGroup.Add(n)
	for i := 0; i < n; i++ {
		go processes[i].transaction(&processes, &waitGroup)
	}
	waitGroup.Wait()

	fmt.Println("All processes have finished.")
}

package main

import (
	"errors"
	"sync"
)

type PointsRepositoryMock struct {
	mutex sync.Mutex

	AmountCents int
}

func (p *PointsRepositoryMock) Charge(chargedAmountInCents int) (int, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Simulate a row-lock SQL query (usually using `FOR UPDATE` statement)
	totalPoints := p.AmountCents

	// Ensure the amount is sufficient to be charged
	if (totalPoints < chargedAmountInCents) {
		return 0, errors.New("insufficient amount")
	}

	// Simulate an atomic update (the atomic part here is that we already lock this transaction using Mutex)
	totalPoints = totalPoints - chargedAmountInCents

	return totalPoints, nil
}

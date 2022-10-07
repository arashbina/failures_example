package logic

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Rule struct{}

var storage = map[uuid.UUID]bool{}

func CreateRule() (uuid.UUID, error) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	// randomly fail this rule creation in 10% of the cases
	if a := r1.Intn(100); a < 10 {
		e := fmt.Sprintf("failed: %d", a)
		return uuid.UUID{}, errors.New(e)
	}

	u := uuid.New()
	storage[u] = true

	return u, nil
}

func Delete(u uuid.UUID) error {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	// randomly fail this rule deletion in 3% of the cases
	if a := r1.Intn(100); a < 3 {
		return errors.New("some error")
	}

	delete(storage, u)
	return nil
}

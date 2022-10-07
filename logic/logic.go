package logic

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"time"
)

type Rule struct{}

var storage = map[uuid.UUID]bool{}

func CreateRule() (uuid.UUID, error) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	u := uuid.New()

	// randomly fail this rule creation in 10% of the cases
	if a := r1.Intn(100); a < 10 {
		e := fmt.Sprintf("failed creation of rule with id: %s", u)
		return uuid.UUID{}, errors.New(e)
	}

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
	log.Printf("deleted a rule with id: %s\n", u)
	return nil
}

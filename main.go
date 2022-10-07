package main

import (
	"encoding/json"
	"failures_example/logic"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Journey struct {
	rules  []uuid.UUID
	vacuum bool
	sync.RWMutex
	errors []error
}

var vacuum = make(chan Journey)

func main() {

	http.HandleFunc("/api/journey", CreateJourney)
	http.HandleFunc("/api/state", GetState)

	ticker := time.NewTicker(2 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				for j := range vacuum {
					if j.TryLock() {
						defer j.Unlock()
						var notClean bool
						for _, r := range j.rules {
							if err := logic.Delete(r); err != nil {
								notClean = true
							}
						}
						if notClean {
							vacuum <- j
						}
					}
				}
			}
		}
	}()

	http.ListenAndServe(":8080", nil)
}

func CreateJourney(w http.ResponseWriter, r *http.Request) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	j := Journey{}

	num := r1.Intn(5)
	fmt.Printf("will create %d rules\n", num)

	j.Lock()
	defer j.Unlock()

	for i := 0; i < num; i++ {
		r, err := logic.CreateRule()
		if err != nil {
			j.errors = append(j.errors, err)
			j.vacuum = true
			continue
		}
		j.rules = append(j.rules, r)
	}

	if j.vacuum {
		vacuum <- j
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(j.errors)
		return
	}

	fmt.Println("create the journey")
	w.WriteHeader(http.StatusOK)
}

func GetState(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("number of items to be vacuumed: %d", len(vacuum))
	w.Write([]byte(s))
	w.WriteHeader(http.StatusOK)
}

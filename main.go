package main

import (
	"failures_example/logic"
	"fmt"
	"github.com/google/uuid"
	"log"
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

var vacuum = make(chan *Journey)

func main() {

	http.HandleFunc("/api/journey", CreateJourney)
	http.HandleFunc("/api/state", GetState)

	go func() {
		for j := range vacuum {
			log.Println("vacuum ran")
			if j.TryLock() {
				defer j.Unlock()
				log.Println("journey locked and will vacuum")
				defer j.Unlock()
				var notClean bool
				for _, r := range j.rules {
					if err := logic.Delete(r); err != nil {
						notClean = true
					}
				}
				if notClean {
					vacuum <- j
				} else {
					log.Println("vacuumed journies")
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

	// create a random number of rules for this journey
	num := r1.Intn(5)
	log.Printf("will create %d rules for this journey\n", num)

	// lock the journey so that no one else can modify it
	j.Lock()
	defer j.Unlock()

	for i := 0; i < num; i++ {
		r, err := logic.CreateRule()
		if err != nil {
			// add the error to the slice of errors
			// and mark the journey to be vacuumed
			j.errors = append(j.errors, err)
			j.vacuum = true
			continue
		}
		j.rules = append(j.rules, r)
	}

	if j.vacuum {
		// add the journey to the vacuum channel
		vacuum <- &j
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("errors: %v\n", j.errors)
		w.Write([]byte("failed creating journey"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("created journey successfully"))
}

func GetState(w http.ResponseWriter, r *http.Request) {
	// return the number of journies that will get vacuumed
	s := fmt.Sprintf("number of items to be vacuumed: %d", len(vacuum))
	w.Write([]byte(s))
	w.WriteHeader(http.StatusOK)
}

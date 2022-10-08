# failures_example

Create journies by calling:

    curl -v localhost:8080/api/journey
    
A journy gets created with a number of rules but 10% of rule creations fail which causes the journey creation to fail. The failed journies get vacummed asynchronously.
Some of rule deletions when a journey is being vacuumed may also fail (about 3% of them) but the journey with the remaining rules get added back to the journies to be vacuumed.

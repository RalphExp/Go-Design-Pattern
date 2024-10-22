.PHONY: test

test_pipeline1:
	go test -v pipeline1_test.go pipeline1.go

test_pipeline2:
	go test -v -bench=. pipeline2_test.go pipeline2.go

test_prime1:
	go test -v pipeline2.go prime.go prime_test.go

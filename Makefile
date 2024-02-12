.PHONY: coverage

default: coverage

go_test_with_coverprofile:
	go test ./... -coverprofile=./coverage/out/coverage.out

exclude_from_coverage:
	./coverage/exclude.sh

go_coverage_html:
	go tool cover -html=./coverage/out/coverage.out -o=./coverage/out/coverage.html

coverage: go_test_with_coverprofile exclude_from_coverage go_coverage_html

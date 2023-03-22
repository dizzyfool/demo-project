-include .env
export

run:
	go run pkg/main.go
test:
	go test mackey/pkg/demo
mfd-xml:
	@mfd-generator xml -c "postgres://postgres:postgres@localhost:5432/demo?sslmode=disable" -m ./docs/model/demo.mfd
mfd-model:
	@mfd-generator model -m ./docs/model/demo.mfd -p db -o ./pkg/db
mfd-repo:
	@mfd-generator repo -m ./docs/model/demo.mfd -p db -o ./pkg/db
.PHONY: all build generate tools quickbuild

BINARY=hrms
TOOL=hrmstools
OUTPUT=build/exec
MAIN=main.go  # Corrected path
DCONFIG := config/connection/dev-config.json 

local: runOnWindowsL
dev_sitl: runOnWindowsSD
prod_sitl: runOnWindowsP
dev: runOnWindowsD
run: runOnWindowsL
all: quickbuild

wbr: winBuild runOnWindows 

server:
	go mod download
	go build -o ${OUTPUT}/${BINARY}.exe ${MAIN}
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags "-w -extldflags '-static' " -o ${OUTPUT}/${BINARY} 

quickbuild:  
	rm -rf ${OUTPUT}/${BINARY}
	mkdir -p ${OUTPUT}
	go build -o ${OUTPUT}/${BINARY} cmd/**/*.go

tools:
	rm -rf ${OUTPUT}/exec/${TOOL}
	mkdir -p ${OUTPUT}/exec
	go build -o ${OUTPUT}/exec/${TOOL} tool/main.go
	
clean:
	rm -rf ${OUTPUT}/exec/${BINARY}
	rm -rf ${OUTPUT}/exec/${TOOL}
	rm -rf ${OUTPUT}/${BINARY}
	rm -rf ${OUTPUT}/${BINARY}.exe

sqlc:
	sqlc generate --file config/sqlc/db_query/db.yaml

generate:
	cd configs/sqlc && sqlc generate -f ./auth-sqlc.yaml 
	cd configs/sqlc && sqlc generate -f ./hr-sqlc.yaml
	cd configs/sqlc && sqlc generate -f ./salary-sqlc.yaml

swag:
	swag init -g main.go --output docs

# Windows build
winBuild:
	go mod download
	go mod tidy
	go build -o ${OUTPUT}/${BINARY}.exe ${MAIN}

# Use 'make runOnWindows' to execute program
# $(MAKE) winBuild

runOnWindowsD:
	go build -o ${OUTPUT}/${BINARY}.exe ${MAIN}
	./${OUTPUT}/${BINARY}.exe -c ./config/connection/dev-config.json

runOnWindowsP:
	go build -o ${OUTPUT}/${BINARY}.exe ${MAIN}
	./${OUTPUT}/${BINARY}.exe -c ./config/connection/dev-config.json --port 8035

lb:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags "-w -extldflags '-static' " -o $(OUTPUT)/$(BINARY) .

lr:
	./$(OUTPUT)/$(BINARY) -c $(DCONFIG)

lbr: lb lr 








MAIN_FILE      = main.go
SQLC_DIR       = config/sqlc/db_query
SQLC_CONFIG    = $(SQLC_DIR)/db.yaml
CONFIG_DIR     = config/connection
DB_MODEL_DIR   = internal/db_model
MODEL_DIR      = internal/model
SERVICE_DIR    = internal/service

path:
	@echo "Main file:            $(MAIN_FILE)"
	@echo "SQLC directory:       $(SQLC_DIR)"
	@echo "SQLC config:          $(SQLC_CONFIG)"
	@echo "Config directory:     $(CONFIG_DIR)"
	@echo "SQLC generated files: $(DB_MODEL_DIR)"
	@echo "Custom Models:        $(MODEL_DIR)"
	@echo "Services:             $(SERVICE_DIR)"
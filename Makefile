MAKEFLAGS += --silent

EXECUTABLE=postgres-command-executor
EXECUTABLE_DIR=bin
SOURCE=cmd/$(EXECUTABLE)/main.go

all: $(EXECUTABLE)


$(EXECUTABLE): create-executable-dir-if-not-exists
	echo "Building Golang project..."
	go build -o $(EXECUTABLE_DIR)/$(EXECUTABLE) $(SOURCE)
	echo "Golang project was build successfully Executable file is located in '$(EXECUTABLE_DIR)/$(EXECUTABLE)'"


create-executable-dir-if-not-exists:
ifeq (, $(wildcard $(EXECUTABLE_DIR)))
	echo "Creating '$(EXECUTABLE_DIR)' directory..."
	mkdir $(EXECUTABLE_DIR)
	echo "'$(EXECUTABLE_DIR)' directory created successfully"
endif

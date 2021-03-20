.PHONY: run
run:
	@echo "Building BiliBan..."
	@mkdir -p ./build
	@go clean
	@rm -f ./build/main
	@go build -o ./build/main .
	@echo "Finish! Ready to start server..."
	@./build/main

build:
    @echo "Building BiliBan..."
	@mkdir -p ./build
	@go clean
	@rm -f ./build/main
	@go build -o ./build/main .
	@echo "Done! Build finish!"
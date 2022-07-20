
# List
default:
    @just --list

# Run the server
run:
    go run ./cmd/server

# Against the locally running server... run a drip
nc-drip:
    echo 'drip 3 3 3' | nc localhost 4246

# Build the teeceepee server
buildserver:
    go build -o teeceepee ./cmd/server

# Attempt to build the library
build:
    go build .
#!/bin/bash

# Default configuration
DEFAULT_DB_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"
DEFAULT_PORT=8080

# Load environment variables from .env if it exists
if [ -f .env ]; then
    source .env
fi

# Use environment variables or defaults
DB_URL=${DATABASE_URL:-$DEFAULT_DB_URL}
BASE_PORT=${BASE_PORT:-$DEFAULT_PORT}

# Function to start a lambda
start_lambda() {
    local lambda_name=$1
    local port=$2
    local dir="lambdas/$lambda_name"
    
    echo "Starting $lambda_name on port $port..."
    
    # Create a temporary .env file for this lambda
    cat > "$dir/.env" << EOF
DATABASE_URL=$DB_URL
PORT=$port
EOF
    
    # Start the lambda in the background
    cd "$dir" && go run main.go &
    
    # Store the PID
    echo $! > "$dir/.pid"
    
    # Return to original directory
    cd - > /dev/null
}

# Function to stop a lambda
stop_lambda() {
    local lambda_name=$1
    local pid_file="lambdas/$lambda_name/.pid"
    
    if [ -f "$pid_file" ]; then
        echo "Stopping $lambda_name..."
        kill $(cat "$pid_file") 2>/dev/null || true
        rm "$pid_file"
    fi
}

# Function to cleanup
cleanup() {
    echo "Cleaning up..."
    for lambda in user_create user_read user_update user_delete send_email log_event; do
        stop_lambda $lambda
        rm -f "lambdas/$lambda/.env"
    done
    exit 0
}

# Set up cleanup on script exit
trap cleanup EXIT INT TERM

# Start each lambda on a different port
start_lambda "user_create" $((BASE_PORT))
start_lambda "user_read" $((BASE_PORT + 1))
start_lambda "user_update" $((BASE_PORT + 2))
start_lambda "user_delete" $((BASE_PORT + 3))

echo "All lambdas started. Press Ctrl+C to stop."
echo
echo "Available endpoints:"
echo "  Create user:   http://localhost:$BASE_PORT/"
echo "  Read user:     http://localhost:$((BASE_PORT + 1))/<id>"
echo "  Update user:   http://localhost:$((BASE_PORT + 2))/<id>"
echo "  Delete user:   http://localhost:$((BASE_PORT + 3))/<id>"

echo
echo "Example usage:"
echo "  # Create a user"
echo "  curl -X POST -H \"Content-Type: application/json\" -d '{\"email\":\"user@example.com\",\"name\":\"John Doe\"}' http://localhost:$BASE_PORT/"
echo
echo "  # Get a user"
echo "  curl http://localhost:$((BASE_PORT + 1))/1"
echo
echo "  # Update a user"
echo "  curl -X PUT -H \"Content-Type: application/json\" -d '{\"email\":\"new@example.com\",\"name\":\"John Updated\"}' http://localhost:$((BASE_PORT + 2))/1"
echo
echo "  # Delete a user"
echo "  curl -X DELETE http://localhost:$((BASE_PORT + 3))/1"
echo


# Keep script running
while true; do
    sleep 1
done 
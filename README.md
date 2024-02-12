# User Statistics and Change Logger

## Overview

This project periodically logs summary statistics and changes in a MongoDB collection. It's written in Go and uses the following libraries:

- `github.com/joho/godotenv` for loading environment variables
- `github.com/robfig/cron/v3` for scheduling tasks
- `go.mongodb.org/mongo-driver/mongo` for interacting with MongoDB

## Features

- Logs the total number of users in the collection at a configurable interval
- Detects changes in the user count and logs the details of new users
- Runs as a background service until interrupted

## Getting Started

1. **Requirements:**
   - Go environment (version 1.17 or later recommended)
   - MongoDB database
2. **Configuration:**
   - Create a `.env` file in the project directory with the following variables:
     - `MONGO_URI`: The connection URI for your MongoDB database
     - `DB_NAME`: The name of the database containing the collection
     - `COLLECTION_NAME`: The name of the collection to monitor
     - `CRON_INTERVAL` (optional): The cron expression for task scheduling (default: `@every 1m`)
3. **Running the Application:**
   - In your terminal, navigate to the project directory and execute:
     ```bash
     go run main.go
     ```

## Usage

The application runs as a background service and logs the following information:

- **Summary Statistics:** Periodically logs the total number of users in the collection.
- **Change Logs:** Logs details of new users whenever a change in the user count is detected.

## Technical Details

- **Programming Language:** Go
- **Dependencies:** See the `go.mod` file for a list of dependencies.

## Contributing

- Feel free to create issues for feature requests or bug reports.
- Pull requests are welcome for contributions.

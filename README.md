
# Forum

## Objectives
This project aims to create a web forum with the following features:
- User communication via posts and comments.
- Post categorization using tags or categories.
- Ability to like or dislike posts and comments.
- Post filtering by categories or other criteria.

## Installation

1. Clone the repository:
   ```bash
   git clone git@github.com:azhaxylyk/forum.git
   ```

2. Navigate to the project directory:
   ```bash
   cd forum
   ```

## Usage

### Run the Program
You can start the program in one of the following ways:

1. **Directly using Go**:
   ```bash
   go run ./cmd
   ```

2. **Using Makefile**:
   - Start the server:
     ```bash
     make run
     ```

3. **Using Docker**:
   - Build the Docker image:
     ```bash
     make build
     ```
   - Run the Docker container:
     ```bash
     make run-img
     ```

### Stop the Server
To stop the server or Docker container:
```bash
make stop
```

## Authors
- [azhaxylyk](https://01.alem.school/git/azhaxylyk)
- [abulatov](https://01.alem.school/git/abulatov)
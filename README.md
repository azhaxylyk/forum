# Forum
### Objectives
This project is a simple forum application where users can create posts, comment on them, like/dislike posts and comments, and manage notifications. It includes features for admins and moderators to manage content, approve/reject user actions, and enforce moderation.

## Installation
Clone the repository:
```bash
git clone git@github.com:azhaxylyk/forum.git
```
Navigate to the project directory:
```bash
cd forum
```

## Running the Project
#### Initialize the Environment
Before running the project, you need to generate SSL certificates and set up environment variables.
```
make init
```
This command will:
- Create the certs/ directory if it does not exist.
- Generate self-signed SSL certificates (certs/localhost.crt and certs/localhost.key).
- Set environment variables SSL_CERT_PATH and SSL_KEY_PATH for secure communication.
#### If you want the environment variables to persist across sessions, run:
```
source ~/.bashrc
```

## Usage
#### Run the Program
You can start the program in one of the following ways:

**1. Directly using Go:**
   ```bash
   go run ./cmd
   ```
**2. Using Makefile:**

Start the server:
```bash
make run
   ```
**3. Using Docker:**

Build the Docker image:
```bash
make build
```
Run the Docker container:
```bash
make run-img
```
#### Stop the Server
To stop the server or Docker container:
   ```bash
   make stop
   ```

## Admin Login Credentials
The default admin login credentials are:
- Email: admin@gmail.com
- Password: Aa123456

## Authors
- [azhaxylyk](https://github.com/azhaxylyk)
- [abulatov](https://github.com/Alish98b)
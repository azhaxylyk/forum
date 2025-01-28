# Forum
## Objectives
This project is a simple forum application where users can create posts, comment on them, like/dislike posts and comments, and manage notifications. It includes features for admins and moderators to manage content, approve/reject user actions, and enforce moderation.

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
   Start the server:
   ```bash
   make run
   ```
3. **Using Docker**:
   Build the Docker image:
   ```bash
   make build
   ```
   Run the Docker container:
   ```bash
   make run-img
   ```
### Stop the Server
To stop the server or Docker container:
   ```bash
   make stop
   ```

## Admin Login Credentials
The default admin login credentials are:
- Email: admin@gmail.com
- Password: Aa123456

## Configuration
Make sure to set the following environment variables for proper functionality:
- SSL_CERT_PATH: Path to your SSL certificate.
- SSL_KEY_PATH: Path to your SSL key.
For example, you can set the environment variables like this:
```bash
export SSL_CERT_PATH="certs/localhost.crt"
export SSL_KEY_PATH="certs/localhost.key"
```
These are necessary for the server to run with HTTPS.

## Authors
- [azhaxylyk](https://github.com/azhaxylyk)
- [abulatov](https://github.com/Alish98b)
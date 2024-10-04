# Shareless

Shareless is an application developed in Go that allows for secure secret sharing. With a focus on privacy and security, the application does not store secrets in plain text; instead, it only stores a hash that enables decryption of the secret when needed.

## Features

- Share secrets with configurable expiration.
- Use temporary tokens for accessing secrets.
- Efficient storage with Redis.

## Requirements

- Go (version 1.16 or higher)
- Redis

## Endpoints

### 1. POST /share

Used to share a new secret.

#### Body
```json
{
  "text": "The secret you want to share",
  "expire_on_opened": "on", // (optional) Indicates whether the secret should expire after being accessed.
  "duration": "1h" // Duration before the secret expires.
}
```

#### Response
```json
{
  "url": "http://host/shared/{id}?token={token}"
}
```

### 2. GET /shared/{id}

Used to access a shared secret.

#### Query Params

- token: string - The token generated at the time of sharing.

#### Response

```json
{
  "text": "plain text of the secret decrypted from the token"
}
```

### How to Run

- Clone the repository:
```bash
git clone https://github.com/omurilo/shareless.git
cd shareless
```

- Install dependencies:
```bash
go mod tidy
```

- Start the server:
```bash
make run # make watch // to run in watching mode
```

- Run in docker-compose (redis included):
```bash
make docker-run
```

- Access the application at `http://localhost:3000`.

### Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

```
Feel free to make any adjustments as needed!
```

### Inspiration
This project is inspired by https://passshare.me, a project of diogo.dev shared on bluesky.

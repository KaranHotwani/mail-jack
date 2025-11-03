# Mail Jack

Lightweight Go-based HTTP service to send emails via providers (currently **AWS SES**).  
It offers **per-recipient tracking**, **synchronous responses**, and **structured JSON errors** — all through a single endpoint.

---

## ✨ Features

- **AWS SES provider** (more coming soon)
- **Synchronous API:** returns actual send status per recipient
- **Detailed logging** — logs every request, response, and error in PostgreSQL for easy debugging and analytics 
- **Docker support** for easy deployment
- **Open source** and self-hostable anywhere

---

## 🧩 Requirements

- Go **1.22+** (module sets **1.24.x** toolchain)
- AWS SES configured (verified sender/domain)
- PostgreSQL (for email logs)

---

## ⚙️ Environment Variables

| Variable | Description | Default |
|-----------|--------------|----------|
| `PORT` | HTTP port | `8080` |
| `EMAIL_PROVIDER` | Email provider name (`SES`) | — |
| `AWS_REGION` | AWS region (e.g. `us-east-1`) | — |
| `AWS_ACCESS_KEY_ID` | AWS Access Key | — |
| `AWS_SECRET_ACCESS_KEY` | AWS Secret Key | — |
| `AWS_SESSION_TOKEN` | Optional session token | — |
| `MAIL_JACK_API_KEY` | API key for authentication | **required** |
| `DATABASE_URL` | PostgreSQL connection string | **required** |

Format:
```
postgres://username:password@host:port/database?sslmode=require
```

---

## 🧱 Run Locally

```bash
go run ./cmd
```

---

## 🐳 Docker

Create a `.env` file with all variables (.sample.env available), then:

```bash
# Build the image
docker build -t mail-jack:latest .

# Run the container
docker run -d -p 8080:8080 --env-file .env mail-jack:latest
```

---

## 📬 API

### POST `/send-email`

#### Headers
```text
Content-Type: application/json
X-API-KEY: your_api_key
```

#### Body (JSON)
```json
{
  "from": "sender@example.com",
  "to": ["user1@example.com", "user2@example.com"],
  "subject": "Hello",
  "body": "Plain text body",
  "html": "<p>HTML body</p>",
  "ccEmails": ["cc1@example.com"]
}
```

#### ✅ Success Response
```json
{
  "status": "SUCCESS",
  "results": [
    {
      "email": "user1@example.com",
      "status": "SUCCESS",
      "messageId": "010e0199b4711bc0-...",
      "error": ""
    }
  ]
}
```

#### ❌ Error Responses
```json
400 { "error": "Invalid request" }
401 { "error": "Invalid API key" }
500 { "error": "Internal server error" }
```

#### Example (curl)
```bash
curl -X POST http://localhost:8080/send-email   -H "Content-Type: application/json"   -H "X-API-KEY: your_secret_api_key"   -d '{"from":"sender@example.com","to":["user@example.com"],"subject":"Hi","body":"Hello","html":"<p>Hello</p>"}'
```

---

## 🧪 Build & Run

```bash
# Build binary
go build -o mail-jack ./cmd

# Run
./mail-jack
```

---

## 🧭 Roadmap

- [x] Basic HTTP server  
- [x] Send Email API (POST /send-email)  
- [x] Split SES calls per recipient  
- [x] PostgreSQL email logs  
- [x] Concurrent sending via goroutines  
- [ ] Async sending via SQS  
- [ ] Web UI for logs  
- [ ] Webhooks for status updates  
- [ ] Add rate limiter
- [ ] Retry and Failure mechanism  
 
---

---

## 📜 License

MIT © Karan Hotwani

**Connect with me:**  
[LinkedIn](https://www.linkedin.com/in/karan-hotwani-a9ba73167/) • [Twitter](https://x.com/Karan151997)



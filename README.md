Mail Jack
=========

Lightweight HTTP service to send emails via providers (currently AWS SES). Includes per-recipient status reporting and JSON error responses.

Features
- AWS SES provider
- Synchronous API: returns actual send status
- Per‑recipient results (success/failed, messageId, error)
- JSON errors with proper HTTP status codes

Requirements
- Go 1.22+ (module sets 1.24.x toolchain)
- AWS SES configured (verified sender/domain)

Environment variables
- PORT: HTTP port (default: 8080)
- EMAIL_PROVIDER: provider name (SES)
- AWS_REGION: AWS region (e.g. us-east-1)
- AWS credentials: via IAM role or env (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, optional `AWS_SESSION_TOKEN`)

Run locally
```bash
go run ./cmd
# or
PORT=8080 EMAIL_PROVIDER=SES AWS_REGION=us-east-1 go run ./cmd
```

API
- POST /send-email

Request body
```json
{
  "from": "sender@example.com",
  "to": ["user1@example.com", "user2@example.com"],
  "subject": "Hello",
  "body": "Plain text body",
  "html": "<p>HTML body</p>"
}
```

Success response (per‑recipient)
```json
{
  "status": "SUCCESS"| "PARTIAL_SUCCESS" | "FAILED",
  "results": [
    {
      "email": "user1@example.com",
      "status": "SUCCESS",
      "messageId": "010e0199b4711bc0-459734c9-418f-4f9d-982a-dbe86dd1f3f5-000000",
      "error": ""
    },
    {
      "email": "user2@example.com",
      "status": "SUCCESS",
      "messageId": "010e0199b4711bc1-8590683f-4390-4d2a-a11b-698a29569699-000000",
      "error": ""
    }
  ]
}
```

Error responses (JSON)
- 400: `{ "error": "Invalid request" }`
- 500: `{ "error": "..." }`

Build
```bash
go build -o mail-jack ./cmd
```

Deploy (Linux binary)
```bash
PORT=8080 EMAIL_PROVIDER=SES AWS_REGION=us-east-1 ./mail-jack
```

Notes
- Imports/module path should match your repo: `module github.com/KaranHotwani/mail-jack`
- If exposing to end users via npm, publish a small JS/TS client that POSTs to `/send-email`.

License
MIT

Roadmap
-------
- [x] Basic HTTP server
- [x] Send email API endpoint (`POST /send-email`)
- [x] Split email sending to individual SES calls per recipient
- Create an npm package to expose a client SDK for triggering the Go API
- Auto-start Go server when the npm package send method is called, if not already running
- [x] Add PostgreSQL to store email logs (requests, responses, errors)
- Test cases — unit and integration tests for all critical functionality.
- [x] Add goroutines to optimize email sending (non-blocking/concurrent)
- Host Go server as a Lambda function triggered via API Gateway for serverless email sending
- Make email sending fully asynchronous using SQS (queue) and SNS (optional notifications)
- Build a simple UI to track email logs in PostgreSQL
- Add webhooks feature to notify external systems on email success/failure
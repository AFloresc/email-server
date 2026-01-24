# 📬 Email Server (Go + Resend + Render)

Servicio backend minimalista escrito en Go, diseñado para recibir mensajes desde el formulario de tu portfolio y enviarlos mediante Resend API.
Optimizado para desplegarse en Render con un servidor ligero y seguro.

---

## 🚀 Características

- Endpoint POST /contact para recibir mensajes JSON.
- Validación básica de campos (`name`, `email`, `message`).
- Envío de correos mediante Resend (HTML + texto plano).
- Soporte para Reply‑To dinámico.
- Middleware CORS para permitir peticiones desde tu frontend.
- Logs claros para depuración en Render.
- Código simple, modular y mantenible.

---

## 📁 Estructura del proyecto

```text
contact-service/
.
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── router/
│   │   └── router.go
│   │
│   ├── handlers/
│   │   └── contact.go
│   │
│   ├── email/
│   │   └── resend.go
│   │
│   ├── security/
│   │   ├── validate.go
│   │   ├── sanitize.go
│   │   ├── patterns.go
│   │   └── ip_hash.go
│   │
│   ├── middleware/
│   │   ├── rate_limit.go
│   │   ├── cooldown.go
│   │   ├── user_agent.go
│   │   ├── body_limit.go
│   │   └── metrics.go
│   │
│   ├── logging/
│   │   └── logger.go
│   │
│   └── metrics/
│       └── metrics.go
│
├── go.mod
├── go.sum
└── Dockerfile
```

---

## 🔧 Requisitos previos

- Go 1.22+
- Cuenta en Resend (API Key)
- Docker
- Cuenta en Render (para el deploy)
- Variables de entorno configuradas

```text
RESEND_API_KEY=tu_api_key
TO_EMAIL=tu_correo_destino
ALLOWED_ORIGINS=http://localhost:5173
PORT=8080 (Render lo inyecta automáticamente)

```

---

## ⚙️ Configuración de entorno en Render

En Render → Dashboard → Environment Variables, añade:

```text
RESEND_API_KEY
TO_EMAIL
ALLOWED_ORIGINS
- (Render añade PORT automáticamente)

```

## 📡 API

### `POST /contact`

Envía un mensaje desde el formulario.

#### Body (JSON)

```json
{
  "name": "Alex",
  "email": "alex@example.com",
  "message": "Hola, quiero contactar contigo."
}

Respuesta exitosa
{
  "status": "ok"
}

Respuesta de error
{
  "error": "Failed to send email"
}
```

## 📨 Envío de correo

El backend envía:

- HTML premium con tu branding
- Texto plano para compatibilidad
- Reply‑To con el email del usuario
- Logs estructurados del resultado de Resend

## 📊 Métricas internas

Disponible en:

### `GET /metric`

ncluye:

- total_requests
- successful_requests
- validation_errors
- honeypot_blocks
- firewall_blocks
- rate_limit_blocks
- cooldown_blocks
- user_agent_blocks
- payload_too_large
- email_send_errors
- total_processing_time_ms

## 🛡️ CORS

Permitido para:

```text
http://localhost:5173
```

Puedes añadir más orígenes según despliegues.

## ▶️ Ejecutar en local

```text
go run ./cmd/server
```

## 🐳 Docker

El proyecto incluye un Dockerfile multi‑stage optimizado:

```text
docker build -t email-server .
docker run -p 8080:8080 email-server
```

## 🚀 Deploy en Render

- Crear nuevo servicio → Web Service
- Seleccionar tu repo
- Runtime: Docker
- Render detecta el Dockerfile automáticamente
- Añadir variables de entorno
- Deploy

```text
    go build -o server .
```

- Start Command:

```text
    ./server
```

- Añadir variables de entorno
- Deploy

## 📜 Licencia

MIT — libre para usar y modificar.

---

## 🛡️ Protecciones incluidas

✔ SQL injection

```text
  SELECT * FROM users
  ' OR '1'='1
  DROP TABLE
```

✔ XSS

```text
  <script>alert(1)</script>
  onerror=alert(1)
```

✔ Directory traversal

```text
  ../../etc/passwd
```

✔ Command injection

```text
  rm -rf /
  $(curl ...)
```

✔ Payloads codificados

```text
  %3Cscript%3E
```

✔ Bots que envían patrones típicos de ataque

---

## 🔄 Diagrama de flujo del request (POST /contact)

```text
┌──────────────────────────────┐
│        Cliente (Frontend)    │
└───────────────┬──────────────┘
                │  POST /contact
                ▼
        ┌───────────────────────┐
        │       Router          │
        └───────────┬──────────┘
                    │
                    ▼
        ┌───────────────────────┐
        │   Middleware: Metrics │
        └───────────┬──────────┘
                    │
                    ▼
        ┌───────────────────────┐
        │ Middleware: RateLimit │───┐
        └───────────┬──────────┘   │ Too many requests
                    │               │ → 429
                    ▼               │
        ┌───────────────────────┐   │
        │ Middleware: Cooldown  │───┘
        └───────────┬──────────┘
                    │
                    ▼
        ┌──────────────────────────────┐
        │ Middleware: User-Agent Check │───┐
        └───────────┬──────────────────┘   │ UA sospechoso
                    │                      │ → 403
                    ▼                      │
        ┌──────────────────────────────┐   │
        │ Middleware: Body Size Limit  │───┘
        └───────────┬──────────────────┘
                    │
                    ▼
        ┌──────────────────────────────┐
        │     Handler: /contact        │
        └───────────┬──────────────────┘
                    │
                    ▼
        ┌──────────────────────────────┐
        │   Decodificar JSON (100KB)   │───┐
        └───────────┬──────────────────┘   │ JSON inválido
                    │                      │ → 400
                    ▼                      │
        ┌──────────────────────────────┐   │
        │  security.ValidateContact    │───┘
        ├──────────────────────────────┤
        │ Honeypot → 403               │
        │ Firewall → 403               │
        │ Email inválido → 400         │
        │ Mensaje inválido → 400       │
        └───────────┬──────────────────┘
                    │
                    ▼
        ┌──────────────────────────────┐
        │  email.SendEmailResend       │───┐
        └───────────┬──────────────────┘   │ Error enviando email
                    │                      │ → 500
                    ▼                      │
        ┌──────────────────────────────┐   │
        │   Resend API (externo)       │   │
        └───────────┬──────────────────┘   │
                    │                      │
                    ▼                      │
        ┌──────────────────────────────┐   │
        │   Respuesta exitosa          │   │
        └───────────┬──────────────────┘   │
                    │
                    ▼
        ┌──────────────────────────────┐
        │     { "status": "ok" }       │
        └──────────────────────────────┘
```

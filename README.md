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
├── main.go
├── go.mod
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
PORT=8080 (Render lo inyecta automáticamente)

```

---

## ⚙️ Configuración de entorno en Render

En Render → Dashboard → Environment Variables, añade:

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
- Logs del resultado de Resend

## 🛡️ CORS

Permitido para:

```text
http://localhost:5173
```

Puedes añadir más orígenes según despliegues.

## ▶️ Ejecutar en local

```text
go run main.go
```

## 🚀 Deploy en Render

- Crear nuevo servicio → Web Service
- Seleccionar tu repo
- Runtime: Go
- Build Command:

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

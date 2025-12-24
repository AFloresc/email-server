# 📬 Contact Service (Go + SMTP + Fly.io)

Servicio backend minimalista escrito en **Go**, diseñado para recibir mensajes desde un formulario web y enviarlos por correo electrónico utilizando **SMTP**.  
Optimizado para desplegarse en **Fly.io** con un contenedor ligero basado en Alpine.

---

## 🚀 Características

- Endpoint `POST /contact` para recibir mensajes JSON.
- Validación básica de campos (`name`, `email`, `message`).
- Envío de correos mediante **SMTP** (Gmail, Outlook, Mailgun, etc.).
- Configuración segura mediante **Fly.io Secrets**.
- Dockerfile optimizado para despliegues rápidos.
- Código simple, mantenible y sin dependencias externas.

---

## 📁 Estructura del proyecto

contact-service/
│
├── main.go
├── go.mod
├── Dockerfile
└── fly.tom

---

## 🔧 Requisitos previos

- Go 1.22+
- Docker
- Cuenta en Fly.io (`flyctl` instalado)
- Credenciales SMTP válidas  
  (Gmail requiere **App Password**, no la contraseña normal)

---

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

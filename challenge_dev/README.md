# Portal de Convocatorias Públicas

Repositorio del portal web para consultar convocatorias públicas, registrar usuarios, autenticar sesiones con JWT y administrar propuestas guardadas.

Arquitectura:

- `backend/`: API HTTP en Go.
- `frontend/`: cliente web estático en HTML/CSS/JavaScript.
- `database/`: script SQL de creación de esquema.

## 1. Resumen técnico

La solución está compuesta por dos capas principales:

- Backend en Go con rutas HTTP protegidas por JWT.
- Frontend estático servido con `python3 -m http.server`.

El backend integra:

- autenticación y sesión de usuario;
- consulta de convocatorias desde una fuente externa de datos abiertos;
- persistencia local de propuestas y asociaciones;
- cierre de sesión mediante invalidación de token;
- validación de esquema de base de datos al arrancar.

La interfaz está en español y usa la paleta:

- `#4A4466` principal
- `#6EADBC` secundario
- `#9FCBAD` acento
- `#F1F7D4` fondo/acento claro

## 2. Requisitos previos

Instalar o disponer de:

- **Go** 1.22 o superior
- **Python 3**
- **PostgreSQL** 14 o superior
- **curl** para pruebas de verificación

Verificación de herramientas:

```bash
go version
python3 --version
psql --version
```

## 3. Ejecución local

### 3.1 Backend

Desde la carpeta `backend/`:

```bash
cd backend
PORT=8888 go run ./cmd/api
```

El backend expone la API en:

- `http://localhost:8888`

Durante el arranque:

- se intenta validar el esquema de PostgreSQL;
- si el esquema existe, la migración se omite;
- si PostgreSQL no está disponible, el servidor continúa con repositorios en memoria para desarrollo local.

### 3.2 Frontend

Desde la carpeta `frontend/`:

```bash
cd frontend
python3 -m http.server 3000
```

La interfaz queda disponible en:

- `http://localhost:3000/pages/`

El archivo `frontend/pages/index.html` redirige automáticamente a `dashboard.html` para evitar el listado de directorios.

## 4. Variables de entorno

El backend consume las siguientes variables:

| Variable | Default | Uso |
|----------|---------|-----|
| `DB_USER` | `sebasdb` | Usuario de PostgreSQL |
| `DB_PASSWORD` | `challenge_app_07032026` | Contraseña de PostgreSQL |
| `DB_HOST` | `localhost` | Host de la base de datos |
| `DB_PORT` | `5432` | Puerto de PostgreSQL |
| `DB_NAME` | `portal_plan_public_app` | Nombre de la base de datos |
| `PORT` | `8080` | Puerto HTTP del backend |
| `SECRET_KEY` | `challenge-secret-key-2026` | Clave para firmar JWT |
| `INTEGRATION_URL` | `https://datos.gov.co/resource/p6dx-8zbt.json` | Fuente externa de convocatorias |

Ejemplo de configuración:

```bash
export DB_USER=sebasdb
export DB_PASSWORD=mi_clave
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=portal_plan_public_app
export PORT=8888
export SECRET_KEY=mi_secret_key
export INTEGRATION_URL=https://datos.gov.co/resource/p6dx-8zbt.json
```

## 5. Rutas y páginas del frontend

### Páginas disponibles

| Página | Ruta | Propósito |
|--------|------|-----------|
| `index.html` | `/pages/` | Redirección por defecto al dashboard |
| `login.html` | `/pages/login.html` | Autenticación de usuario |
| `user-create.html` | `/pages/user-create.html` | Registro de usuario |
| `dashboard.html` | `/pages/dashboard.html` | Listado principal de convocatorias |
| `saved-proposals.html` | `/pages/saved-proposals.html` | Convocatorias guardadas |

### Flujo operativo

1. Acceder a `/pages/` o abrir `login.html`.
2. Autenticarse.
3. Consultar convocatorias desde el dashboard.
4. Guardar una propuesta.
5. Revisar propuestas guardadas.
6. Cerrar sesión con el botón `Cerrar sesión`.

## 6. Estructura general del repositorio

```text
challenge_dev/
├── backend/
│   ├── cmd/api/
│   └── internal/
├── database/
│   └── public_calls_database.sql
├── frontend/
│   ├── pages/
│   ├── scripts/
│   └── sheets/
└── specs_project/
```

## 7. Endpoints principales

| Método | Ruta | Auth | Descripción |
|--------|------|:----:|-------------|
| GET | `/health` | No | Health check |
| POST | `/user-create` | No | Registro de usuario |
| POST | `/login` | No | Inicio de sesión |
| POST | `/logout` | Sí | Invalida la sesión actual |
| POST | `/user-info` | Sí | Perfil del usuario autenticado |
| POST | `/user-modify` | Sí | Actualización de perfil |
| GET | `/public-proposals` | Sí | Convocatorias desde datos abiertos |
| GET | `/saved_proposals` | Sí | Propuestas guardadas con datos completos |
| POST | `/saved-proposals` | Sí | Asociación de propuesta al usuario |
| POST | `/proposal-save` | Sí | Upsert de propuesta completa |

## 8. Verificación rápida

```bash
curl http://localhost:8888/health
```

Respuesta esperada:

```json
{"status":"ok"}
```

## 9. Notas técnicas

- La interfaz es completamente estática.
- El backend usa JWT para proteger las rutas privadas.
- El logout invalida el token actual del usuario.
- El arranque del backend intenta validar el esquema en PostgreSQL antes de servir tráfico.
- El repositorio soporta ejecución local incluso si PostgreSQL no está disponible, mediante mocks en memoria.


# SECOP II SODA API - datos.gov.co Documentation

## Overview

- **Portal:** https://datos.gov.co
- **Dataset:** SECOP II - Procesos de Contratación
- **Dataset ID:** `p6dx-8zbt`
- **Protocol:** Socrata Open Data API (SODA) v2.1
- **Authentication:** None required (public data)
- **Base URL:** `https://datos.gov.co/resource/p6dx-8zbt`
- **Format:** JSON (standard), also CSV and XML supported
- **Owner:** Agencia Nacional de Contratación Pública - Colombia Compra Eficiente

## Endpoint

### GET /resource/p6dx-8zbt.json

Query procurement processes (convocatorias públicas) recorded in the SECOP II platform since launch.

```
GET https://datos.gov.co/resource/p6dx-8zbt.json
```

## Input Parameters (SODA Query Language)

| Parameter | Type | Description | Example |
|---|---|---|---|
| `$limit` | integer | Maximum rows to return | `$limit=10` |
| `$offset` | integer | Row offset for pagination | `$offset=0` |
| `$select` | string | Comma-separated columns to return | `$select=entidad,nombre_del_procedimiento,fase` |
| `$where` | string | SoQL WHERE clause for filtering | `$where=fase = 'Presentación de oferta'` |
| `$order` | string | Sort by column (ASC/DESC) | `$order=fecha_de_publicacion_del DESC` |
| `$group` | string | Group by column | `$group=fase` |
| `$q` | string | Full-text search | `$q=convocatoria` |

### Filtering by Convocatorias Públicas

Use SoQL text operators on the `nombre_del_procedimiento` or `descripci_n_del_procedimiento` fields:

```
$where=lower(nombre_del_procedimiento) like '%convocatoria%'
```

## Available Columns (59 total)

Key columns for querying public calls:

| Field Name | Data Type | Description |
|---|---|---|
| `entidad` | text | Entity publishing the procurement process |
| `nit_entidad` | text | Tax ID of the publishing entity |
| `departamento_entidad` | text | Department of the entity |
| `ciudad_entidad` | text | City of the entity |
| `ordenentidad` | text | Entity order (Nacional, Regional) |
| `id_del_proceso` | text | Unique process ID from platform |
| `referencia_del_proceso` | text | Process reference from entity |
| `nombre_del_procedimiento` | text | Name of the procurement procedure |
| `descripci_n_del_procedimiento` | text | Description of the procedure |
| `fase` | text | Current phase of the process |
| `fecha_de_publicacion_del` | calendar_date | Initial publication date |
| `fecha_de_ultima_publicaci` | calendar_date | Last publication date |
| `modalidad_de_contratacion` | text | Selection modality |
| `precio_base` | number | Base estimated price |
| `duracion` | number | Estimated duration |
| `unidad_de_duracion` | text | Duration unit |
| `fecha_de_recepcion_de` | calendar_date | Response submission deadline |
| `estado_del_procedimiento` | text | Current status |
| `adjudicado` | text | Whether awarded (Si/No) |
| `nombre_del_proveedor` | text | Awarded supplier name |
| `valor_total_adjudicacion` | number | Total award value |
| `urlproceso` | url | URL to the process on SECOP II |
| `codigo_principal_de_categoria` | text | UNSPSC main category code |
| `tipo_de_contrato` | text | Contract type |
| `estado_de_apertura_del_proceso` | text | Information opening status |
| `estado_resumen` | text | Summary status |
| `proveedores_invitados` | number | Total invited suppliers |
| `proveedores_que_manifestaron` | number | Suppliers who expressed interest |
| `respuestas_al_procedimiento` | number | Total responses |
| `numero_de_lotes` | number | Number of item lots |

## Output Format

The API returns a JSON array of objects. Each object represents one procurement process.

### Example Response (single record)

```json
{
  "entidad": "DEPARTAMENTO ADMINISTRATIVO NACIONAL DE ESTADISTICA (DANE)",
  "nit_entidad": "899999027",
  "departamento_entidad": "Distrito Capital de Bogotá",
  "ciudad_entidad": "Bogotá",
  "ordenentidad": "Nacional",
  "id_del_proceso": "CO1.REQ.2577563",
  "referencia_del_proceso": "EDP-545-2022",
  "nombre_del_procedimiento": "Convocatoria DSNFT-0001-FEEC-2024",
  "descripci_n_del_procedimiento": "Aunar esfuerzos entre las partes para ejecutar el proyecto...",
  "fase": "Presentación de oferta",
  "fecha_de_publicacion_del": "2024-02-19T00:00:00.000",
  "fecha_de_ultima_publicaci": "2024-03-12T00:00:00.000",
  "modalidad_de_contratacion": "Licitación pública",
  "precio_base": 500000000.00,
  "estado_del_procedimiento": "Publicado",
  "urlproceso": "https://community.secop.gov.co/..."
}
```

All `calendar_date` values use ISO 8601 format: `YYYY-MM-DDTHH:mm:ss.SSS`
All `number` values use decimal format without currency symbol.
All `text` values are UTF-8 strings.

## Example Queries

### 1. Get latest public calls
```
GET https://datos.gov.co/resource/p6dx-8zbt.json
  ?$limit=10
  &$order=fecha_de_publicacion_del DESC
```

### 2. Search convocatorias by keyword
```
GET https://datos.gov.co/resource/p6dx-8zbt.json
  ?$limit=20
  &$select=entidad,nombre_del_procedimiento,descripci_n_del_procedimiento,fase,fecha_de_publicacion_del
  &$where=lower(nombre_del_procedimiento) like '%convocatoria%'
  &$order=fecha_de_publicacion_del DESC
```

### 3. Filter by phase
```
GET https://datos.gov.co/resource/p6dx-8zbt.json
  ?$limit=50
  &$where=fase = 'Presentación de oferta'
```

### 4. Get processes by entity
```
GET https://datos.gov.co/resource/p6dx-8zbt.json
  ?$limit=10
  &$where=lower(entidad) like '%dane%'
```

### 5. Get processes by modality (licitación pública)
```
GET https://datos.gov.co/resource/p6dx-8zbt.json
  ?$limit=10
  &$where=lower(modalidad_de_contratacion) like '%licitaci%n%'
```

### 6. Paginate results
```
GET https://datos.gov.co/resource/p6dx-8zbt.json
  ?$limit=100
  &$offset=0
  &$order=fecha_de_publicacion_del DESC
```

### 7. Count/aggregate by phase
```
GET https://datos.gov.co/resource/p6dx-8zbt.json
  ?$select=fase,count(*)
  &$group=fase
```

### 8. Processes published in 2024 with open phase
```
GET https://datos.gov.co/resource/p6dx-8zbt.json
  ?$limit=100
  &$where=fecha_de_publicacion_del >= '2024-01-01T00:00:00.000' AND fase in('Presentación de oferta','Manifestación de interés')
  &$order=fecha_de_publicacion_del DESC
```

## Limitations

- No API key required, but rate limits apply (undocumented; reasonable usage recommended).
- Maximum 50,000 rows per query without a custom application token.
- Pagination required for large result sets (use `$offset` with `$limit`).
- Calendar dates are returned in Colombia timezone.
- Some special characters in column names (accents, ñ) are preserved.

## Implementation Notes for Go

- Use `net/http` standard library for HTTP GET requests.
- Parse JSON response into a struct slice.
- URL-encode query parameters with `url.Values`.
- Calendar dates should be parsed as `time.Time` from RFC3339.
- Numbers should be parsed as `float64` or `json.Number`.
- See `test.go` for a complete working example.
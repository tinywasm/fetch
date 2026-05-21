# PLAN — Constructor público `NewResponse` para uso desde Service Workers

> Auto-contenido. Cambio aditivo (no breaking).

## Problema

`fetch.Response` se construye hoy únicamente de forma interna por `doRequest`
(la respuesta a una llamada saliente). Su campo `body` es privado: sólo hay
getters (`Body()`, `Text()`, `GetHeader()`), ningún constructor ni setter
público.

Un Service Worker handler (`tinywasm/js`) necesita **construir** una
`fetch.Response` para devolverla desde `OnFetch(ctx, req) (*fetch.Response, error)`.
Hoy es imposible fuera del paquete `fetch`.

## Solución — añadir un constructor público

```go
// NewResponse construye una Response servible (p.ej. desde un Service Worker
// handler que responde un FetchEvent interceptado).
//
//   status  — código HTTP (200, 404, ...).
//   headers — cabeceras de respuesta ({Key, Value}).
//   body    — cuerpo crudo ya serializado.
func NewResponse(status int, headers []Header, body []byte) *Response {
    return &Response{
        Status:  status,
        Headers: headers,
        body:    body,
    }
}
```

Ubicación sugerida: junto a la definición de `Response` en
[fetch.go](../fetch.go). `RequestURL` y `Method` quedan en cero — no aplican
a una respuesta construida localmente; si un caso futuro los necesita, se
añade una variante o se setean por campo público (ya son exportados).

## Justificación

- **Cambio aditivo, no breaking:** sólo agrega una función; no toca firmas
  ni comportamiento existente.
- **Mantiene `body` privado:** se preserva la invariante de que el cuerpo
  sólo se lee vía `Body()`/`Text()`. El constructor es el único punto de
  escritura público.
- **Evita que `tinywasm/js` defina un `Response` propio:** reusar el tipo
  canónico del ecosistema mantiene una sola representación de respuesta HTTP.

## Tests

| Archivo | Test | Verifica |
|---|---|---|
| `tests/response_test.go` (nuevo o anexo) | `TestNewResponse_RoundTrip` | `NewResponse(200, hdrs, body)` → `Status==200`, `Body()` devuelve `body`, `GetHeader` resuelve las cabeceras pasadas |
| idem | `TestNewResponse_EmptyBody` | `NewResponse(204, nil, nil)` no panica; `Body()` devuelve slice vacío/nil |

Ejecución: `gotest ./...` (skill `testing`). El paquete tiene tests dual
stdlib/wasm — verificar que `NewResponse` (lógica pura, sin `syscall/js`)
pasa en ambos modos.

## Reglas de dependencias

`tinywasm/fetch` ya cumple la regla del ecosistema (usa `tinywasm/fmt`, sin
stdlib vetada). `NewResponse` no introduce imports nuevos.

## Stages

| # | Tarea | Done |
|---|---|---|
| 1 | Añadir `func NewResponse(status int, headers []Header, body []byte) *Response` en `fetch.go` | [x] |
| 2 | Tests `TestNewResponse_RoundTrip` y `TestNewResponse_EmptyBody` | [x] |
| 3 | `gotest ./...` verde (stdlib + wasm) | [x] |
| 4 | Actualizar `docs/API.md` documentando el constructor | [x] |
| 5 | Publicar nueva versión con `gopush` | [x] |

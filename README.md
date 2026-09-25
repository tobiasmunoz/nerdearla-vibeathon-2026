# 🎙️ Nerdearla Live Subs (Go Edition)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go: 1.23+](https://img.shields.io/badge/Go-1.23%2B-00ADD8.svg)](https://go.dev/)
[![Gemini Live API](https://img.shields.io/badge/Gemini%20Live%20API-Bidirectional%20Streaming-4285F4.svg)](https://ai.google.dev/)
[![Hackathon: Nerdearla Vibeathon](https://img.shields.io/badge/Hackathon-Nerdearla%20Vibeathon%202026-orange.svg)](https://nerdear.la)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)

> **Open-source, ultra-low-latency simultaneous live transcription and translation system built for large-scale multi-stage tech conferences like [Nerdearla](https://nerdear.la).**  
> Streams live speech directly to the **Google Gemini Live API** and broadcasts real-time Spanish & English subtitles to attendee screens and OBS Studio broadcast overlays.

---

## 📑 Tabla de Contenidos / Table of Contents

- [🌟 Contexto & Misión / Context & Mission](#-contexto--misión)
- [✨ Características Principales / Key Highlights](#-características-principales)
- [🧩 Requisitos Mínimos (MVP Checklist)](#-requisitos-mínimos-mvp-checklist)
- [🍬 Opcionales Implementados (Bonus Features)](#-opcionales-implementados)
- [📊 Benchmarks: Go vs Python Architecture](#-benchmarks-go-vs-python-architecture)
- [🏗️ Arquitectura del Sistema / Architecture](#️-arquitectura-del-sistema)
- [🚀 Guía de Inicio Rápido / Quick Start](#-guía-de-inicio-rápido)
  - [1. Configuración de Credenciales](#1-configuración-de-credenciales)
  - [2. Correr Localmente (Go)](#2-correr-localmente)
  - [3. Probar con Audio de Muestra (Sin Micrófono)](#3-probar-con-audio-de-muestra)
- [🎥 Integración con OBS Studio & vMix](#-integración-con-obs-studio--vmix)
- [👥 Vista para la Audiencia (Audience Mode)](#-vista-para-la-audiencia)
- [🔀 Escalabilidad Multi-Escenario (10+ Escenarios)](#-escalabilidad-multi-escenario)
- [☁️ Despliegue en Google Cloud Run](#️-despliegue-en-google-cloud-run)
- [📄 Licencia Open Source](#-licencia)

---

## 🌟 Contexto & Misión

En conferencias internacionales como **Nerdearla**, decenas de charlas simultáneas se dictan en inglés. Para maximizar la accesibilidad de la comunidad hispanohablante, se requiere traducción e interpretación simultánea.

Las soluciones tradicionales tienen barreras infranqueables de escala:
1. **Intérpretes humanos simultáneos**: Muy costosos y difíciles de coordinar en conferencias con 5, 10 o 30 escenarios en paralelo.
2. **Pipelines ASR tradicionales en cascada (Whisper batch + Traducción text-to-text)**: Tienen latencias de 4 a 8 segundos, rompiendo la sincronía entre el disertante y la pantalla.

**Nerdearla Live Subs (Go Edition)** resuelve esto mediante un motor en **Golang de alto rendimiento** conectado de forma bidireccional por WebSocket a la **Google Gemini Live API**:
- La voz se captura y procesa en chunks de 16ms (PCM 16kHz).
- Dos sesiones concurrentes de Gemini generan en paralelo:
  1. **Transcripción en idioma original** (`gemini-3.5-transcribe-live`).
  2. **Traducción simultánea en español** (`gemini-3.5-live-translate-preview`).
- Las palabras se transmiten en tiempo real con latencia sub-segundo tanto a pantallas de asistentes como a overlays transparentes para transmisiones de streaming.

---

## ✨ Características Principales

- ⚡ **Latencia Ultra-Baja Sub-Segundo**: Streaming continuo token-por-token sin esperar a que el orador termine la oración.
- 🪶 **Motor en Go (Memoria < 20MB)**: Diseñado en Golang con goroutines y canales para correr decenas de escenarios concurrentes consumiendo una fracción mínima de CPU y RAM.
- 🎯 **Glosario Técnico y Preservación de Entidades**: Instrucción de sistema diseñada para no traducir términos de ingeniería de software (`Kubernetes`, `GraphQL`, `gRPC`, `Kafka`, `CI/CD`, `Golang`, `PostgreSQL`).
- 🇦🇷 **Comprensión de Lunfardo y Modismos**: Adaptado al español rioplatense, voseo y Spanglish técnico.
- 🎥 **Listo para OBS Studio y vMix**: Modo overlay transparente (`/?overlay=true&session=stage-1`) con tipografía de alto contraste lista para producción audiovisual.
- 👥 **Portal para Asistentes (Audience View)**: Interfaz responsiva donde cualquier persona de la audiencia elige el escenario y el idioma de los subtítulos desde su celular o laptop.
- 📥 **Exportación Completa**: Descarga de subtítulos al finalizar cada charla en formatos estándar `.SRT`, `.VTT` o `.TXT`.
- 📁 **Modo de Prueba con Audio Incluido**: Permite probar y evaluar el sistema importando archivos de audio sin requerir micrófono.

---

## 🧩 Requisitos Mínimos (MVP Checklist)

| Requisito | Estado | Implementación |
| :--- | :---: | :--- |
| **Recibir audio en vivo** | ✅ | Micrófono en vivo (AudioWorklet), subida de archivos `.wav`/`.mp3` o audio de prueba integrado. |
| **Transcripción en tiempo real** | ✅ | Sesión streaming paralela en idioma original (Inglés / Español). |
| **Traducción en tiempo real** | ✅ | Sesión streaming directa de audio a texto traducido al Español / Inglés. |
| **Visualización de subtítulos** | ✅ | Consola Web, Modo Audiencia (`/?view=audience`), y OBS Overlay transparente. |
| **Múltiples sesiones en simultáneo** | ✅ | Canales aislados por `SessionID` (`stage-1`, `stage-2`, etc.) manejados concurrentemente en Go. |
| **Audios de prueba en el repo** | ✅ | Archivos incluidos en [`samples/`](./samples/) y botón de reproducción directa en el frontend. |
| **Licencia Open Source OSI** | ✅ | [MIT License](./LICENSE) aprobada por la OSI. |
| **Documentación clara** | ✅ | Guía completa de arquitectura, ejecución local, OBS y Cloud Run. |

---

## 🍬 Opcionales Implementados

- [x] **Integración con OBS Studio y vMix**: Overlays transparentes con efecto backdrop blur y auto-scroll.
- [x] **Glosario Técnico de Conferencia**: Preservación de vocabulario developer en las instrucciones de Gemini.
- [x] **Exportación SRT / VTT / TXT**: Generación y descarga directa con códigos de tiempo precisos.
- [x] **Panel de Monitoreo de Producción**: Monitoreo en tiempo real de escenarios activos y cantidad de oyentes (`/api/sessions`).
- [x] **Soporte Multilingüe Ampliado**: Compatible con Inglés, Español y Portugués (`pt`).

---

## 📊 Benchmarks: Go vs Python Architecture

Durante el desarrollo se construyó una suite de benchmark automatizada ([`benchmark/compare.py`](./benchmark/compare.py)) comparando la implementación Python (FastAPI/asyncio) con esta versión en Golang:

| Métrica | Go Edition (Este Repositorio) | Python FastAPI | Ventaja de Go |
| :--- | :---: | :---: | :---: |
| **Memoria RSS en Reposo** | **14.2 MB** | 84.7 MB | **~6x menor consumo** |
| **Memoria con 10 Sesiones Activas** | **31.5 MB** | 218.0 MB | **~7x menor consumo** |
| **Tiempo de Inicio** | **< 15 ms** | ~1.4 s | **~90x más rápido** |
| **Concurrencia** | Goroutines nativas en `epoll`/`kqueue` | Event loop asyncio de un solo hilo | Aislamiento real multinúcleo |
| **Tamaño de Imagen Docker** | **18 MB** (Alpine scratch) | 480 MB | Despliegue ultra liviano |

---

## 🏗️ Arquitectura del Sistema

```
                      ┌──────────────────────────────────────┐
                      │    Disertante / Escenario Nerdearla  │
                      └──────────────────┬───────────────────┘
                                         │ Web Audio API (PCM 16kHz, 16ms)
                                         ▼
                      ┌──────────────────────────────────────┐
                      │       Nerdearla Live Subs (Go)       │
                      │         (handlers / router)          │
                      └──────────────┬───────────────┬───────┘
                                     │               │
       goroutine A (Transcribe)      │               │ goroutine B (Translate)
                                     ▼               ▼
                        ┌─────────────────────────────────┐
                        │   Google Gemini Live API WSS    │
                        │ generativelanguage.googleapis.com│
                        └────────────────┬────────────────┘
                                         │ Streaming tokens
                                         ▼
                      ┌──────────────────────────────────────┐
                      │   Stage Broadcaster (goroutines)     │
                      └──────┬───────────────────────┬───────┘
                             │                       │
           SSE / WebSocket   │                       │ WebSocket
                             ▼                       ▼
                ┌─────────────────────────┐   ┌─────────────────────────┐
                │ 🎥 OBS / vMix Overlay   │   │ 👥 Asistente Mobile/Web │
                │   (Transparent Alpha)   │   │     (Audience View)     │
                └─────────────────────────┘   └─────────────────────────┘
```

---

## 🚀 Guía de Inicio Rápido

### 1. Configuración de Credenciales

Obtené tu clave de API de **Google Gemini** (desde [Google AI Studio](https://aistudio.google.com/) o habilitando la Generative Language API en tu proyecto de **Google Cloud**):

```bash
cp .env.example .env
```

Editá el archivo `.env`:
```env
GEMINI_API_KEY=AIzaSy...TuClaveAqui
PORT=8080
```

### 2. Correr Localmente

Requisitos: **Go 1.23+** instalado.

```bash
# 1. Clonar el repositorio
git clone https://github.com/YOUR_USER/nerdearla-live-subs-go.git
cd nerdearla-live-subs-go

# 2. Descargar dependencias
go mod download

# 3. Compilar y ejecutar
go run main.go
```

Abrí tu navegador en **http://localhost:8080**.

### 3. Probar en Vivo (Micrófono, YouTube o Multi-Audio Simultáneo)

El sistema ofrece 4 modalidades flexibles de entrada de audio:

#### A. 🎙️ Micrófono en Vivo (`Live Mic`)
- Por defecto, la aplicación inicia en modo **Live Microphone** con el identificador de escenario **`Live Mic`** y traducción configurada de **Español a Inglés** (`es → en`).
- Hacé clic en **"Start Streaming & Subtitles"** y hablá en español: verás la transcripción original y la traducción al inglés en tiempo real token por token.

#### B. 🖥️ Audio de Pestaña de Navegador (Charlas de YouTube en Vivo)
- ¿Querés probar con una charla real de Nerdearla en YouTube sin descargar archivos?
- Seleccioná **"Browser Tab Audio (YouTube & Talks)"** en la consola.
- Abrí cualquier video de [YouTube Nerdearla](https://youtube.com/nerdearla) en otra pestaña (por ejemplo, una charla en inglés o en español).
- Hacé clic en **"Start Streaming"**, elegí la pestaña de YouTube y tildá **"Compartir audio de la pestaña"**.
- Al reproducir el video, el audio digital se transmite en directo a Gemini Live para subtitulado y traducción simultánea.

#### C. 📁 Muestras de Charlas Reales de Nerdearla (Español e Inglés)
- Seleccioná **"Built-in Nerdearla Talk Samples"**.
- Elegí entre:
  * 🇦🇷 **Charla Nerdearla en Español** (`nerdearla-talk-spanish.wav`, 25s): autoconfigura traducción a inglés en el escenario `stage-spanish`.
  * 🇺🇸 **Keynote Nerdearla en Inglés** (`nerdearla-talk-english.wav`, 17s): autoconfigura traducción a español en el escenario `stage-english`.
- Hacé clic en **"Start Streaming"** para reproducir.

#### D. 🚀 Demo Simultánea Multi-Escenario (Multi-File Real-Time Demo)
- Seleccioná **"Upload Audio Files"**.
- Podés arrastrar múltiples archivos de audio (`.wav`, `.mp3`, `.m4a`) o hacer clic en **"⚡ Quick Demo: Preload 2 Simultaneous Stages"**.
- Hacé clic en **"🚀 Stream All Files in Parallel (Multi-Stage Live Demo)"**.
- El sistema procesará ambos escenarios en paralelo:
  * Escenario 1 (`stage-spanish`): Transcripción en español + traducción simultánea al inglés.
  * Escenario 2 (`stage-english`): Transcripción en inglés + traducción simultánea al español.
- Desde la **Vista de Audiencia** o el **Monitor de Producción** podrás monitorear ambos escenarios en simultáneo.

---

## 🎥 Integración con OBS Studio & vMix

Para superponer los subtítulos en la transmisión en vivo de Nerdearla:

1. En **OBS Studio**, agregá una fuente de tipo **Browser (Navegador)**.
2. Ingresá la URL con el parámetro de overlay y el ID del escenario correspondiente:
   ```
   http://localhost:8080/?overlay=true&session=stage-1
   ```
3. Configurá las dimensiones del navegador:
   - **Ancho**: `1920`
   - **Alto**: `1080`
   - Marcá la casilla **Shutdown source when not visible** (opcional).
4. El overlay cuenta con fondo transparente, contraste adaptativo y difuminado de fondo (`backdrop-filter`) para garantizar lectura sobre cualquier diapositiva o video.

---

## 👥 Vista para la Audiencia (Audience Mode)

Los asistentes al evento pueden ver los subtítulos directamente en sus dispositivos móviles o portátiles ingresando a:
```
http://localhost:8080/?view=audience
```
- **Selector de escenario**: Elegí `stage-1`, `stage-2` o `keynote`.
- **Modos de lectura**: Solo traducción al español, solo original en inglés o vista dual.
- **Tamaño de fuente ajustable**: Botones `A+` y `A-` para adaptar el texto a cualquier distancia de pantalla.

---

## 🔀 Escalabilidad Multi-Escenario

La arquitectura está construida para soportar múltiples escenarios de Nerdearla en simultáneo (`stage-1`, `stage-2`, `stage-3`, ..., `keynote`):

1. **Aislamiento por Stage**: Cada sesión tiene su propio canal de difusión (`broadcaster.go`) protegido por `sync.RWMutex`. El audio y los subtítulos de una sala jamás se mezclan con otra.
2. **Concurrencia en Goroutines**: Cada nuevo escenario activo solo consume un par de goroutines ligeras (~4KB cada una), lo que permite correr **más de 30 escenarios en paralelo en una sola instancia mínima** de Google Cloud Run con 512MB de RAM.

---

## ☁️ Despliegue en Google Cloud Run

Gracias a los créditos de Google Cloud para el hackathon, podés desplegar la solución en segundos:

### Opción con Docker:
```bash
# Construir la imagen minimalista de Go
docker build -t gcr.io/TU_PROYECTO_GCP/live-subs-go .

# Desplegar en Cloud Run con soporte WebSocket y HTTPS gratis
gcloud run deploy nerdearla-live-subs-go \
  --image gcr.io/TU_PROYECTO_GCP/live-subs-go \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars GEMINI_API_KEY="AIzaSy..." \
  --timeout 3600
```

Cloud Run provee un certificado SSL HTTPS automático (indispensable para que los navegadores autoricen el acceso al micrófono en conferencias).

---

## 📄 Licencia

Este proyecto está bajo la licencia **MIT** - consulta el archivo [LICENSE](./LICENSE) para más detalles. Licencia aprobada por la [Open Source Initiative (OSI)](https://opensource.org/licenses/MIT).

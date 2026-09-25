# 🎙️ Nerdearla Live Subs (Go Edition)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go: 1.23+](https://img.shields.io/badge/Go-1.23%2B-00ADD8.svg)](https://go.dev/)
[![Gemini Live API](https://img.shields.io/badge/Gemini%20Live%20API-Bidirectional%20Streaming-4285F4.svg)](https://ai.google.dev/)
[![Hackathon: Nerdearla Vibeathon](https://img.shields.io/badge/Hackathon-Nerdearla%20Vibeathon%202026-orange.svg)](https://nerdear.la)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)

> **Open-source, ultra-low-latency simultaneous live transcription and translation system built for large-scale multi-stage tech conferences like [Nerdearla](https://nerdear.la).**  
> Streams live speech directly to the **Google Gemini Live API** (`gemini-3.5-transcribe-live` + `gemini-3.5-live-translate-preview`) and broadcasts real-time Spanish, English, and Portuguese subtitles to attendee screens, production monitors, and OBS Studio broadcast overlays.

---

## 📑 Tabla de Contenidos / Table of Contents

- [🌟 Contexto & Misión / Context & Mission](#-contexto--misión)
- [✨ Características Principales / Key Highlights](#-características-principales)
- [🧩 Requisitos Mínimos (MVP Checklist)](#-requisitos-mínimos-mvp-checklist)
- [🍬 Opcionales Implementados (Bonus Features)](#-opcionales-implementados)
- [🎬 Multi-Stage Studio & Video en Vivo / Studio Mode](#-multi-stage-studio--video-en-vivo)
- [📊 Benchmarks: Go vs Python Architecture](#-benchmarks-go-vs-python-architecture)
- [🏗️ Arquitectura del Sistema / Architecture](#️-arquitectura-del-sistema)
- [🚀 Guía de Inicio Rápido / Quick Start](#-guía-de-inicio-rápido)
  - [1. Configuración de Credenciales](#1-configuración-de-credenciales)
  - [2. Correr Localmente (Go)](#2-correr-localmente)
  - [3. Probar Charlas Reales de Nerdearla en Vivo](#3-probar-charlas-reales-de-nerdearla-en-vivo)
- [🎥 Integración con OBS Studio & vMix](#-integración-con-obs-studio--vmix)
- [👥 Vista para la Audiencia (Audience Mode)](#-vista-para-la-audiencia)
- [📊 Monitor de Producción (Stage Monitor)](#-monitor-de-producción)
- [🔀 Escalabilidad Multi-Escenario (30+ Escenarios)](#-escalabilidad-multi-escenario)
- [☁️ Despliegue en Google Cloud Run](#️-despliegue-en-google-cloud-run)
- [📄 Licencia Open Source](#-licencia)

---

## 🌟 Contexto & Misión

En conferencias técnicas internacionales como **Nerdearla**, decenas de charlas simultáneas se dictan en inglés. Para maximizar la accesibilidad de la comunidad hispanohablante, se requiere transcripción e interpretación simultánea.

Las soluciones comerciales tradicionales tienen barreras infranqueables de escala:
1. **Intérpretes humanos simultáneos**: Muy costosos y difíciles de coordinar en conferencias con 5, 10 o 30 escenarios en paralelo.
2. **Pipelines ASR tradicionales en cascada (Whisper batch + Traducción text-to-text)**: Tienen latencias de 4 a 8 segundos, rompiendo la sincronía entre el disertante y la pantalla.

**Nerdearla Live Subs (Go Edition)** resuelve esto mediante un motor en **Golang de alto rendimiento** conectado de forma bidireccional por WebSocket a la **Google Gemini Live API**:
- La voz se captura y procesa en chunks de 16ms (PCM 16kHz) con downsampling nativo en `AudioWorkletNode`.
- Dos sesiones concurrentes de Gemini generan en paralelo:
  1. **Transcripción en idioma original** (`gemini-3.5-transcribe-live`).
  2. **Traducción simultánea** (`gemini-3.5-live-translate-preview`).
- Las palabras se transmiten en tiempo real con latencia sub-segundo tanto a pantallas de asistentes como a overlays transparentes para transmisiones de streaming.

---

## ✨ Características Principales

- ⚡ **Latencia Ultra-Baja Sub-Segundo**: Streaming continuo token-por-token sin esperar a que el orador termine la frase.
- 🪶 **Motor en Go (Memoria < 20MB)**: Diseñado en Golang con goroutines y canales para correr decenas de escenarios concurrentes consumiendo una fracción mínima de CPU y RAM.
- 🎬 **Multi-Stage Studio con Video en Vivo**: Vista principal que monta la transmisión de video real de las pestañas de Chrome dentro de las tarjetas del estudio (`● LIVE TAB FEED`).
- 🎧 **Optimización para Auriculares Bluetooth / Galaxy Buds**: Desactiva el filtro de software del navegador para preservar los fonemas de la voz y añade un slider de **Mic Boost** (+3.5 dB).
- 🎯 **Glosario Técnico y Preservación de Entidades**: Instrucción de sistema diseñada para no traducir términos de ingeniería de software (`Kubernetes`, `GraphQL`, `gRPC`, `Kafka`, `CI/CD`, `Golang`, `PostgreSQL`).
- 🇦🇷 **Comprensión de Lunfardo y Modismos**: Adaptado al español rioplatense, voseo y Spanglish técnico.
- 🎥 **Listo para OBS Studio y vMix**: Modo overlay transparente (`/?overlay=true&session=stage-1`) con tipografía de alto contraste lista para producción audiovisual.
- 👥 **Portal para Asistentes (Audience View)**: Interfaz responsiva donde cualquier persona de la audiencia elige el escenario y el idioma de los subtítulos desde su celular o laptop.
- 📥 **Exportación Completa**: Descarga de subtítulos al finalizar cada charla en formatos estándar `.SRT`, `.VTT` o `.TXT`.
- 📊 **Monitor de Producción en Tiempo Real**: Estadísticas en vivo de conexiones, espectadores por sala y salud de los WebSockets (`/api/sessions`).

---

## 🧩 Requisitos Mínimos (MVP Checklist)

| Requisito | Estado | Implementación |
| :--- | :---: | :--- |
| **Recibir audio en vivo** | ✅ | Micrófono en vivo (AudioWorklet), captura de audio de pestañas de navegador (YouTube), o subida de archivos `.wav`/`.mp3`. |
| **Transcripción en tiempo real** | ✅ | Sesión streaming paralela en idioma original (Inglés / Español) vía `gemini-3.5-transcribe-live`. |
| **Traducción en tiempo real** | ✅ | Sesión streaming directa de audio a texto traducido vía `gemini-3.5-live-translate-preview`. |
| **Visualización de subtítulos** | ✅ | Multi-Stage Studio, Consola Web, Modo Audiencia (`/?view=audience`), y OBS Overlay transparente. |
| **Múltiples sesiones en simultáneo** | ✅ | 4 sesiones activas en simultáneo (1 Micrófono Presentador + 3 Charlas YouTube en paralelo), escalable a 30+. |
| **Audios de prueba en el repo** | ✅ | Archivos incluidos en [`samples/`](./samples/) y enlaces directos a charlas reales de Nerdearla. |
| **Licencia Open Source OSI** | ✅ | [MIT License](./LICENSE) aprobada por la OSI. |
| **Documentación clara** | ✅ | Guía completa de arquitectura, ejecución local, OBS y Cloud Run. |

---

## 🍬 Opcionales Implementados

- [x] **Integración con OBS Studio y vMix**: Overlays transparentes con efecto backdrop blur y auto-scroll con 1-clic (`📺 OBS`).
- [x] **Glosario Técnico de Conferencia**: Preservación de vocabulario developer en las instrucciones de Gemini.
- [x] **Exportación SRT / VTT / TXT**: Generación y descarga directa con códigos de tiempo precisos.
- [x] **Panel de Monitoreo de Producción**: Monitoreo en tiempo real de escenarios activos y cantidad de oyentes (`/api/sessions`).
- [x] **Soporte Multilingüe Ampliado**: Compatible con Inglés, Español y Portugués (`pt`).
- [x] **Transmisión de Video en Vivo en la App**: Captura simultánea de video y audio desde pestañas de Chrome dentro de las tarjetas del estudio.

---

## 🎬 Multi-Stage Studio & Video en Vivo

La pantalla principal de la aplicación (**Multi-Stage Studio**) está diseñada para producción en vivo y demostraciones de streaming en OBS (formato 16:9 de 1400px):

```
┌────────────────────────────────────────────────────────────────────────────────────────────────┐
│  🎙️ Nerdearla Live Subs [Go Engine]     [🎬 Multi-Stage Studio] [🎛️ Presenter Console] [👥]    │
├────────────────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                                │
│  ┌──────────────────────────────────────────────────────────────────────────────────────────┐  │
│  │ 🎙️ LIVE PRESENTER BROADCAST (YOU)           [Ready / Live On Air] [VU Meter ■■■■■■□□□□]  │  │
│  ├──────────────────────────────────────────┬───────────────────────────────────────────────┤  │
│  │ • Mic: [🎧 Galaxy Buds Pro 3 / Built-in] │ Live Presenter Subtitles               ES → EN │  │
│  │ • Lang: [🇦🇷 Español] ⇄ [🇺🇸 English]    │ Original: "Bienvenidos a Nerdearla 2026..."   │  │
│  │ • Boost: 1.5x [────●────]                │ Trans:    "Welcome to Nerdearla 2026..."      │  │
│  │ [▶ Start Presenter Mic]  [📺 OBS Overlay]│                                               │  │
│  └──────────────────────────────────────────┴───────────────────────────────────────────────┘  │
│                                                                                                │
│  ┌──────────────────────────────────────────────────────────────────────────────────────────┐  │
│  │ 🎥 Simultaneous Video Stages (Real YouTube Tabs)                  [🔗 Open All 3 Tabs ↗]   │  │
│  │ Quick Open: [🏛️ 1. Keynote (EN→ES) ↗] [🎙️ 2. Main Talk (ES→EN) ↗] [🔊 3. Noisy Room ↗]   │  │
│  └──────────────────────────────────────────────────────────────────────────────────────────┘  │
│                                                                                                │
│  ┌──────────────────────────┐  ┌──────────────────────────┐  ┌──────────────────────────┐      │
│  │ 🏛️ Stage 1: Keynote       │  │ 🎙️ Stage 2: Main Talk    │  │ 🔊 Stage 3: Noisy Room   │      │
│  │ [Ready / Live Tab] [VU]  │  │ [Ready / Live Tab] [VU]  │  │ [Ready / Live Tab] [VU]  │      │
│  │ ┌──────────────────────┐ │  │ ┌──────────────────────┐ │  │ ┌──────────────────────┐ │      │
│  │ │ ● LIVE TAB VIDEO FEED│ │  │ │ ● LIVE TAB VIDEO FEED│ │  │ │ ● LIVE TAB VIDEO FEED│ │      │
│  │ │ (Plays real video!)  │ │  │ │ (Plays real video!)  │ │  │ │ (Plays real video!)  │ │      │
│  │ └──────────────────────┘ │  │ └──────────────────────┘ │  │ └──────────────────────┘ │      │
│  │ 🔗 Open Tab  📺 OBS      │  │ 🔗 Open Tab  📺 OBS      │  │ 🔗 Open Tab  📺 OBS      │      │
│  │ [🇺🇸 EN] → [🇦🇷 ES]       │  │ [🇦🇷 ES] → [🇺🇸 EN]       │  │ [🇦🇷 ES] → [🇺🇸 EN]       │      │
│  │ [🖥️ Capture Video&Audio] │  │ [🖥️ Capture Video&Audio] │  │ [🖥️ Capture Video&Audio] │      │
│  │ ┌──────────────────────┐ │  │ ┌──────────────────────┐ │  │ ┌──────────────────────┐ │      │
│  │ │ Orig: K8s architecture │  │ │ Orig: Modelos de IA    │  │ │ Orig: Con ruido de fondo│ │      │
│  │ │ Trans: Arquitectura K8s│  │ │ Trans: AI Models in... │  │ │ Trans: With background..│ │      │
│  │ └──────────────────────┘ │  │ └──────────────────────┘ │  │ └──────────────────────┘ │      │
│  └──────────────────────────┘  └──────────────────────────┘  └──────────────────────────┘      │
└────────────────────────────────────────────────────────────────────────────────────────────────┘
```

### Charlas Reales de Nerdearla Preconfiguradas:
1. **Stage 1 (Keynote)**: [Cloud Native & Kubernetes Talk](https://www.youtube.com/watch?v=GVadbxHks_A&t=92s) (🇺🇸 Inglés → 🇦🇷 Español).
2. **Stage 2 (Main Talk)**: [Inteligencia Artificial y Open Source](https://www.youtube.com/watch?v=ObElMurNobI) (🇦🇷 Español → 🇺🇸 Inglés).
3. **Stage 3 (Noisy Room)**: [Charla con Ruido Ambiente de Sala](https://www.youtube.com/watch?v=_bV_-ZLNWho) (🇦🇷 Español → 🇺🇸 Inglés).

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
git clone https://github.com/tobiasmunoz/nerdearla-vibeathon-2026.git
cd nerdearla-vibeathon-2026

# 2. Descargar dependencias
go mod download

# 3. Compilar y ejecutar
go run main.go
```

Abrí tu navegador en **http://localhost:8080**.

### 3. Probar Charlas Reales de Nerdearla en Vivo

1. En la barra superior, hacé clic en los botones de acceso rápido para abrir las charlas de YouTube en pestañas separadas:
   - `🏛️ Stage 1: Keynote (EN → ES) ↗`
   - `🎙️ Stage 2: Main Talk (ES → EN) ↗`
   - `🔊 Stage 3: Noisy Room (ES → EN) ↗`
2. En cada tarjeta del escenario, hacé clic en **`🖥️ Capture Tab Video & Audio`**.
3. Seleccioná la pestaña correspondiente y asegurate de tildar **"Compartir audio de la pestaña"** (*Also share tab audio*).
4. El video en vivo de la pestaña se montará directamente dentro de la tarjeta de la app y Gemini Live subtitulará en tiempo real.
5. Encendé tu micrófono en la tarjeta superior (**Live Presenter Broadcast**) para hablar en vivo simultáneamente con todas las charlas corriendo.

---

## 🎥 Integración con OBS Studio & vMix

Para superponer los subtítulos en la transmisión en vivo de Nerdearla:

1. En **OBS Studio**, agregá una fuente de tipo **Browser (Navegador)**.
2. Ingresá la URL con el parámetro de overlay y el ID del escenario correspondiente (o usá el botón **📺 OBS** en cualquier tarjeta):
   ```
   http://localhost:8080/?overlay=true&session=stage-keynote
   ```
3. Configurá las dimensiones del navegador:
   - **Ancho**: `1920`
   - **Alto**: `1080`
4. El overlay cuenta con fondo transparente, contraste adaptativo y difuminado de fondo (`backdrop-filter`) para garantizar lectura sobre cualquier diapositiva o video.

---

## 👥 Vista para la Audiencia (Audience Mode)

Los asistentes al evento pueden ver los subtítulos directamente en sus dispositivos móviles o portátiles ingresando a:
```
http://localhost:8080/?view=audience
```
- **Selector de escenario**: Elegí `stage-keynote`, `stage-workshop`, `stage-noisy` o `Live Mic`.
- **Modos de lectura**: Solo traducción, solo original o vista dual.
- **Tamaño de fuente ajustable**: Botones `A+` y `A-` para adaptar el texto a cualquier distancia de pantalla.

---

## 📊 Monitor de Producción

El equipo de producción audiovisual y los organizadores pueden supervisar la salud del sistema ingresando a la pestaña **`📊 Stage Monitor`** (`/?view=monitor`):
- Muestra el estado activo/inactivo de cada sala en tiempo real.
- Cantidad de asistentes conectados por sala.
- Acceso directo a los enlaces de overlay de OBS.

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

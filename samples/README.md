# 🎙️ Sample Audio Files for Testing

This directory contains test audio files to evaluate **Nerdearla Live Subs** in real time without needing a live microphone.

## Available Sample Files
1. **`nerdearla-talk-spanish.wav`** (25s, 16kHz mono PCM)
   - Real conference speech in Spanish discussing open source, AI in streaming, and distributed systems.
   - Ideal for testing Spanish transcription and Spanish → English translation.
2. **`nerdearla-talk-english.wav`** (17s, 16kHz mono PCM)
   - Real keynote speech in English discussing cloud architecture, real-time data pipelines, and accessibility.
   - Ideal for testing English transcription and English → Spanish translation.
3. **`sample-talk-test.wav`** (5s, 16kHz mono PCM)
   - Quick test tone sample.

---

## Testing Methods

### Method 1: Simultaneous Multi-Stage Demo (Recommended for Judges!)
1. Open the web interface at `http://localhost:8080`.
2. Select **"Upload Audio Files"** in the Audio Input Source selector.
3. Click **"⚡ Quick Demo: Preload 2 Simultaneous Stages (Spanish & English Talks)"**.
4. Click **"🚀 Stream All Files in Parallel (Multi-Stage Live Demo)"**.
5. Both talks stream concurrently to Gemini Live:
   - **`stage-spanish`**: Spanish speech → translated to English in real time.
   - **`stage-english`**: English speech → translated to Spanish in real time.
6. Switch tabs in the Presenter Console or open the **Audience View** and **Stage Monitor** to see both live streams simultaneously!

### Method 2: Live Browser Tab Audio (YouTube Videos)
1. In the web interface, choose **"🖥️ Browser Tab Audio (YouTube & Talks)"**.
2. Open any Nerdearla talk on [YouTube](https://youtube.com/nerdearla) in another tab.
3. Click **▶ Start Streaming & Subtitles**.
4. In the browser dialog, select **Chrome Tab**, pick your YouTube tab, and check **"Also share tab audio"**.
5. Press Play on YouTube: the audio will stream directly to Gemini Live for simultaneous subtitles.

### Method 3: Upload Custom Audio Files
Upload one or more `.wav`, `.mp3`, or `.m4a` files directly from your computer. You can configure individual stages and languages for each file and run them in parallel or sequentially.

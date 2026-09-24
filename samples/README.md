# 🎙️ Sample Audio Files for Testing

This directory contains test audio files to evaluate **Nerdearla Live Subs** without needing a live microphone.

## Files
- `sample-talk-test.wav`: A 16kHz 16-bit mono PCM audio sample formatted specifically for the Gemini Live API.

## How to Test with Sample Audio

### Method 1: Using the Web UI (Recommended)
1. Open the web interface at `http://localhost:8080`.
2. In the **Audio Source** selector, choose **"Upload Audio File"** or **"Use Sample Audio"**.
3. Select `samples/sample-talk-test.wav` (or any `.wav` / `.mp3` / `.m4a` file from a past Nerdearla conference talk).
4. Click **Start Translating**.
5. The audio is decoded in the browser, resampled to 16kHz PCM, and streamed in real time to the server and Gemini Live API.

### Method 2: Testing with Real Nerdearla Talks
You can download audio from any past Nerdearla conference talk (e.g. from [YouTube Nerdearla](https://youtube.com/nerdearla)) using `yt-dlp` or `ffmpeg`:

```bash
# Extract 16kHz mono WAV from any video
ffmpeg -i your_talk_video.mp4 -ar 16000 -ac 1 -c:a pcm_s16le samples/talk-demo.wav
```
Then load it through the web interface!

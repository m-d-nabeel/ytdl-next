import { randomUUID } from "crypto";
import { NextResponse } from "next/server";
import ytdl from "ytdl-core";
import fs from "fs";
import Ffmpeg from "fluent-ffmpeg";

export async function POST(req: Request) {
  try {
    const { url } = await req.json();
    if (!url) {
      return new NextResponse("Missing url", { status: 400 });
    }

    if (!ytdl.validateURL(url)) {
      return new NextResponse("Invalid youtube url", { status: 400 });
    }

    const info = await ytdl.getBasicInfo(url);
    const { title: videoTitle } = info.videoDetails;

    const { videoPath, audioPath } = await downloadMedia(url);

    const mixedPath = `/tmp/downloaded/${videoTitle.replace(
      /[^a-z0-9]/gi,
      "_"
    )}.mp4`;

    await mergeAudioAndVideo(videoPath, audioPath, mixedPath);

    deleteTimer(mixedPath);

    return NextResponse.json({
      message: "Download Success",
      title: `${videoTitle}.mp4`,
      info,
    });
  } catch (error) {
    console.error("Download Server Error:", error);
    return new NextResponse("Internal Error", { status: 500 });
  }
}

async function downloadMedia(
  url: string
): Promise<{ videoPath: string; audioPath: string }> {
  const videoPath = `/tmp/downloaded/${randomUUID()}.mp4`;
  const audioPath = `/tmp/downloaded/${randomUUID()}.mp3`;

  const video = ytdl(url, { quality: "highestvideo" });
  const audio = ytdl(url, { quality: "highestaudio" });

  const downloadVideo = streamToFile(video, videoPath);
  const downloadAudio = streamToFile(audio, audioPath);

  await Promise.all([downloadVideo, downloadAudio]);

  return { videoPath, audioPath };
}

function streamToFile(stream: any, filePath: string): Promise<void> {
  return new Promise((resolve, reject) => {
    stream.pipe(fs.createWriteStream(filePath));
    stream.on("end", () => {
      console.log(`${filePath.split("/").pop()} finished downloading.`);
      resolve();
    });
    stream.on("error", (error: any) => {
      console.error(`Error downloading ${filePath.split("/").pop()}:`, error);
      reject(error);
    });
  });
}

async function mergeAudioAndVideo(
  videoPath: string,
  audioPath: string,
  mixedPath: string
): Promise<void> {
  return new Promise((resolve, reject) => {
    Ffmpeg.setFfmpegPath("/usr/bin/ffmpeg");
    Ffmpeg()
      .input(videoPath)
      .input(audioPath)
      .outputOptions(["-acodec copy", "-vcodec copy"])
      .output(mixedPath)
      .on("end", () => {
        console.log("Audio and video merged successfully:", mixedPath);
        fs.unlink(audioPath, (err) => {
          if (err) {
            console.log("Error deleting audio file:");
            throw err;
          }
        });
        fs.unlink(videoPath, (err) => {
          if (err) {
            console.log("Error deleting video file:");
            throw err;
          }
        });
        resolve();
      })
      .on("error", (err) => {
        console.error("Error merging audio and video:", err);
        reject(err);
      })
      .run();
  });
}

function deleteTimer(path: string) {
  console.log("Timer set");
  setTimeout(() => {
    fs.unlink(path, (error) => {
      if (error) {
        console.error("Error deleting file:", error);
      } else {
        console.log("File deleted successfully.");
      }
    });
  }, 15 * 60 * 1000);
}

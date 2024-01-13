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

    const video = ytdl(url, {
      quality: "highestvideo",
    });
    const audio = ytdl(url, {
      quality: "highestaudio",
    });

    const videoTitle = info.videoDetails.title.replace(/[^a-z0-9]/gi, "_");
    const videoPath = `/tmp/video/${randomUUID()}.mp4`;
    const audioPath = `/tmp/video/${randomUUID()}.mp3`;
    const mixedPath = `/tmp/video/${videoTitle}.mp4`;

    const videoPromise = new Promise<void>((resolve, reject) => {
      video.pipe(fs.createWriteStream(videoPath));
      video.on("end", () => {
        console.log("Video finished downloading.");
        resolve();
      });
      video.on("error", (error) => {
        console.error("Error downloading video:", error);
        reject(error);
      });
    });
    const audioPromise = new Promise<void>((resolve, reject) => {
      audio.pipe(fs.createWriteStream(audioPath));
      audio.on("end", () => {
        console.log("Audio finished downloading.");
        resolve();
      });
      audio.on("error", (error) => {
        console.error("Error downloading audio:", error);
        reject(error);
      });
    });
    try {
      await Promise.all([videoPromise, audioPromise]);
      console.log("Both audio and video downloaded successfully.");
    } catch (error) {
      console.error("Error downloading audio and video:", error);
      return new NextResponse("Download Failed in server", { status: 500 });
    }
    
    Ffmpeg.setFfmpegPath("/usr/bin/ffmpeg");
    Ffmpeg()
      .input(videoPath)
      .input(audioPath)
      .outputOptions(["-acodec copy", "-vcodec copy"])
      .output(mixedPath)
      .on("end", () => {
        console.log("Audio and video merged successfully.", mixedPath);
        fs.rm(audioPath, (error) => {
          !!error && console.error("Error deleting audio file:", error);
        });
        fs.rm(videoPath, (error) => {
          !!error && console.error("Error deleting video file:", error);
        });
      })
      .on("error", (error: any) => {
        console.error("Error merging audio and video:", error);
        return new NextResponse("Internal Error", { status: 500 });
      })
      .run();

    deleteTimer(mixedPath);

    return NextResponse.json({
      message: "Download Success",
      title: videoTitle + ".mp4",
      info,
    });
  } catch (error) {
    console.log("Download Server Error");
    return new NextResponse("Internal Error", { status: 500 });
  }
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

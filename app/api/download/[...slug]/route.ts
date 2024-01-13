import { NextResponse } from "next/server";
import ytdl from "ytdl-core";
import fs from "fs";
import readline from "readline";
import Ffmpeg from "fluent-ffmpeg";

export async function GET(
  req: Request,
  { params }: { params: { slug: string[] } }
) {
  try {
    console.log(params.slug);
    const { quality } = await req.json();
    const url = params.slug[0];
    if (!url) {
      return new NextResponse("Missing url", { status: 400 });
    }

    if (!ytdl.validateURL(url)) {
      return new NextResponse("Invalid youtube url", { status: 400 });
    }
    const info = await ytdl.getBasicInfo(url);
    const { videoDetails } = info;
    const resObj: YTVideoDetail = {
      title: videoDetails.title,
      embed: {
        url: videoDetails.embed?.iframeUrl,
        width: videoDetails.embed?.width,
        height: videoDetails.embed?.height,
      },
      thumbnail_url: videoDetails.thumbnails.map((thumb) => thumb.url),
      quality,
    };
    return NextResponse.json(resObj, { status: 200 });
  } catch (error) {
    console.log("[download/[...slug]]", error);
    return new NextResponse("Internal server error", { status: 500 });
  }
}

/********************************************************************/

const QUALITY_MAP: Record<string, string> = {
  low: "lowest",
  medium: "highest",
  high: "highestaudio_highestvideo",
  audio_only: "highestaudio",
};

export async function POST(
  req: Request,
  { params }: { params: { slug: string[] } }
) {
  try {
    const url = params.slug[0];
    const { quality } = await req.json();
    const qualityOption = QUALITY_MAP[quality.toLowerCase()];
    console.log(qualityOption, url, quality);
    if (!qualityOption) {
      return new NextResponse("Invalid quality", { status: 400 });
    }
    if (!url) {
      return new NextResponse("Missing url", { status: 400 });
    }

    if (!ytdl.validateURL(url)) {
      return new NextResponse("Invalid youtube url", { status: 400 });
    }
    const info = await ytdl.getBasicInfo(url);

    let title = info.videoDetails.title.replace(/[^a-z0-9]/gi, "_");
    const videoPath = `/tmp/video/${title}video.mp4`;
    const audioPath = `/tmp/video/${title}audio.mp3`;
    const mixedPath = `/tmp/video/${title}mixed.mp4`;
    Ffmpeg.setFfmpegPath("/usr/bin/ffmpeg");
    ////////////////////////////////////////////////////
    if (qualityOption === "highestaudio_highestvideo") {
      const video = ytdl(url, {
        quality: "highestvideo",
      });
      const audio = ytdl(url, {
        quality: "highestaudio",
      });

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
        mergeAudioVideo(videoPath, audioPath, mixedPath);
        deleteTimer(mixedPath);
        title = title + "mixed.mp4";
        console.log("Both audio and video downloaded successfully.");
      } catch (error) {
        console.error("Error downloading audio and video:", error);
        return new NextResponse("Download Failed in server", { status: 500 });
      }
    } else if (qualityOption === "highestaudio") {
      let stream = ytdl(url, {
        quality: "highestaudio",
      });

      let start = Date.now();
      Ffmpeg(stream)
        .audioBitrate(128)
        .save(audioPath)
        .on("progress", (p) => {
          readline.cursorTo(process.stdout, 0);
          process.stdout.write(`${p.targetSize}kb downloaded`);
        })
        .on("end", () => {
          console.log(`\ntook ${(Date.now() - start) / 1000}s to complete`);
        });
      title = title + "audio.mp3";
      deleteTimer(audioPath);
    } else if (qualityOption === "lowest") {
      const video = ytdl(url, {
        quality: "lowest",
      });
      video.pipe(fs.createWriteStream(videoPath));
      video.on("end", () => {
        console.log("Video finished downloading.");
      });
      video.on("error", (error) => {
        console.error("Error downloading video:", error);
      });
      title = title + "video.mp4";
      deleteTimer(videoPath);
    } else if (qualityOption === "highest") {
      const video = ytdl(url, {
        quality: "highest",
      });
      video.pipe(fs.createWriteStream(videoPath));
      video.on("end", () => {
        console.log("Video finished downloading.");
      });
      video.on("error", (error) => {
        console.error("Error downloading video:", error);
      });
      title = title + "video.mp4";
      deleteTimer(videoPath);
    }
    return NextResponse.json({
      message: "Download Success",
      title,
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

const mergeAudioVideo = (
  videoPath: string,
  audioPath: string,
  mixedPath: string
) => {
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
};

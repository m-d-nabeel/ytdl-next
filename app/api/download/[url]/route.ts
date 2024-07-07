import { NextResponse } from "next/server";
import ytdl from "ytdl-core";
import fs from "fs";
import Ffmpeg from "fluent-ffmpeg";
import { YTVideoDetail } from "@/types";
import { Readable } from "stream";
import path from "path";

const QUALITY_MAP: Record<string, string> = {
  low: "lowest",
  medium: "highest",
  high: "highestaudio_highestvideo",
  audio_only: "highestaudio",
} as const;

async function downloadFile(
  url: string,
  qualityOption: string,
  title: string
): Promise<string> {

  console.log("Downloading file with quality:", qualityOption);

  const isAudio = qualityOption.includes("audio");
  const filename = isAudio ? `${title}.mp3` : `${title}.mp4`;
  const downloadPath = path.join("/tmp/downloaded", filename);

  // Get the file stream
  // const writeStream = fs.createWriteStream(downloadPath);

  await new Promise<void>((resolve, reject) => {
    ytdl(url, { quality: qualityOption })
      .pipe(fs.createWriteStream(downloadPath))
      .on("finish", resolve)
      .on("error", reject)
      .on("end", resolve)
  });

  console.log(downloadPath);

  // const fileStream = ytdl(url, { quality: qualityOption });

  // await new Promise<void>((resolve, reject) => {
  //   Ffmpeg(fileStream)
  //     .toFormat("mp3")
  //     .saveToFile(downloadPath)
  //     .on("error", (err) => {
  //       console.error("Error downloading file:", err);
  //       reject();
  //     })
  //     .on("end", () => {
  //       console.log("Finished merging!");
  //       resolve();
  //     });
  // });
  return downloadPath;
}

async function deleteFile(path: string): Promise<void> {
  try {
    await fs.promises.unlink(path);
    console.log("File deleted successfully:", path);
  } catch (error) {
    console.error("Error deleting file:", error);
  }
}

async function mergeAudioVideo(
  videoStream: Readable,
  audioPath: string,
  mixedPath: string
): Promise<void> {
  await new Promise<void>((resolve, reject) => {
    Ffmpeg(videoStream)
      .input(audioPath)
      .saveToFile(mixedPath)
      .on("error", (err) => {
        reject(err);
      })
      .on("end", () => {
        console.log("Finished merging!");
        resolve();
      });
  });

  console.log("Audio and video merged successfully.");
}

export async function GET(
  req: Request,
  { params }: { params: { url: string } }
) {
  try {
    const { url } = params;
    if (!url) { return new NextResponse("Missing URL", { status: 400 }); }

    if (!ytdl.validateURL(url))
      return new NextResponse("Invalid YouTube URL", { status: 400 });

    const { quality } = await req.json();
    const info = await ytdl.getBasicInfo(url);

    const resObj: YTVideoDetail = {
      title: info.videoDetails.title,
      embed: {
        url: info.videoDetails.embed?.iframeUrl,
        width: info.videoDetails.embed?.width,
        height: info.videoDetails.embed?.height,
      },
      thumbnail_url: info.videoDetails.thumbnails.map((thumb) => thumb.url),
      quality,
    };

    return NextResponse.json(resObj, { status: 200 });
  } catch (error) {
    console.error("[GET] Error:", error);
    return new NextResponse("Internal Server Error", { status: 500 });
  }
}

export async function POST(
  req: Request,
  { params }: { params: { url: string } }
) {
  try {
    const { url } = params;
    if (!url) {
      return new NextResponse("Missing URL", { status: 400 });
    }

    if (!ytdl.validateURL(url)) {
      return new NextResponse("Invalid YouTube URL", { status: 400 });
    }

    const { quality } = await req.json();
    const qualityOption = QUALITY_MAP[quality.toLowerCase()];

    if (!qualityOption) {
      return new NextResponse("Invalid Quality", { status: 400 });
    }

    const info = await ytdl.getBasicInfo(url);
    let title = info.videoDetails.title.replace(/[^a-z0-9]/gi, "_");

    // Create a directory to store downloaded files
    if (!fs.existsSync("/tmp/downloaded")) {
      fs.mkdirSync("/tmp/downloaded", { recursive: true });
      console.log("Directory created successfully.");
    } else {
      console.log("Directory already exists.");
    }

    let downloadPath: string;
    if (qualityOption === "highestaudio_highestvideo") {
      const audioPath = await downloadFile(url, "highestaudio", title);
      const videoStream = ytdl(url, { quality: "highestvideo" });
      const mixedPath = `/tmp/downloaded/${title}mixed.mp4`;
      await mergeAudioVideo(videoStream, audioPath, mixedPath);
      await deleteFile(audioPath);
      setTimeout(() => deleteFile(mixedPath), 15 * 60 * 1000);
      title += ".mp4";
    } else {
      downloadPath = await downloadFile(url, qualityOption, title);
      setTimeout(() => deleteFile(downloadPath), 15 * 60 * 1000);
      title += qualityOption.includes("audio") ? ".mp3" : ".mp4";
    }

    return NextResponse.json({ message: "Download Success", title, info });
  } catch (error) {
    console.error("[POST] Error:", error);
    return new NextResponse("Internal Server Error", { status: 500 });
  }
}

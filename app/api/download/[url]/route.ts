import { NextResponse } from "next/server";
import ytdl from "ytdl-core";
import fs from "fs";
import Ffmpeg from "fluent-ffmpeg";
import { YTVideoDetail } from "@/types";
import { Readable } from "stream";
import path from "path";
import readline from "readline";

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
  const starttime = Date.now();
  const isAudio = qualityOption.includes("audio");
  const filename = isAudio ? `${title}.mp3` : `${title}.mp4`;
  const downloadPath = path.join("/tmp/downloaded", filename);

  if (fs.existsSync(downloadPath)) {
    return downloadPath;
  }

  await new Promise<void>((resolve, reject) => {
    ytdl(url, { quality: qualityOption })
      .pipe(fs.createWriteStream(downloadPath))
      .on("finish", resolve)
      .on("error", reject)
      .on('progress', (chunkLength, downloaded, total) => {
        const percent = downloaded / total;
        const downloadedMinutes = (Date.now() - starttime) / 1000 / 60;
        const estimatedDownloadTime = (downloadedMinutes / percent) - downloadedMinutes;
        readline.cursorTo(process.stdout, 0);
        process.stdout.write(`${(percent * 100).toFixed(2)}% downloaded `);
        process.stdout.write(`(${(downloaded / 1024 / 1024).toFixed(2)}MB of ${(total / 1024 / 1024).toFixed(2)}MB)\n`);
        process.stdout.write(`running for: ${downloadedMinutes.toFixed(2)}minutes`);
        process.stdout.write(`, estimated time left: ${estimatedDownloadTime.toFixed(2)}minutes `);
        readline.moveCursor(process.stdout, 0, -1);
      })
      .on('end', function () {
        console.log("Download finished!");
        resolve();
      })
  });
  console.log(downloadPath);
  return downloadPath;
}

async function mergeAudioVideo(
  videoStream: Readable,
  audioPath: string,
  mixedPath: string
): Promise<void> {
  const starttime = Date.now();
  await new Promise<void>((resolve, reject) => {
    Ffmpeg(videoStream)
      .input(audioPath)
      .save(mixedPath)
      .on('start', function (commandLine) {
        console.log('Spawned FFmpeg with command: ' + commandLine);
      })
      .on('progress', function (progress) {
        // Clear the current line and move the cursor to the start
        process.stdout.write('\r\x1b[K');

        // Display progress bar and additional information
        const percent = progress.percent / 100;
        process.stdout.write(drawProgressBar(percent));
        process.stdout.write(` ${progress.timemark} processed`);
      })
      .on('end', function () {
        console.log('Merging finished!');
        resolve();
      })
      .on('error', function (err, stdout, stderr) {
        console.error('Error: ' + err.message);
        console.error('FFmpeg stderr: ' + stderr);
        reject(err);
      });
  });

  console.log("Audio and video merged successfully.");
}


export async function POST(req: Request, { params }: { params: { url: string } }
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
      if (fs.existsSync(`/tmp/downloaded/${title}mixed.mp4`)) {
        return NextResponse.json({ message: "Download Success", title: `${title}mixed.mp4`, info });
      }
      const audioPath = await downloadFile(url, "highestaudio", title);
      const videoStream = ytdl(url, { quality: "highestvideo" });
      const mixedPath = `/tmp/downloaded/${title}mixed.mp4`;
      await mergeAudioVideo(videoStream, audioPath, mixedPath);
      await deleteFile(audioPath);
      deleteAfterTimeout(mixedPath, 15 * 60 * 1000);
      title += "mixed.mp4";
    } else {
      downloadPath = await downloadFile(url, qualityOption, title);
      deleteAfterTimeout(downloadPath, 15 * 60 * 1000);
      title += qualityOption.includes("audio") ? ".mp3" : ".mp4";
    }

    return NextResponse.json({ message: "Download Success", title, info });
  } catch (error) {
    console.error("[POST] Error:", error);
    return new NextResponse("Internal Server Error", { status: 500 });
  }
}


// Get video details
export async function GET(req: Request, { params }: { params: { url: string } }) {
  try {
    const { url } = params;
    if (!url) {
      return new NextResponse("Missing URL", { status: 400 });
    }

    if (!ytdl.validateURL(url)) {
      return new NextResponse("Invalid YouTube URL", { status: 400 });
    }

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


// Delete file after a certain timeout
function deleteAfterTimeout(path: string, timeout: number) {
  setTimeout(() => deleteFile(path), timeout);
  console.log("File will be deleted after", timeout, "milliseconds.");
}

async function deleteFile(path: string): Promise<void> {
  try {
    await fs.promises.unlink(path);
    console.log("File deleted successfully:", path);
  } catch (error) {
    console.error("Error deleting file:", error);
  }
}

// Function to draw the progress bar
function drawProgressBar(percent: number) {
  const barLength = 30; // Length of the progress bar
  const filledLength = Math.round(barLength * percent);
  const emptyLength = barLength - filledLength;

  const bar = '█'.repeat(filledLength) + '░'.repeat(emptyLength);
  return `|${bar}| ${Math.round(percent * 100)}%`;
}
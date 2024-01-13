import fs from "fs";
import path from "path";
import { NextResponse } from "next/server";

export const revalidate = 0;

async function getVideoStream(title: string) {
  try {
    const videoDir = path.join("/", "tmp", "video");
    const filePath = path.join(videoDir, title);

    console.log("getVideoStream path:", filePath);

    sleep(1500);

    if (!fs.existsSync(filePath)) {
      console.log("File not found");
      return new Response("File not found", { status: 404 });
    }

    const stat = fs.statSync(filePath);
    const fileSize = stat.size;

    const contType = title.endsWith(".mp4") ? "video/mp4" : "audio/mp3";

    const headers = {
      "Accept-Ranges": "bytes",
      "Content-Type": contType,
      "Content-Length": `${fileSize}`,
    };

    const videoStream: any = fs.createReadStream(filePath);
    console.log("videoStream created");
    return new Response(videoStream, { status: 200, headers });
  } catch (error) {
    console.error("getVideoStream Error:", error);
    return new Response("Internal Server Error", { status: 500 });
  }
}

export async function GET(
  _req: Request,
  { params }: { params: { title: string } }
) {
  try {
    const response = await getVideoStream(params.title);
    return response;
  } catch (error: any) {
    console.error("GET Error:", error);
    return NextResponse.json({ error: error.message }, { status: 500 });
  }
}

function sleep(ms: number) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}

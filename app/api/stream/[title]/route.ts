import fs from "fs";
import path from "path";
import { NextResponse } from "next/server";

async function getVideoStream(title: string) {
  try {
    const videoDir = path.join("/", "tmp", "video");
    const filePath = path.join(videoDir, title);

    if (!fs.existsSync(filePath)) {
      return new Response("File not found", { status: 404 });
    }

    const stat = fs.statSync(filePath);
    const fileSize = stat.size;

    const headers = {
      "Accept-Ranges": "bytes",
      "Content-Type": "video/mp4",
      "Content-Length": `${fileSize}`,
    };

    const videoStream: any = fs.createReadStream(filePath);
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

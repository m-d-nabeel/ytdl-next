"use client";

import Link from "next/link";
import React, { useEffect, useState } from "react";

export const revalidate = 0;

function VideoPlayer({ title }: { title: string }) {
  const [videoUrl, setVideoUrl] = useState<string | null>(null);

  useEffect(() => {
    const fetchVideo = async () => {
      try {
        const response = await fetch(
          `/api/stream/${encodeURIComponent(title)}`,
          {
            method: "GET",
          }
        );

        if (!response.ok) {
          console.log("Illegal Response from Server", response);
          throw new Error("Network response was not ok");
        }

        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        setVideoUrl(url);
      } catch (error) {
        console.error("There was a problem fetching the video:", error);
      }
    };

    fetchVideo()
      .then((_) => {
        console.log("Video Fetched");
      })
      .catch((_) => {
        console.log("Video Fetch Error");
      });
  }, [title]);

  if (!videoUrl)
    return (
      <div className="w-full flex flex-col items-center">
        <p className="text-2xl">Loading...</p>
      </div>
    );

  return (
    <div className="w-full flex flex-col items-center">
      <video controls className="h-[360px]">
        <source src={videoUrl ?? ""} type="video/mp4" />
        Your browser does not support the video tag.
      </video>
      <Link
        download
        href={videoUrl ?? "#"}
        target="_blank"
        rel="noopener noreferrer"
        referrerPolicy="no-referrer"
        className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
      >
        Download
      </Link>
    </div>
  );
}

export default VideoPlayer;

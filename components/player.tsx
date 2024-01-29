"use client";
import Link from "next/link";
import React, { useEffect, useState } from "react";

export const revalidate = 0;

function Player({ title }: { title: string }) {
  const [fileStreamUrl, setFileStreamUrl] = useState<string | null>(null);

  useEffect(() => {
    const fetchFileStream = async () => {
      try {
        const response = await fetch(
          `/api/stream/${encodeURIComponent(title)}`,
          {
            method: "GET",
          }
        );

        if (!response.ok) {
          console.log("Response Not OK", response);
          throw new Error("Network response was not ok");
        }

        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        setFileStreamUrl(url);
      } catch (error) {
        console.error("There was a problem fetching the FileStream:", error);
      }
    };

    fetchFileStream();
  }, [title]);

  if (!fileStreamUrl) {
    return (
      <div className="w-full flex flex-col items-center">
        <p className="text-2xl">Loading...</p>
      </div>
    );
  }
  return (
    <div className="w-full flex flex-col items-center">
      <iframe
        src={fileStreamUrl ?? ""}
        referrerPolicy="no-referrer"
        rel="noreferrer"
        title="Embedded youtube"
        sandbox="allow-same-origin allow-scripts allow-popups allow-forms"
        loading="lazy"
        allowFullScreen
        className="w-auto h-[360px]"
      ></iframe>
      <Link
        download={`${title}`}
        href={fileStreamUrl ?? "#"}
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

export default Player;

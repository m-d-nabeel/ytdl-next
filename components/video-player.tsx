"use client";

import Link from "next/link";
import { useEffect, useState } from "react";

export default function VideoPlayer({ path }: { path: string }) {
  const [isMounted, setIsMounted] = useState(false);
  useEffect(() => {
    setIsMounted(true);
  }, [path]);

  if (!isMounted) {
    return (
      <div className="flex flex-col items-center justify-center w-full h-full">
        <div className="w-32 h-32 border-8 border-gray-300 border-t-slate-800 rounded-full animate-spin" />
      </div>
    );
  }

  return (
    <div className="flex flex-col items-center justify-center w-full h-full">
      <video autoPlay={false} controls className="w-auto h-[360px]">
        <source src={path} type="video/mp4" />
        Your browser does not support the video tag.
      </video>
      <Link
        href={path}
        download
        className="border rounded-md px-4 py-2 bg-blue-500 text-white hover:bg-blue-600"
      >
        Download
      </Link>
    </div>
  );
}

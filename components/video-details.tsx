// "use client";
import { useStore } from "@/hooks/use-store";

export default function VideoDetails({ title }: { title: string }) {
  // if (!videoDetails) {
  //   return null;
  // }

  return (
    <div className="w-full flex justify-center items-center flex-col">
      <h1 className="text-2xl font-bold mb-4">{title}</h1>
      {/* {videoDetails.embed ? (
        <iframe
          src={videoDetails.embed.url}
          height={videoDetails.embed.height / 3}
          width={videoDetails.embed.width / 3}
          referrerPolicy="no-referrer"
          rel="noreferrer"
          title="Embedded youtube"
          itemType="video"
          sandbox="allow-same-origin allow-scripts allow-popups allow-forms"
          loading="lazy"
          allowFullScreen
        ></iframe>
      ) : (
        <p>Sorry, the video embed is not available.</p>
      )} */}
    </div>
  );
}

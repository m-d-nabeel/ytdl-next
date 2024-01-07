import VideoPlayer from "@/components/video-player";

type VideoProps = {
  params: {
    title: string;
  };
};
export default function Video({ params }: VideoProps) {
  const title = decodeURIComponent(params.title).replaceAll("=", "_");
  return (
    <div className="flex justify-center items-center h-full w-full">
      <VideoPlayer title={title} />
    </div>
  );
}

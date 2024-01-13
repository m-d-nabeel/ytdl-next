import VideoDetails from "@/components/video-details";
import VideoPlayer from "@/components/video-player";

type VideoProps = {
  params: {
    title: string;
  };
};
export default function Video({ params }: VideoProps) {
  const title = decodeURIComponent(params.title);
  return (
    <div className="flex flex-col justify-center items-center h-full w-full">
      <VideoDetails title={title} />
      <VideoPlayer title={title} />
    </div>
  );
}

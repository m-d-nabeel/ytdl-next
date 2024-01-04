import VideoPlayer from "@/components/video-player";

type VideoProps = {
  params: {
    slug: string[];
  };
};
export default function Video({ params }: VideoProps) {
  const title = decodeURIComponent(params.slug[0]).replaceAll("=", "_");
  const info = decodeURIComponent(params.slug[1]);
  console.log(JSON.parse(JSON.stringify(info)));
  const path = "/video/" + title;
  console.log(path);
  return (
    <div>
      <VideoPlayer path={path} />
    </div>
  );
}

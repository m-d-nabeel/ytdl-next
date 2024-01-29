import Player from "@/components/player";

type VideoProps = {
  params: {
    title: string;
  };
};
export default function Video({ params }: VideoProps) {
  const title = decodeURIComponent(params.title);
  return (
    <div className="flex flex-col justify-center items-center h-full w-full">
      <Player title={title} />
    </div>
  );
}

import LinkForm from "@/components/link-form";
import VideoDetails from "@/components/video-details";

export default function Home() {
  return (
    <main className="flex min-h-screen w-full flex-col items-center justify-between p-24 text-black">
      <div className="flex flex-col justify-center items-center w-full sm:w-[540px]">
        <LinkForm />
      </div>
    </main>
  );
}

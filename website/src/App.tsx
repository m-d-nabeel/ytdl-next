import { AnimatePresence, motion } from "framer-motion";
import {
  ChevronDown,
  Clock,
  Download,
  Eye,
  Music,
  ThumbsUp,
  Video,
  X,
} from "lucide-react";
import { useState } from "react";
import { AudioSelector } from "./AudioSelector";
import { Format, formatFileSize, VideoData } from "./util";

const getCodecInfo = (format: Format) => {
  const hasVideo = format.vcodec && format.vcodec !== "none";
  const hasAudio = format.acodec !== "none";
  return { hasVideo, hasAudio };
};

const isYouTubeUrl = (url: string) => {
  return url.match(/^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be)\/.+$/);
};

const CodecIndicator = ({ format }: { format: Format }) => {
  const { hasVideo, hasAudio } = getCodecInfo(format);

  return (
    <div className="flex gap-2 items-center text-xs">
      <div
        className="flex items-center gap-1"
        title={`Video: ${hasVideo ? format.vcodec : "No video"}`}
      >
        {hasVideo ? (
          <Video className="w-3 h-3 text-green-600" />
        ) : (
          <X className="w-3 h-3 text-red-600" />
        )}
      </div>
      <div
        className="flex items-center gap-1"
        title={`Audio: ${hasAudio ? format.acodec : "No audio"}`}
      >
        {hasAudio ? (
          <Music className="w-3 h-3 text-green-600" />
        ) : (
          <X className="w-3 h-3 text-red-600" />
        )}
      </div>
    </div>
  );
};

export default function App() {
  const [expanded, setExpanded] = useState(false);
  const [mediaUrl, setMediaUrl] = useState("");
  const [status, setStatus] = useState("");
  const [videoData, setVideoData] = useState<VideoData | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [showAudioSelector, setShowAudioSelector] = useState(false);
  const [selectedVideoFormatId, setSelectedVideoFormatId] =
    useState<string>("");
  const [error, setError] = useState<string>("");

  const handleUrlChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setMediaUrl(e.target.value);
    setError("");
    setStatus("");
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!mediaUrl.trim()) return;

    if (!isYouTubeUrl(mediaUrl)) {
      setError("Please enter a valid YouTube URL");
      return;
    }

    setStatus("Fetching video information...");
    setIsLoading(true);
    setVideoData(null);
    setError("");

    try {
      const response = await fetch(
        `/api/yt/info?url=${encodeURIComponent(mediaUrl)}`
      );
      if (!response.ok) throw new Error("Failed to fetch video information");

      const data: VideoData = await response.json();
      setVideoData(data);
      setStatus("Select a quality to download:");
    } catch (error) {
      setError(
        `Error: ${error instanceof Error ? error.message : "An unknown error occurred"
        }`
      );
    } finally {
      setIsLoading(false);
    }
  };

  const handleDownload = (format: Format) => {
    if (!videoData) return;

    const { hasAudio } = getCodecInfo(format);
    if (!hasAudio) {
      setSelectedVideoFormatId(format.format_id);
      setShowAudioSelector(true);
    } else {
      startDownload(format.format_id);
    }
  };

  const startDownload = (formatId: string, audioFormatId?: string) => {
    if (!videoData) return;

    let finalFormatId = formatId;
    if (audioFormatId) {
      finalFormatId = `${audioFormatId}+${formatId}`;
    }

    const encodedFormatId = encodeURIComponent(finalFormatId);
    const encodedUrl = encodeURIComponent(mediaUrl);
    const downloadUrl = `/api/yt/download?format_id=${encodedFormatId}&url=${encodedUrl}`;

    const link = document.createElement("a");
    link.href = downloadUrl;
    link.download = `${videoData.title || "video"}.mp4`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    setStatus("Download started. Check your browser downloads.");
    setShowAudioSelector(false);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-700 via-pink-500 to-red-500 flex items-center justify-center p-4">
      <motion.div
        initial={{ scale: 0.8, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.5 }}
        className="bg-white rounded-lg shadow-2xl overflow-hidden w-full max-w-2xl"
      >
        <form onSubmit={handleSubmit} className="p-6">
          <div className="flex gap-4">
            <input
              type="url"
              value={mediaUrl}
              onChange={handleUrlChange}
              placeholder="Enter YouTube URL"
              className="flex-1 px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
              required
            />
            <motion.button
              type="submit"
              disabled={isLoading}
              className="bg-purple-600 text-white px-6 py-2 rounded-lg font-semibold disabled:opacity-50"
              whileHover={{ scale: 1.05 }}
              whileTap={{ scale: 0.95 }}
            >
              {isLoading ? "Loading..." : "Fetch"}
            </motion.button>
          </div>
          {error && <p className="mt-2 text-red-500">{error}</p>}
          {status && !error && <p className="mt-2 text-gray-600">{status}</p>}
        </form>

        <AnimatePresence>
          {videoData && (
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.3 }}
            >
              <div className="p-6 pt-0">
                <h2 className="text-2xl font-bold mb-4 text-gray-800 break-words">
                  {videoData.title}
                </h2>
                <div className="flex flex-wrap items-center justify-between gap-4 mb-6">
                  <p className="text-gray-600">{videoData.uploader}</p>
                  <div className="flex flex-wrap items-center gap-4">
                    <span className="flex items-center text-gray-600">
                      <Eye className="w-4 h-4 mr-1" />
                      {videoData.view_count.toLocaleString()}
                    </span>
                    <span className="flex items-center text-gray-600">
                      <ThumbsUp className="w-4 h-4 mr-1" />
                      {videoData.like_count.toLocaleString()}
                    </span>
                    <span className="flex items-center text-gray-600">
                      <Clock className="w-4 h-4 mr-1" />
                      {videoData.duration}s
                    </span>
                  </div>
                </div>
                <div className="relative">
                  <motion.img
                    src={`https://img.youtube.com/vi/${videoData.id}/maxresdefault.jpg`}
                    alt={videoData.title}
                    className="w-full h-64 object-cover rounded-lg"
                    initial={{ y: 20, opacity: 0 }}
                    animate={{ y: 0, opacity: 1 }}
                    transition={{ delay: 0.2 }}
                  />
                </div>
              </div>

              <motion.div
                initial={false}
                animate={{ height: expanded ? "auto" : "0" }}
                transition={{ duration: 0.3 }}
                className="overflow-hidden bg-gray-100"
              >
                <div className="p-6">
                  <h3 className="text-xl font-semibold mb-4">
                    Available Formats
                  </h3>
                  <div className="space-y-2">
                    {videoData.formats
                      .filter((format) => format.filesize && format.resolution)
                      .map((format) => (
                        <motion.div
                          key={format.format_id}
                          className="bg-white p-4 rounded-lg shadow flex items-center justify-between"
                          whileHover={{ scale: 1.02 }}
                          transition={{ duration: 0.2 }}
                        >
                          <div className="flex-1">
                            <div className="flex items-center gap-2">
                              <p className="font-semibold">
                                {format.format_note}
                              </p>
                              <CodecIndicator format={format} />
                            </div>
                            <p className="text-sm text-gray-600">
                              {format.resolution} • {format.ext.toUpperCase()} •{" "}
                              {formatFileSize(format.filesize)}
                            </p>
                          </div>
                          <motion.button
                            onClick={() => handleDownload(format)}
                            whileHover={{ scale: 1.05 }}
                            whileTap={{ scale: 0.95 }}
                            className="bg-blue-500 text-white px-4 py-2 rounded-full text-sm font-semibold ml-4 flex items-center"
                          >
                            <Download className="w-4 h-4 mr-2" />
                            Download
                          </motion.button>
                        </motion.div>
                      ))}
                  </div>
                </div>
              </motion.div>

              <motion.button
                onClick={() => setExpanded(!expanded)}
                className="w-full py-4 bg-gray-200 text-gray-800 font-semibold flex items-center justify-center"
                whileHover={{ backgroundColor: "#e5e7eb" }}
              >
                <ChevronDown
                  className={`w-5 h-5 mr-2 transform transition-transform ${expanded ? "rotate-180" : ""
                    }`}
                />
                {expanded ? "Hide Formats" : "Show Formats"}
              </motion.button>
            </motion.div>
          )}
        </AnimatePresence>

        <AudioSelector
          isOpen={showAudioSelector}
          onClose={() => setShowAudioSelector(false)}
          onSelectAudio={(audioFormatId) => {
            if (selectedVideoFormatId) {
              startDownload(selectedVideoFormatId, audioFormatId);
            }
          }}
          audioFormats={
            videoData?.formats.filter(
              (f) => f.vcodec == "none" && f.acodec !== "none"
            ) || []
          }
        />
      </motion.div>
    </div>
  );
}

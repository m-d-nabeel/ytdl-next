import { motion } from "framer-motion"
import { useState } from "react"
import { ChevronDown, Eye, ThumbsUp, Clock, Video, Music, X } from "lucide-react"

interface Format {
  format_id: string
  resolution: string
  ext: string
  filesize: number
  format_note: string
  acodec: string
  vcodec: string
}

interface VideoData {
  id: string
  title: string
  duration: number
  uploader: string
  view_count: number
  like_count: number
  formats: Format[]
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return "0 Bytes"
  const k = 1024
  const sizes = ["Bytes", "KB", "MB", "GB"]
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Number.parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i]
}

const getCodecInfo = (format: Format) => {
  const hasVideo = format.vcodec !== "none"
  const hasAudio = format.acodec !== "none"
  return { hasVideo, hasAudio }
}

const CodecIndicator = ({ format }: { format: Format }) => {
  const { hasVideo, hasAudio } = getCodecInfo(format)

  return (
    <div className="flex gap-2 items-center text-xs">
      <div className="flex items-center gap-1" title={`Video: ${hasVideo ? format.vcodec : 'No video'}`}>
        {hasVideo ? (
          <Video className="w-3 h-3 text-green-600" />
        ) : (
          <X className="w-3 h-3 text-red-600" />
        )}
      </div>
      <div className="flex items-center gap-1" title={`Audio: ${hasAudio ? format.acodec : 'No audio'}`}>
        {hasAudio ? (
          <Music className="w-3 h-3 text-green-600" />
        ) : (
          <X className="w-3 h-3 text-red-600" />
        )}
      </div>
    </div>
  )
}

export default function MediaDownloader() {
  const [expanded, setExpanded] = useState(false)
  const [mediaUrl, setMediaUrl] = useState('')
  const [status, setStatus] = useState('')
  const [videoData, setVideoData] = useState<VideoData | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!mediaUrl.trim()) return

    setStatus('Fetching video information...')
    setIsLoading(true)
    setVideoData(null)

    try {
      const response = await fetch(`/api/yt/info?url=${encodeURIComponent(mediaUrl)}`)
      if (!response.ok) throw new Error('Failed to fetch video information')

      const data: VideoData = await response.json()
      setVideoData(data)
      setStatus('Select a quality to download:')
    } catch (error) {
      setStatus(`Error: ${error instanceof Error ? error.message : 'An unknown error occurred'}`)
    } finally {
      setIsLoading(false)
    }
  }

  const handleDownload = (formatId: string) => {
    if (!videoData) return

    const downloadUrl = `/api/yt/download?format_id=${formatId}&url=${encodeURIComponent(mediaUrl)}`

    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = `${videoData.title || 'video'}.${videoData.formats.find(f => f.format_id === formatId)?.ext || 'mp4'}`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    setStatus('Download started. Check your browser downloads.')
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-700 via-pink-500 to-red-500 flex items-center justify-center p-4">
      <motion.div
        initial={{ scale: 0.8, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ duration: 0.5 }}
        className="bg-white rounded-lg shadow-2xl overflow-hidden w-full max-w-2xl"
      >
        {/* URL Input Form */}
        <form onSubmit={handleSubmit} className="p-6">
          <div className="flex gap-4">
            <input
              type="url"
              value={mediaUrl}
              onChange={(e) => setMediaUrl(e.target.value)}
              placeholder="Enter media URL"
              className="flex-1 px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
              required
            />
            <button
              type="submit"
              disabled={isLoading}
              className="bg-purple-600 text-white px-6 py-2 rounded-lg font-semibold disabled:opacity-50"
            >
              {isLoading ? 'Loading...' : 'Fetch'}
            </button>
          </div>
          {status && <p className="mt-2 text-gray-600">{status}</p>}
        </form>

        {videoData && (
          <>
            <div className="p-6 pt-0">
              <h2 className="text-3xl font-bold mb-4 text-gray-800">{videoData.title}</h2>
              <div className="flex items-center justify-between mb-6">
                <p className="text-gray-600">{videoData.uploader}</p>
                <div className="flex items-center space-x-4">
                  <span className="flex items-center text-gray-600">
                    <Eye className="w-4 h-4 mr-1" /> {videoData.view_count.toLocaleString()}
                  </span>
                  <span className="flex items-center text-gray-600">
                    <ThumbsUp className="w-4 h-4 mr-1" /> {videoData.like_count.toLocaleString()}
                  </span>
                  <span className="flex items-center text-gray-600">
                    <Clock className="w-4 h-4 mr-1" /> {videoData.duration}s
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
                <h3 className="text-xl font-semibold mb-4">Available Formats</h3>
                <div className="space-y-2">
                  {videoData.formats
                    .filter(format => format.filesize && format.resolution)
                    .map((format) => (
                      <div
                        key={format.format_id}
                        className="bg-white p-4 rounded-lg shadow flex items-center justify-between"
                      >
                        <div className="flex-1">
                          <div className="flex items-center gap-2">
                            <p className="font-semibold">{format.format_note}</p>
                            <CodecIndicator format={format} />
                          </div>
                          <p className="text-sm text-gray-600">
                            {format.resolution} • {format.ext.toUpperCase()} • {formatFileSize(format.filesize)}
                          </p>
                          <p className="text-xs text-gray-500 mt-1">
                            {format.vcodec !== "none" && `Video: ${format.vcodec}`}
                            {format.vcodec !== "none" && format.acodec !== "none" && " • "}
                            {format.acodec !== "none" && `Audio: ${format.acodec}`}
                          </p>
                        </div>
                        <motion.button
                          onClick={() => handleDownload(format.format_id)}
                          whileHover={{ scale: 1.05 }}
                          whileTap={{ scale: 0.95 }}
                          className="bg-blue-500 text-white px-4 py-2 rounded-full text-sm font-semibold ml-4"
                        >
                          Download
                        </motion.button>
                      </div>
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
                className={`w-5 h-5 mr-2 transform transition-transform ${expanded ? "rotate-180" : ""}`}
              />
              {expanded ? "Hide Formats" : "Show Formats"}
            </motion.button>
          </>
        )}
      </motion.div>
    </div>
  )
}

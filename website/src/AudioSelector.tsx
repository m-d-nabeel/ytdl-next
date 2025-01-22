import { motion, AnimatePresence } from "framer-motion"
import { Music, X } from "lucide-react"

interface Format {
  format_id: string
  ext: string
  filesize: number
  format_note: string
  acodec: string
}

interface AudioSelectorProps {
  isOpen: boolean
  onClose: () => void
  onSelectAudio: (formatId: string) => void
  audioFormats: Format[]
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return "0 Bytes"
  const k = 1024
  const sizes = ["Bytes", "KB", "MB", "GB"]
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Number.parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i]
}

export default function AudioSelector({ isOpen, onClose, onSelectAudio, audioFormats }: AudioSelectorProps) {
  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
          onClick={onClose}
        >
          <motion.div
            initial={{ scale: 0.9, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{ scale: 0.9, opacity: 0 }}
            className="bg-white rounded-lg p-6 w-full max-w-md"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-xl font-semibold">Select Audio Format</h3>
              <button onClick={onClose} className="text-gray-500 hover:text-gray-700">
                <X className="w-6 h-6" />
              </button>
            </div>
            <div className="space-y-2">
              {audioFormats.map((format) => (
                <motion.button
                  key={format.format_id}
                  className="w-full bg-gray-100 p-3 rounded-lg flex items-center justify-between hover:bg-gray-200 transition-colors"
                  onClick={() => onSelectAudio(format.format_id)}
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.98 }}
                >
                  <div className="flex items-center">
                    <Music className="w-5 h-5 mr-3 text-blue-500" />
                    <div className="text-left">
                      <p className="font-medium">{format.format_note}</p>
                      <p className="text-sm text-gray-600">
                        {format.ext.toUpperCase()} • {formatFileSize(format.filesize)}
                      </p>
                    </div>
                  </div>
                  <span className="text-blue-500">Select</span>
                </motion.button>
              ))}
            </div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  )
}


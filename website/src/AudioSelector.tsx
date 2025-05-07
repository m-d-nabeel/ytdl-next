import { AnimatePresence, motion } from "framer-motion";
import { Music, X, CheckCircle } from "lucide-react";
import { Format, formatFileSize } from "./util";

export const AudioSelector = ({
  isOpen,
  onClose,
  onSelectAudio,
  audioFormats,
}: {
  isOpen: boolean;
  onClose: () => void;
  onSelectAudio: (formatId: string) => void;
  audioFormats: Format[];
}) => (
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
            <button
              onClick={onClose}
              className="text-gray-500 hover:text-gray-700"
            >
              <X className="w-6 h-6" />
            </button>
          </div>
          
          <div className="mb-4">
            <div className="bg-blue-50 p-3 rounded-lg text-sm text-blue-700 flex items-start">
              <CheckCircle className="w-5 h-5 mr-2 mt-0.5 flex-shrink-0" />
              <span>Formats marked as "Compatible" will work on most devices and browsers.</span>
            </div>
          </div>
          
          <div className="space-y-2">
            {audioFormats.map((format) => (
              <motion.button
                key={format.format_id}
                className={`w-full p-3 rounded-lg flex items-center justify-between transition-colors
                  ${format.is_compatible 
                    ? "bg-green-50 hover:bg-green-100 border border-green-200" 
                    : "bg-gray-100 hover:bg-gray-200"}`}
                onClick={() => onSelectAudio(format.format_id)}
                whileHover={{ scale: 1.02 }}
                whileTap={{ scale: 0.98 }}
              >
                <div className="flex items-center">
                  <Music className={`w-5 h-5 mr-3 ${format.is_compatible ? "text-green-500" : "text-blue-500"}`} />
                  <div className="text-left">
                    <div className="flex items-center">
                      <p className="font-medium">{format.quality_label || format.format_note}</p>
                      {format.is_compatible && (
                        <span className="ml-2 bg-green-100 text-green-800 text-xs px-2 py-0.5 rounded-full">
                          Compatible
                        </span>
                      )}
                    </div>
                    <p className="text-sm text-gray-600">
                      {format.acodec} • {format.ext.toUpperCase()} • {formatFileSize(format.filesize)}
                    </p>
                  </div>
                </div>
                <span className={format.is_compatible ? "text-green-500" : "text-blue-500"}>Select</span>
              </motion.button>
            ))}
          </div>
        </motion.div>
      </motion.div>
    )}
  </AnimatePresence>
);

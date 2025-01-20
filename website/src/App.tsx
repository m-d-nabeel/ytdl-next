import { useState } from 'react';

interface Format {
  format_note: string;
  resolution?: string;
  ext: string;
  format_id: string;
}

const App = () => {
  const [mediaUrl, setMediaUrl] = useState('');
  const [status, setStatus] = useState('');
  const [qualities, setQualities] = useState<Format[]>([]);

  // Handle URL submission
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!mediaUrl.trim()) return;

    setStatus('Fetching qualities...');
    setQualities([]);

    try {
      const response = await fetch(`/api/yt/info?url=${encodeURIComponent(mediaUrl)}`);
      if (!response.ok) throw new Error('Failed to fetch qualities.');

      const data = await response.json();
      setQualities(data.formats);
      setStatus('Select a quality:');
    } catch (error: unknown) {
      if (error instanceof Error) {
        setStatus(`Error: ${error.message}`);
      } else {
        setStatus('An unknown error occurred.');
      }
    }
  };

  // Handle download

  const handleDownload = (formatId: string) => {
    const downloadUrl = `/api/yt/download?format_id=${formatId}&url=${encodeURIComponent(mediaUrl)}`;

    const link = document.createElement('a');
    link.href = downloadUrl;
    link.download = "media-download";
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    setStatus('Download started. Check your browser downloads.');
  };

  return (
    <div className="container mx-auto p-6">
      <h1 className="text-4xl font-semibold text-center mb-6">Media Downloader</h1>

      <form onSubmit={handleSubmit} className="max-w-xl mx-auto bg-white p-6 rounded-lg shadow-md">
        <input
          type="url"
          id="mediaUrl"
          placeholder="Enter media URL here..."
          className="w-full px-4 py-2 mb-4 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          value={mediaUrl}
          onChange={(e) => setMediaUrl(e.target.value)}
          required
        />
        <button
          type="submit"
          className="w-full bg-blue-500 text-white py-2 rounded-md hover:bg-blue-600 focus:ring-2 focus:ring-blue-500"
        >
          Fetch Qualities
        </button>
      </form>

      <div className="mt-4 text-center">{status}</div>

      {qualities.length > 0 && (
        <div className="mt-6 grid grid-cols-1 gap-4">
          {qualities.map(({ format_note, resolution, ext, format_id }) => (
            <button
              key={format_id}
              onClick={() => handleDownload(format_id)}
              className="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600"
            >
              {resolution || 'Audio only'} - {format_note} ({ext.toUpperCase()})
            </button>
          ))}
        </div>
      )}
    </div>
  );
};

export default App;

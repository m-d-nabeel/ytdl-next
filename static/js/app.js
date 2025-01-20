const form = document.getElementById('downloadForm');
const urlInput = document.getElementById('mediaUrl');
const downloadBtn = document.getElementById('downloadBtn');
const status = document.getElementById('status');

form.addEventListener('submit', async (e) => {
  e.preventDefault();

  const mediaURL = urlInput.value.trim();
  if (!mediaURL) return;

  try {
    downloadBtn.disabled = true;
    status.className = 'status downloading';
    status.textContent = 'Starting download...';

    // Create form data with the URL
    const formData = new FormData();
    formData.append('url', mediaURL);

    // Create the download link
    const downloadLink = document.createElement('a');
    downloadLink.href = `/api/download?url=${encodeURIComponent(mediaURL)}`;
    document.body.appendChild(downloadLink);
    downloadLink.click();
    document.body.removeChild(downloadLink);

    status.textContent = 'Download started. Check your browser downloads.';

  } catch (error) {
    status.className = 'status error';
    status.textContent = `Error: ${error.message}`;
  } finally {
    downloadBtn.disabled = false;
  }
});

export interface Format {
    format_id: string;
    resolution?: string;
    ext: string;
    filesize: number;
    format_note: string;
    acodec: string;
    vcodec?: string;
}

export interface VideoData {
    id: string;
    title: string;
    duration: number;
    uploader: string;
    view_count: number;
    like_count: number;
    formats: Format[];
}

export const formatFileSize = (bytes: number) => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
};

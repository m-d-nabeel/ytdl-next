export type YTVideoDetail = {
  title: string;
  embed:
    | {
        url: string;
        width: number;
        height: number;
      }
    | undefined;
  thumbnail_url: string[];
  quality: string;
};

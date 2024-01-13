import { create } from "zustand";

type Store = {
  videoDetails: YTVideoDetail | null;
  setVideoDetails: (videoDetails: YTVideoDetail) => void;
};

export const useStore = create<Store>()((set) => ({
  videoDetails: null,
  setVideoDetails: (videoDetails) => set({ videoDetails }),
}));

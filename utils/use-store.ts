import { create } from "zustand";

type Store = {
  info: string;
  title: string;
  setInfo: (info: string) => void;
  setTitle: (title: string) => void;
};

export const useStore = create<Store>()((set) => ({
  info: "",
  title: "",
  setInfo: (_info: string) => set((state) => ({ ...state, info: _info })),
  setTitle: (_title: string) => set((state) => ({ ...state, title: _title })),
}));

import { create } from "zustand";
import type { User } from "@/types/api";
import { authApi, ApiError } from "@/lib/api";

interface AuthState {
  me: User | null;
  status: "idle" | "loading" | "ready";
  bootstrap: () => Promise<void>;
  refreshMe: () => Promise<void>;
  setMe: (u: User | null) => void;
  logout: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set) => ({
  me: null,
  status: "idle",
  setMe: (u) => set({ me: u }),
  bootstrap: async () => {
    set({ status: "loading" });
    try {
      const me = await authApi.me();
      set({ me, status: "ready" });
    } catch (e) {
      if (e instanceof ApiError && (e.status === 401 || e.status === 404)) {
        set({ me: null, status: "ready" });
      } else {
        set({ me: null, status: "ready" });
      }
    }
  },
  refreshMe: async () => {
    try {
      const me = await authApi.me();
      set({ me });
    } catch {
      set({ me: null });
    }
  },
  logout: async () => {
    try { await authApi.logout(); } catch { /* ignore */ }
    set({ me: null });
  },
}));

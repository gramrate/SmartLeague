import { create } from "zustand";
import { getMe } from "../api/users";
import { ClubState, Role } from "../types/enums";

type AuthState = {
  isReady: boolean;
  isAuthed: boolean;
  userId?: string;
  role?: Role;
  email?: string;
  clubId?: string | null;
  clubState?: ClubState;
  nickname?: string;
  name?: string;
  showName?: boolean;
  description?: string | null;
  init: () => Promise<void>;
  clear: () => void;
};

export const useAuthStore = create<AuthState>((set) => ({
  isReady: false,
  isAuthed: false,
  init: async () => {
    try {
      const me = await getMe();
      set({
        isReady: true,
        isAuthed: true,
        userId: me.id,
        role: me.role,
        email: me.email,
        clubId: me.club_id ?? null,
        clubState: me.club_state,
        nickname: me.nickname,
        name: me.name,
        showName: me.show_name,
        description: me.description
      });
    } catch {
      set({ isReady: true, isAuthed: false, userId: undefined, role: undefined, email: undefined, clubId: undefined, clubState: undefined, nickname: undefined, name: undefined, showName: undefined, description: undefined });
    }
  },
  clear: () => set({ isReady: true, isAuthed: false, userId: undefined, role: undefined, email: undefined, clubId: undefined, clubState: undefined, nickname: undefined, name: undefined, showName: undefined, description: undefined })
}));

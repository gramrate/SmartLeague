import { ClubState, GameStatus } from "@/types/api";

export const CLUB_STATE_LABEL: Record<ClubState, string> = {
  [ClubState.None]: "—",
  [ClubState.Member]: "Участник",
  [ClubState.Resident]: "Резидент",
  [ClubState.Leader]: "Лидер",
  [ClubState.President]: "Президент",
};

export const CLUB_STATE_TONE: Record<ClubState, string> = {
  [ClubState.None]: "bg-muted text-muted-foreground",
  [ClubState.Member]: "bg-secondary text-secondary-foreground",
  [ClubState.Resident]: "bg-accent/20 text-accent",
  [ClubState.Leader]: "bg-primary/25 text-primary-foreground border border-primary/50",
  [ClubState.President]: "bg-gradient-to-r from-primary to-primary-glow text-primary-foreground shadow-[var(--shadow-glow)]",
};

export function isClubManager(state?: ClubState | null) {
  if (state == null) return false;
  return state === ClubState.Leader || state === ClubState.President;
}

export function canManageClub(
  me: { club_id?: string | null; club_state?: ClubState | null } | null,
  clubId: string,
) {
  return !!me && me.club_id === clubId && isClubManager(me.club_state ?? undefined);
}

export const GAME_STATUS_LABEL: Record<GameStatus, string> = {
  [GameStatus.Draft]: "Черновик",
  [GameStatus.InProgress]: "Идет",
  [GameStatus.Finished]: "Завершена",
};

export function displayUserName(u: {
  nickname?: string;
  name?: string;
  show_name?: boolean;
  id?: string;
} | null | undefined): string {
  if (!u) return "Неизвестно";
  if (u.nickname) return u.nickname;
  if (u.show_name && u.name) return u.name;
  if (u.name) return u.name;
  return "Неизвестный игрок";
}

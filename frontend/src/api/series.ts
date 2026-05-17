import { apiFetch } from "./client";
import type { Game, LeaderboardRow, PaginationInfo, Series, User, UUID } from "../types/dto";

export interface CreateSeriesRequest {
  name: string;
  description: string;
  start_at: string;
  end_at: string;
  price_rub: number;
  is_closed: boolean;
  game_type: number;
}

export function createSeries(req: CreateSeriesRequest) {
  return apiFetch<Series>("/api/v1/series", { method: "POST", body: JSON.stringify(req) });
}

export function getSeries(id: UUID) {
  return apiFetch<Series>(`/api/v1/series/${id}`, { method: "GET" });
}

export interface GetSeriesFullResponse {
  series: Series;
  participants: GetSeriesParticipantsResponse;
  games: {
    items: Game[];
    pagination: PaginationInfo;
  };
  leaderboard: GetSeriesLeaderboardResponse;
}

export function getSeriesFull(
  id: UUID,
  params: {
    participants_limit?: number;
    participants_offset?: number;
    games_limit?: number;
    games_offset?: number;
    leaderboard_limit?: number;
    leaderboard_offset?: number;
  } = {}
) {
  const q = new URLSearchParams();
  if (params.participants_limit != null) q.set("participants_limit", String(params.participants_limit));
  if (params.participants_offset != null) q.set("participants_offset", String(params.participants_offset));
  if (params.games_limit != null) q.set("games_limit", String(params.games_limit));
  if (params.games_offset != null) q.set("games_offset", String(params.games_offset));
  if (params.leaderboard_limit != null) q.set("leaderboard_limit", String(params.leaderboard_limit));
  if (params.leaderboard_offset != null) q.set("leaderboard_offset", String(params.leaderboard_offset));
  return apiFetch<GetSeriesFullResponse>(`/api/v1/series/${id}/full?${q.toString()}`, { method: "GET" });
}

export function updateSeries(id: UUID, patch: Partial<CreateSeriesRequest>) {
  return apiFetch<Series>(`/api/v1/series/${id}`, { method: "PATCH", body: JSON.stringify(patch) });
}

export function deleteSeries(id: UUID) {
  return apiFetch<void>(`/api/v1/series/${id}`, { method: "DELETE" });
}

export interface GetClubSeriesResponse {
  items: Series[];
  pagination: PaginationInfo;
}

export function getClubSeries(clubId: UUID, params: { limit?: number; offset?: number }) {
  const q = new URLSearchParams();
  if (params.limit != null) q.set("limit", String(params.limit));
  if (params.offset != null) q.set("offset", String(params.offset));
  return apiFetch<GetClubSeriesResponse>(`/api/v1/club/${clubId}/series?${q.toString()}`, { method: "GET" });
}

export interface AllSeriesItem {
  id: UUID;
  club_id: UUID;
  club_name: string;
  name: string;
  description: string;
  start_at: string;
  end_at: string;
  games_count: number;
}

export interface GetAllSeriesResponse {
  items: AllSeriesItem[];
  pagination: PaginationInfo;
}

export function listAllSeries(params: { limit?: number; offset?: number }) {
  const q = new URLSearchParams();
  if (params.limit != null) q.set("limit", String(params.limit));
  if (params.offset != null) q.set("offset", String(params.offset));
  return apiFetch<GetAllSeriesResponse>(`/api/v1/series/all?${q.toString()}`, { method: "GET" });
}

export function joinSeries(id: UUID) {
  return apiFetch<void>(`/api/v1/series/${id}/join`, { method: "POST" });
}

export function leaveSeries(id: UUID) {
  return apiFetch<void>(`/api/v1/series/${id}/leave`, { method: "POST" });
}

export interface GetSeriesParticipantsResponse {
  items: User[];
  pagination: PaginationInfo;
}

export function getSeriesParticipants(id: UUID, params: { limit?: number; offset?: number; q?: string }) {
  const q = new URLSearchParams();
  if (params.limit != null) q.set("limit", String(params.limit));
  if (params.offset != null) q.set("offset", String(params.offset));
  if (params.q) q.set("q", params.q);
  return apiFetch<GetSeriesParticipantsResponse>(`/api/v1/series/${id}/participants?${q.toString()}`, { method: "GET" });
}

export interface GetSeriesLeaderboardResponse {
  items: LeaderboardRow[];
  pagination: PaginationInfo;
}

export function getLeaderboard(id: UUID, params: { limit?: number; offset?: number }) {
  const q = new URLSearchParams();
  if (params.limit != null) q.set("limit", String(params.limit));
  if (params.offset != null) q.set("offset", String(params.offset));
  return apiFetch<GetSeriesLeaderboardResponse>(`/api/v1/series/${id}/leaderboard?${q.toString()}`, { method: "GET" });
}

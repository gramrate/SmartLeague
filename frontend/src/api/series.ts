import { apiFetch } from "./client";
import type { LeaderboardRow, PaginationInfo, Series, User, UUID } from "../types/dto";

export interface CreateSeriesRequest {
  name: string;
  scoring_rules: string;
  start_at: string;
  end_at: string;
  description?: string | null;
  price_rub: number;
  is_closed: boolean;
  game_type: number;
  status: number;
}

export function createSeries(req: CreateSeriesRequest) {
  return apiFetch<Series>("/api/v1/series", { method: "POST", body: JSON.stringify(req) });
}

export function getSeries(id: UUID) {
  return apiFetch<Series>(`/api/v1/series/${id}`, { method: "GET" });
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

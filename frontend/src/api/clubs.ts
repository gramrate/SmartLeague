import { apiFetch } from "./client";
import type { Club, PaginationInfo, User, UUID } from "../types/dto";
import { ClubState } from "../types/enums";

export interface CreateClubRequest {
  name: string;
  description?: string | null;
}

export function createClub(req: CreateClubRequest) {
  return apiFetch<Club>("/api/v1/club", { method: "POST", body: JSON.stringify(req) });
}

export function getClub(id: UUID) {
  return apiFetch<Club>(`/api/v1/club/${id}`, { method: "GET" });
}

export interface GetAllClubsResponse {
  items: Club[];
  pagination: PaginationInfo;
}

export function listClubs(params: { limit?: number; offset?: number }) {
  const q = new URLSearchParams();
  if (params.limit != null) q.set("limit", String(params.limit));
  if (params.offset != null) q.set("offset", String(params.offset));
  return apiFetch<GetAllClubsResponse>(`/api/v1/club/all?${q.toString()}`, { method: "GET" });
}

export interface GetClubMembersResponse {
  items: User[];
  pagination: PaginationInfo;
}

export function getClubMembers(id: UUID, params: { limit?: number; offset?: number }) {
  const q = new URLSearchParams();
  if (params.limit != null) q.set("limit", String(params.limit));
  if (params.offset != null) q.set("offset", String(params.offset));
  return apiFetch<GetClubMembersResponse>(`/api/v1/club/${id}/members?${q.toString()}`, { method: "GET" });
}

export function joinClub(id: UUID) {
  return apiFetch<void>(`/api/v1/club/${id}/join`, { method: "POST" });
}

export function leaveClub() {
  return apiFetch<void>(`/api/v1/club/leave`, { method: "POST" });
}

export function setLeader(clubId: UUID, memberId: UUID) {
  return apiFetch<void>(`/api/v1/club/${clubId}/leader/${memberId}`, { method: "POST" });
}

export function setClubMemberRole(clubId: UUID, memberId: UUID, state: ClubState) {
  return apiFetch<void>(`/api/v1/club/${clubId}/member/${memberId}/role`, { method: "POST", body: JSON.stringify({ state }) });
}

export function kickClubMember(clubId: UUID, memberId: UUID) {
  return apiFetch<void>(`/api/v1/club/${clubId}/member/${memberId}/kick`, { method: "POST" });
}

export function blockClubMember(clubId: UUID, memberId: UUID) {
  return apiFetch<void>(`/api/v1/club/${clubId}/member/${memberId}/block`, { method: "POST" });
}

export function updateClub(id: UUID, patch: Partial<CreateClubRequest>) {
  return apiFetch<Club>(`/api/v1/club/${id}`, { method: "PATCH", body: JSON.stringify(patch) });
}

export function deleteClub(id: UUID) {
  return apiFetch<void>(`/api/v1/club/${id}`, { method: "DELETE" });
}

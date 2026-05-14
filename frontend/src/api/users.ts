import { apiFetch } from "./client";
import type { PaginationInfo, User } from "../types/dto";
import type { Role } from "../types/enums";

export interface GetMeResponse {
  id: string;
  nickname: string;
  email: string;
  password_hash?: string;
  name: string;
  show_name: boolean;
  description?: string | null;
  club_id?: string | null;
  club_state: number;
  role: Role;
}

export function getMe() {
  return apiFetch<GetMeResponse>("/api/v1/user", { method: "GET" });
}

export function getUserById(id: string) {
  return apiFetch<GetMeResponse>(`/api/v1/user/${id}`, { method: "GET" });
}

export function updateMe(patch: { nickname?: string; name?: string; show_name?: boolean; description?: string | null }) {
  return apiFetch<GetMeResponse>("/api/v1/user", { method: "PATCH", body: JSON.stringify(patch) });
}

export function changePassword(req: { old_password: string; new_password: string }) {
  return apiFetch<void>("/api/v1/user/password", { method: "POST", body: JSON.stringify(req) });
}

export interface GetAllUsersResponse {
  items: User[];
  pagination: PaginationInfo;
}

export function getUsers(params: { limit?: number; offset?: number; role?: number; q?: string; email_prefix?: string }) {
  const q = new URLSearchParams();
  if (params.limit != null) q.set("limit", String(params.limit));
  if (params.offset != null) q.set("offset", String(params.offset));
  if (params.role != null) q.set("role", String(params.role));
  if (params.q) q.set("q", params.q);
  if (params.email_prefix) q.set("email_prefix", params.email_prefix);
  return apiFetch<GetAllUsersResponse>(`/api/v1/user/all?${q.toString()}`, { method: "GET" });
}

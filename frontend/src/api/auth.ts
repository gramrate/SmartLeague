import { apiFetch } from "./client";
import type { User } from "../types/dto";

export interface RegisterUserRequest {
  nickname?: string;
  name: string;
  show_name?: boolean;
  description?: string | null;
  email: string;
  password: string;
}
export interface RegisterUserResponse extends User {
  refresh_token: string;
}

export interface LoginUserRequest {
  email: string;
  password: string;
}
export interface LoginUserResponse {
  refresh_token: string;
}

export function registerUser(req: RegisterUserRequest) {
  return apiFetch<RegisterUserResponse>("/api/v1/user/register", {
    method: "POST",
    body: JSON.stringify(req)
  });
}

export function loginUser(req: LoginUserRequest) {
  return apiFetch<LoginUserResponse>("/api/v1/user/login", {
    method: "POST",
    body: JSON.stringify(req)
  });
}

export function logoutUser() {
  return apiFetch<void>("/api/v1/user/logout", { method: "POST" });
}

export function refresh() {
  return apiFetch<void>("/api/v1/auth/refresh", { method: "POST", retryOn401: false });
}

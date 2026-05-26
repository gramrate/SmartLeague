// Unified API client for SmartLeague backend.
// Cookie-based auth: requests use credentials:"include". On 401 we try refresh once.

import type {
  Club, Paged, User, Series, AllSeriesItem, Game, GameFull, LeaderboardRow,
  SeriesFull, PlayerGame, PlayerSeries, LoginRequest, RegisterRequest,
  CreateClubRequest, UpdateClubRequest, CreateSeriesRequest, CreateGameRequest,
  UpdateCurrentUserRequest, GameResultRow, UpdateSeriesRequest, ManageGameRow,
} from "@/types/api";

export const API_BASE_URL =
  (import.meta as any).env?.VITE_API_BASE_URL || "http://localhost:8000";

export class ApiError extends Error {
  status: number;
  body: unknown;
  constructor(status: number, message: string, body?: unknown) {
    super(message);
    this.status = status;
    this.body = body;
  }
}

interface FetchOpts extends Omit<RequestInit, "body"> {
  body?: unknown;
  query?: Record<string, string | number | boolean | undefined | null>;
  // internal — prevent infinite refresh loop
  _retried?: boolean;
  // skip auto-refresh (e.g. for the /user "me" probe)
  skipRefresh?: boolean;
}

function buildUrl(path: string, query?: FetchOpts["query"]) {
  const url = new URL(path.startsWith("http") ? path : `${API_BASE_URL}${path}`);
  if (query) {
    for (const [k, v] of Object.entries(query)) {
      if (v !== undefined && v !== null && v !== "") {
        url.searchParams.set(k, String(v));
      }
    }
  }
  return url.toString();
}

async function refreshTokens(): Promise<boolean> {
  try {
    const res = await fetch(`${API_BASE_URL}/api/v1/auth/refresh`, {
      method: "POST",
      credentials: "include",
    });
    return res.ok;
  } catch {
    return false;
  }
}

export async function apiFetch<T = unknown>(
  path: string,
  opts: FetchOpts = {},
): Promise<T> {
  const { body, query, headers, _retried, skipRefresh, ...rest } = opts;
  const init: RequestInit = {
    ...rest,
    credentials: "include",
    headers: {
      Accept: "application/json",
      ...(body !== undefined ? { "Content-Type": "application/json" } : {}),
      ...(headers as Record<string, string> | undefined),
    },
    body: body !== undefined ? JSON.stringify(body) : undefined,
  };

  let res: Response;
  try {
    res = await fetch(buildUrl(path, query), init);
  } catch (e: any) {
    throw new ApiError(0, `Network error: ${e?.message ?? "unreachable"}`);
  }

  if (res.status === 401 && !_retried && !skipRefresh) {
    const ok = await refreshTokens();
    if (ok) {
      return apiFetch<T>(path, { ...opts, _retried: true });
    }
  }

  if (res.status === 204) return undefined as T;

  let payload: any = null;
  const text = await res.text();
  if (text) {
    try { payload = JSON.parse(text); } catch { payload = text; }
  }

  if (!res.ok) {
    const message =
      (payload && typeof payload === "object" && (payload.message || payload.error)) ||
      (typeof payload === "string" && payload) ||
      `HTTP ${res.status}`;
    throw new ApiError(res.status, String(message), payload);
  }

  return payload as T;
}

// ---------- Auth ----------
export const authApi = {
  me: () => apiFetch<User>("/api/v1/user", { skipRefresh: false }),
  login: (data: LoginRequest) =>
    apiFetch<User>("/api/v1/user/login", { method: "POST", body: data, skipRefresh: true }),
  register: (data: RegisterRequest) =>
    apiFetch<User>("/api/v1/user/register", { method: "POST", body: data, skipRefresh: true }),
  logout: () => apiFetch<void>("/api/v1/user/logout", { method: "POST" }),
  updateMe: (data: UpdateCurrentUserRequest) =>
    apiFetch<User>("/api/v1/user", { method: "PATCH", body: data }),
  changePassword: (old_password: string, new_password: string) =>
    apiFetch<void>("/api/v1/user/password", {
      method: "POST",
      body: { old_password, new_password },
    }),
};

// ---------- Users ----------
export const usersApi = {
  get: (id: string) => apiFetch<User>(`/api/v1/user/${id}`),
  search: (params: { q?: string; club?: string; club_state?: number; limit?: number; offset?: number }) =>
    apiFetch<Paged<User>>("/api/v1/user/all", { query: params }),
  games: (id: string, limit = 50, offset = 0) =>
    apiFetch<Paged<PlayerGame>>(`/api/v1/user/${id}/games`, { query: { limit, offset } }),
  series: (id: string, params?: {
    q?: string;
    from?: string;
    to?: string;
    is_rating?: boolean;
    show_past?: boolean;
    show_closed?: boolean;
    limit?: number;
    offset?: number;
  }) =>
    apiFetch<Paged<PlayerSeries>>(`/api/v1/user/${id}/series`, { query: params }),
};

// ---------- Clubs ----------
export const clubsApi = {
  all: (params?: { q?: string; limit?: number; offset?: number }) =>
    apiFetch<Paged<Club>>("/api/v1/club/all", { query: params }),
  get: (id: string) => apiFetch<Club>(`/api/v1/club/${id}`),
  create: (data: CreateClubRequest) =>
    apiFetch<Club>("/api/v1/club", { method: "POST", body: data }),
  update: (id: string, data: UpdateClubRequest) =>
    apiFetch<Club>(`/api/v1/club/${id}`, { method: "PATCH", body: data }),
  delete: (id: string) => apiFetch<void>(`/api/v1/club/${id}`, { method: "DELETE" }),
  members: (id: string, params?: { q?: string; club_state?: number; limit?: number; offset?: number }) =>
    apiFetch<Paged<User>>(`/api/v1/club/${id}/members`, { query: params }),
  games: (id: string, params?: { limit?: number; offset?: number }) =>
    apiFetch<Paged<PlayerGame>>(`/api/v1/club/${id}/games`, { query: params }),
  series: (id: string, limit = 100, offset = 0) =>
    apiFetch<Paged<Series>>(`/api/v1/club/${id}/series`, { query: { limit, offset } }),
  join: (id: string) =>
    apiFetch<void>(`/api/v1/club/${id}/join`, { method: "POST" }),
  leave: () => apiFetch<void>("/api/v1/club/leave", { method: "POST" }),
  setRole: (id: string, memberId: string, state: number) =>
    apiFetch<void>(`/api/v1/club/${id}/member/${memberId}/role`, {
      method: "POST", body: { state },
    }),
  setLeader: (id: string, memberId: string) =>
    apiFetch<void>(`/api/v1/club/${id}/leader/${memberId}`, { method: "POST" }),
  kick: (id: string, memberId: string) =>
    apiFetch<void>(`/api/v1/club/${id}/member/${memberId}/kick`, { method: "POST" }),
  block: (id: string, memberId: string) =>
    apiFetch<void>(`/api/v1/club/${id}/member/${memberId}/block`, { method: "POST" }),
};

// ---------- Series ----------
export const seriesApi = {
  all: (params?: {
    q?: string;
    club?: string;
    from?: string;
    to?: string;
    is_rating?: boolean;
    show_past?: boolean;
    show_closed?: boolean;
    limit?: number;
    offset?: number;
  }) =>
    apiFetch<Paged<AllSeriesItem>>("/api/v1/series/all", { query: params }),
  get: (id: string) => apiFetch<Series>(`/api/v1/series/${id}`),
  full: (id: string) => apiFetch<SeriesFull>(`/api/v1/series/${id}/full`),
  create: (clubId: string, data: CreateSeriesRequest) => {
    // Series are created under a club — backend infers club from creator's club_id
    // (per swagger: POST /api/v1/series); clubId is currently unused but kept for clarity.
    void clubId;
    return apiFetch<Series>("/api/v1/series", { method: "POST", body: data });
  },
  update: (id: string, data: UpdateSeriesRequest) =>
    apiFetch<Series>(`/api/v1/series/${id}`, { method: "PATCH", body: data }),
  delete: (id: string) => apiFetch<void>(`/api/v1/series/${id}`, { method: "DELETE" }),
  games: (id: string, limit = 100, offset = 0) =>
    apiFetch<Paged<Game>>(`/api/v1/series/${id}/games`, { query: { limit, offset } }),
  participants: (id: string, limit = 100, offset = 0) =>
    apiFetch<Paged<User>>(`/api/v1/series/${id}/participants`, { query: { limit, offset } }),
  leaderboard: (id: string, limit = 100, offset = 0) =>
    apiFetch<Paged<LeaderboardRow>>(`/api/v1/series/${id}/leaderboard`, { query: { limit, offset } }),
  join: (id: string) => apiFetch<void>(`/api/v1/series/${id}/join`, { method: "POST" }),
  leave: (id: string) => apiFetch<void>(`/api/v1/series/${id}/leave`, { method: "POST" }),
  createGame: (seriesId: string, data: CreateGameRequest) =>
    apiFetch<Game>(`/api/v1/series/${seriesId}/games`, { method: "POST", body: data }),
  createGameDraft: (seriesId: string, data: CreateGameRequest) =>
    apiFetch<Game>(`/api/v1/series/${seriesId}/games/draft`, { method: "POST", body: data }),
};

// ---------- Games ----------
export const gamesApi = {
  get: (id: string) => apiFetch<Game>(`/api/v1/game/${id}`),
  full: (id: string) => apiFetch<GameFull>(`/api/v1/game/${id}/full`),
  delete: (id: string) => apiFetch<void>(`/api/v1/game/${id}`, { method: "DELETE" }),
  setParticipants: (id: string, participant_ids: string[]) =>
    apiFetch<void>(`/api/v1/game/${id}/participants`, {
      method: "POST", body: { participant_ids },
    }),
  setResults: (id: string, rows: GameResultRow[]) =>
    apiFetch<void>(`/api/v1/game/${id}/results`, { method: "POST", body: { rows } }),
  saveDraft: (id: string, rows: ManageGameRow[]) =>
    apiFetch<void>(`/api/v1/game/${id}/draft`, { method: "POST", body: { rows } }),
  publish: (id: string, rows: ManageGameRow[]) =>
    apiFetch<void>(`/api/v1/game/${id}/publish`, { method: "POST", body: { rows } }),
};

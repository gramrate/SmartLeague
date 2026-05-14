import { queryClient } from "../shared/queryClient";

type RequestOptions = Omit<RequestInit, "headers"> & {
  headers?: Record<string, string>;
  retryOn401?: boolean;
};

const baseUrl = import.meta.env.VITE_API_BASE_URL || "";

async function refreshTokens(): Promise<boolean> {
  try {
    const resp = await fetch(`${baseUrl}/api/v1/auth/refresh`, {
      method: "POST",
      credentials: "include"
    });
    return resp.status === 204;
  } catch {
    return false;
  }
}

export async function apiFetch<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const { retryOn401 = true, headers, ...rest } = opts;

  const resp = await fetch(`${baseUrl}${path}`, {
    credentials: "include",
    ...rest,
    headers: {
      "Content-Type": "application/json",
      ...(headers || {})
    }
  });

  if (resp.status === 204) {
    return undefined as unknown as T;
  }

  if (resp.status === 401 && retryOn401) {
    const ok = await refreshTokens();
    if (ok) {
      queryClient.invalidateQueries();
      return apiFetch<T>(path, { ...opts, retryOn401: false });
    }
  }

  if (!resp.ok) {
    const body = await resp.json().catch(() => null);
    throw body || { code: resp.status, message: resp.statusText };
  }

  return (await resp.json()) as T;
}


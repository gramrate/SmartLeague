import { ApiError } from "./api";

const translations: Record<string, string> = {
  "invalid email or password": "Неверный email или пароль",
  "password mismatch": "Старый пароль неверен",
  "new password should differ from the old": "Новый пароль должен отличаться от текущего",
  "permission denied": "Нет доступа",
  "profile not found": "Пользователь не найден",
  "invalid request": "Некорректный запрос",
};

export function translateError(e: unknown, fallback = "Ошибка"): string {
  const msg = e instanceof ApiError ? e.message : fallback;
  return translations[msg] ?? msg;
}

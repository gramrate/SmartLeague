import { createFileRoute, useNavigate, Link } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { gamesApi, seriesApi, ApiError } from "@/lib/api";
import { LoadingBlock, ErrorBlock } from "@/components/site/States";
import { canManageClub, displayUserName } from "@/lib/roles";
import { useAuthStore } from "@/lib/auth-store";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import type { MafiaRole, ManageGameRow } from "@/types/api";
import { normalizeSearchText } from "@/lib/search";

export const Route = createFileRoute("/game/$id/manage")({ component: ManageGamePage });

type RowState = {
  slot: number;
  nickname: string;
  profile_id?: string;
  role?: MafiaRole;
  best_move?: string;
  yellow_cards: string;
  removed: string;
  extra_points: string;
  total_points: string;
};

type PersistedDraft = {
  rows: RowState[];
};

const ROLE_OPTIONS: { value: MafiaRole; label: string }[] = [
  { value: "civilian", label: "Мирный" },
  { value: "mafia", label: "Мафия" },
  { value: "don", label: "Дон" },
  { value: "sheriff", label: "Шериф" },
];

function toNum(value: string): number {
  const normalized = value.trim().replace(/\s+/g, "").replace(/,/g, ".");
  const n = Number(normalized);
  return Number.isFinite(n) ? n : 0;
}

function normalizeBestMove(raw: string | undefined): string | undefined {
  if (!raw) return undefined;
  const clean = raw.replace(/[^0-9]/g, "").slice(0, 3);
  return clean.length ? clean : undefined;
}

function ManageGamePage() {
  const { id } = Route.useParams();
  const me = useAuthStore((s) => s.me);
  const status = useAuthStore((s) => s.status);
  const navigate = useNavigate();
  const qc = useQueryClient();

  const game = useQuery({ queryKey: ["game", id, "full"], queryFn: () => gamesApi.full(id) });
  const series = useQuery({
    queryKey: ["series", game.data?.series_id],
    queryFn: () => seriesApi.get(game.data!.series_id),
    enabled: !!game.data?.series_id,
  });
  const participants = useQuery({
    queryKey: ["series", game.data?.series_id, "participants"],
    queryFn: () => seriesApi.participants(game.data!.series_id),
    enabled: !!game.data?.series_id,
  });

  const canManage = series.data ? canManageClub(me, series.data.club_id) : null;
  useEffect(() => {
    if (status === "ready" && canManage === false) navigate({ to: "/game/$id", params: { id } });
  }, [canManage, status, navigate, id]);

  const [rows, setRows] = useState<RowState[]>([]);
  const [activeSlot, setActiveSlot] = useState<number | null>(null);
  const [savingDraft, setSavingDraft] = useState(false);
  const [publishing, setPublishing] = useState(false);

  const draftStorageKey = `smartleague:game-manage-draft:${id}`;

  useEffect(() => {
    if (!game.data) return;

    const byProfile = new Map((game.data.results ?? []).map((r) => [r.profile_id, r]));
    const pids = game.data.participant_ids ?? [];

    const init: RowState[] = Array.from({ length: 10 }, (_, i) => {
      const slot = i + 1;
      const pid = pids[i];
      const rr = pid ? byProfile.get(pid) : undefined;
      return {
        slot,
        nickname: pid ? (participants.data?.items ?? []).find((p) => p.id === pid)?.nickname
          || (participants.data?.items ?? []).find((p) => p.id === pid)?.name
          || ""
          : "",
        profile_id: pid,
        role: rr?.role ?? "civilian",
        best_move: rr?.best_move ?? "",
        yellow_cards: rr ? String(rr.yellow_cards ?? 0) : "0",
        removed: rr ? String(rr.removed ?? 0) : "0",
        extra_points: rr ? String(rr.extra_points ?? 0) : "0",
        total_points: rr ? String(rr.total_points ?? 0) : "0",
      };
    });

    if (typeof window === "undefined") {
      setRows(init);
      return;
    }

    let persisted: PersistedDraft | null = null;
    try {
      const raw = window.localStorage.getItem(draftStorageKey);
      if (raw) persisted = JSON.parse(raw) as PersistedDraft;
    } catch {
      persisted = null;
    }

    if (!persisted?.rows?.length) {
      setRows(init);
      return;
    }

    const persistedBySlot = new Map(persisted.rows.map((r) => [r.slot, r]));
    const merged = init.map((r) => {
      const pr = persistedBySlot.get(r.slot);
      if (!pr) return r;
      const persistedNickname = pr.nickname ?? "";
      const hasManualNickname = persistedNickname.trim().length > 0;
      return {
        ...r,
        nickname: persistedNickname || r.nickname,
        // If nickname is manually entered without selecting a player, keep profile_id empty.
        profile_id: hasManualNickname ? pr.profile_id : (pr.profile_id ?? r.profile_id),
        role: pr.role ?? r.role,
        best_move: pr.best_move ?? r.best_move,
        yellow_cards: pr.yellow_cards ?? r.yellow_cards,
        removed: pr.removed ?? r.removed,
        extra_points: pr.extra_points ?? r.extra_points,
        total_points: pr.total_points ?? r.total_points,
      };
    });

    setRows(merged);
  }, [game.data, participants.data, draftStorageKey]);

  useEffect(() => {
    if (typeof window === "undefined") return;
    try {
      const payload: PersistedDraft = { rows };
      window.localStorage.setItem(draftStorageKey, JSON.stringify(payload));
    } catch {
      // ignore localStorage errors
    }
  }, [rows, draftStorageKey]);

  const setCell = (slot: number, patch: Partial<RowState>) => {
    setRows((prev) => prev.map((r) => (r.slot === slot ? { ...r, ...patch } : r)));
  };

  const onNicknameChange = (slot: number, raw: string) => {
    const value = raw.trimStart();
    const normalized = normalizeSearchText(value);
    const matched = (participants.data?.items ?? []).find((p) => normalizeSearchText(displayUserName(p)) === normalized);
    setCell(slot, {
      nickname: value,
      profile_id: matched?.id,
    });
  };

  const participantItems = participants.data?.items ?? [];

  const payloadRows: ManageGameRow[] = rows.map((r) => ({
    slot: r.slot,
    profile_id: r.profile_id,
    role: r.role,
    best_move: normalizeBestMove(r.best_move),
    yellow_cards: toNum(r.yellow_cards),
    removed: toNum(r.removed),
    extra_points: toNum(r.extra_points),
    total_points: toNum(r.total_points),
  }));

  const draftOnlySlots = rows
    .filter((r) => r.nickname.trim().length > 0 && !r.profile_id)
    .map((r) => r.slot);
  const slotsByProfile = rows.reduce<Map<string, number[]>>((acc, r) => {
    if (!r.profile_id) return acc;
    const prev = acc.get(r.profile_id) ?? [];
    prev.push(r.slot);
    acc.set(r.profile_id, prev);
    return acc;
  }, new Map<string, number[]>());
  const duplicateParticipantSlots = Array.from(slotsByProfile.values())
    .filter((slots) => slots.length > 1)
    .flat()
    .sort((a, b) => a - b);

  const saveDraft = async () => {
    if (duplicateParticipantSlots.length > 0) {
      toast.error(`Нельзя сохранить: один и тот же игрок выбран в нескольких слотах (${duplicateParticipantSlots.join(", ")}).`);
      return;
    }
    setSavingDraft(true);
    try {
      await gamesApi.saveDraft(id, payloadRows);
      toast.success("Черновик сохранен");
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : "Ошибка");
    } finally {
      setSavingDraft(false);
    }
  };

  const publish = async () => {
    if (duplicateParticipantSlots.length > 0) {
      toast.error(`Нельзя опубликовать: один и тот же игрок выбран в нескольких слотах (${duplicateParticipantSlots.join(", ")}).`);
      return;
    }
    if (draftOnlySlots.length > 0) {
      toast.error(`Нельзя опубликовать: в слотах ${draftOnlySlots.join(", ")} указан только текстовый никнейм. Выберите игрока из списка.`);
      return;
    }
    setPublishing(true);
    try {
      await gamesApi.publish(id, payloadRows);
      if (typeof window !== "undefined") {
        window.localStorage.removeItem(draftStorageKey);
      }
      qc.invalidateQueries({ queryKey: ["game", id, "full"] });
      toast.success("Игра опубликована");
      navigate({ to: "/game/$id", params: { id } });
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : "Ошибка");
    } finally {
      setPublishing(false);
    }
  };

  if (game.isLoading) return <PageShell><LoadingBlock /></PageShell>;
  if (game.error) return <PageShell><ErrorBlock error={game.error} /></PageShell>;
  if (!game.data) return null;

  return (
    <PageShell>
      <PageHeader
        eyebrow="Управление игрой"
        title={game.data.name || `Игра #${game.data.number}`}
        actions={<Button variant="outline" asChild><Link to="/game/$id" params={{ id }}>Назад к игре</Link></Button>}
      />

      <section className="rounded-2xl border border-border/60 bg-card/60 p-6">
        <div className="overflow-x-auto">
          <table className="w-full min-w-[1100px] text-sm">
            <thead>
              <tr className="border-b border-border/60 text-left text-xs uppercase text-muted-foreground">
                <th className="px-2 py-2">Слот</th>
                <th className="px-2 py-2">Никнейм</th>
                <th className="px-2 py-2">Роль</th>
                <th className="w-24 px-2 py-2">Лучший ход</th>
                <th className="px-2 py-2">ЖК</th>
                <th className="px-2 py-2">Удаление</th>
                <th className="px-2 py-2">Доп балл</th>
                <th className="px-2 py-2">Итоговый балл</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((r) => (
                <tr key={r.slot} className="border-b border-border/30">
                  <td className="px-2 py-2 font-semibold text-primary">{r.slot}</td>
                  <td className="px-2 py-2">
                    <div className="relative">
                      <Input
                        value={r.nickname}
                        onChange={(e) => onNicknameChange(r.slot, e.target.value)}
                        onFocus={() => setActiveSlot(r.slot)}
                        onBlur={() => setTimeout(() => setActiveSlot((s) => (s === r.slot ? null : s)), 120)}
                        placeholder="Введите или выберите"
                        className="h-8"
                      />
                      {activeSlot === r.slot && (
                        <div className="absolute z-20 mt-1 max-h-52 w-full overflow-auto rounded-md border border-border bg-popover p-1 shadow-md">
                          {participantItems
                            .filter((p) => {
                              const q = normalizeSearchText(r.nickname);
                              if (!q) return true;
                              return normalizeSearchText(displayUserName(p)).includes(q);
                            })
                            .map((p) => (
                              <button
                                type="button"
                                key={p.id}
                                className="block w-full rounded-sm px-2 py-1.5 text-left text-sm hover:bg-accent"
                                onMouseDown={(e) => e.preventDefault()}
                                onClick={() => {
                                  setCell(r.slot, { nickname: displayUserName(p), profile_id: p.id });
                                  setActiveSlot(null);
                                }}
                              >
                                {displayUserName(p)}
                              </button>
                            ))}
                        </div>
                      )}
                    </div>
                  </td>
                  <td className="px-2 py-2">
                    <Select value={r.role ?? ""} onValueChange={(v) => setCell(r.slot, { role: (v || undefined) as MafiaRole | undefined })}>
                      <SelectTrigger className="h-8">
                        <SelectValue placeholder="Роль" />
                      </SelectTrigger>
                      <SelectContent>
                        {ROLE_OPTIONS.map((opt) => (
                          <SelectItem key={opt.value} value={opt.value}>{opt.label}</SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </td>
                  <td className="w-24 px-2 py-2">
                    <Input
                      value={r.best_move ?? ""}
                      onChange={(e) => setCell(r.slot, { best_move: e.target.value.replace(/[^0-9]/g, "").slice(0, 3) })}
                      placeholder="до 3 цифр"
                      className="h-8"
                    />
                  </td>
                  <td className="px-2 py-2"><Input value={r.yellow_cards} onChange={(e) => setCell(r.slot, { yellow_cards: e.target.value })} className="h-8" /></td>
                  <td className="px-2 py-2"><Input value={r.removed} onChange={(e) => setCell(r.slot, { removed: e.target.value })} className="h-8" /></td>
                  <td className="px-2 py-2"><Input value={r.extra_points} onChange={(e) => setCell(r.slot, { extra_points: e.target.value })} className="h-8" /></td>
                  <td className="px-2 py-2"><Input value={r.total_points} onChange={(e) => setCell(r.slot, { total_points: e.target.value })} className="h-8" /></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div className="mt-5 flex flex-wrap gap-2">
          <Button variant="outline" disabled={savingDraft || publishing} onClick={() => void saveDraft()}>
            {savingDraft ? "Сохранение..." : "Сохранить как черновик"}
          </Button>
          <Button disabled={savingDraft || publishing || draftOnlySlots.length > 0 || duplicateParticipantSlots.length > 0} onClick={() => void publish()}>
            {publishing ? "Публикация..." : "Опубликовать"}
          </Button>
        </div>

        <p className="mt-3 text-xs text-muted-foreground">
          Черновики видны только лидерам и президентам клуба. После публикации игра становится доступной всем.
        </p>
        {draftOnlySlots.length > 0 && (
          <p className="mt-1 text-xs text-amber-600">
            Публикация недоступна: в слотах {draftOnlySlots.join(", ")} не выбран игрок из списка.
          </p>
        )}
        {duplicateParticipantSlots.length > 0 && (
          <p className="mt-1 text-xs text-destructive">
            Сохранение и публикация недоступны: один и тот же игрок выбран в нескольких слотах ({duplicateParticipantSlots.join(", ")}).
          </p>
        )}
      </section>
    </PageShell>
  );
}

import { createFileRoute, useNavigate, Link } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { gamesApi, seriesApi, ApiError } from "@/lib/api";
import { LoadingBlock, ErrorBlock } from "@/components/site/States";
import { canManageClub, displayUserName } from "@/lib/roles";
import { useAuthStore } from "@/lib/auth-store";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from "@/components/ui/alert-dialog";
import { useEffect, useRef, useState } from "react";
import { toast } from "sonner";
import type { MafiaRole } from "@/types/api";
import { normalizeSearchText } from "@/lib/search";

export const Route = createFileRoute("/game/$id/manage")({ component: ManageGamePage });

type RowState = {
  slot: number;
  nickname: string;
  profile_id?: string;
  guest_id?: string; // populated when loading an existing guest slot
  role?: MafiaRole;
  best_move?: string;
  compensation: string;
  yellow_cards: string;
  removed: string;
  extra_points: string;
  total_points: string;
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
  const judgesQuery = useQuery({
    queryKey: ["series", game.data?.series_id, "judges"],
    queryFn: () => seriesApi.judges(game.data!.series_id),
    enabled: !!game.data?.series_id && !!me,
  });

  const isJudge = !!me && (judgesQuery.data?.items ?? []).some((j) => j.profile_id === me.id);
  const canManage = series.data ? (canManageClub(me, series.data.club_id) || isJudge) : null;
  const queriesReady = status === "ready" && !!series.data && (!me || !judgesQuery.isLoading);
  useEffect(() => {
    if (queriesReady && canManage === false) navigate({ to: "/game/$id", params: { id } });
  }, [queriesReady, canManage, navigate, id]);

  // "unset" = not chosen yet (blocks publish in tournament), "none" = Не указан, uuid = specific judge
  const [judgeSelection, setJudgeSelection] = useState<"unset" | "none" | string>("unset");

  const [rows, setRows] = useState<RowState[]>([]);
  const [activeSlot, setActiveSlot] = useState<number | null>(null);
  const [savingDraft, setSavingDraft] = useState(false);
  const [publishing, setPublishing] = useState(false);
  const [backDialogOpen, setBackDialogOpen] = useState(false);
  const initialSnapshot = useRef<string>("");

  // Load rows: prefer server-side draft, fall back to published results
  useEffect(() => {
    if (!game.data) return;

    const draft = game.data.draft;
    if (draft?.rows?.length) {
      const participantMap = new Map((participants.data?.items ?? []).map((p) => [p.id, p]));
      const draftRows: RowState[] = Array.from({ length: 10 }, (_, i) => {
        const slot = i + 1;
        const dr = draft.rows.find((r) => r.slot === slot);
        const profile = dr?.profile_id ? participantMap.get(dr.profile_id) : undefined;
        return {
          slot,
          nickname: dr?.guest_nickname ?? (profile ? (profile.nickname || profile.name || "") : ""),
          profile_id: dr?.profile_id ?? undefined,
          guest_id: undefined,
          role: (dr?.role as MafiaRole | undefined) ?? "civilian",
          best_move: dr?.best_move ?? "",
          compensation: String(dr?.compensation ?? 0),
          yellow_cards: String(dr?.yellow_cards ?? 0),
          removed: String(dr?.removed ?? 0),
          extra_points: String(dr?.extra_points ?? 0),
          total_points: String(dr?.total_points ?? 0),
        };
      });
      setRows(draftRows);
      const loadedJudge = draft.judge_id ?? (draft.judge_confirmed ? "none" : "unset");
      setJudgeSelection(loadedJudge);
      initialSnapshot.current = JSON.stringify({ rows: draftRows, judge: loadedJudge });
      return;
    }

    // Fall back to published results
    const bySlot = new Map((game.data.results ?? []).filter((r) => r.place != null).map((r) => [r.place!, r]));
    const byProfile = new Map((game.data.results ?? []).filter((r) => r.profile_id).map((r) => [r.profile_id!, r]));
    const pids = game.data.participant_ids ?? [];

    const init: RowState[] = Array.from({ length: 10 }, (_, i) => {
      const slot = i + 1;
      const rr = bySlot.get(slot) ?? (pids[i] ? byProfile.get(pids[i]) : undefined);
      const legacyPid = !bySlot.size ? pids[i] : undefined;

      if (rr?.guest_id) {
        return {
          slot,
          nickname: rr.guest_nickname ?? "",
          profile_id: undefined,
          guest_id: rr.guest_id,
          role: rr.role ?? "civilian",
          best_move: rr.best_move ?? "",
          compensation: String(rr.compensation ?? 0),
          yellow_cards: String(rr.yellow_cards ?? 0),
          removed: String(rr.removed ?? 0),
          extra_points: String(rr.extra_points ?? 0),
          total_points: String(rr.total_points ?? 0),
        };
      }

      const pid = rr?.profile_id ?? legacyPid;
      const profile = pid ? (participants.data?.items ?? []).find((p) => p.id === pid) : undefined;
      return {
        slot,
        nickname: pid ? (profile?.nickname || profile?.name || "") : "",
        profile_id: pid,
        guest_id: undefined,
        role: rr?.role ?? "civilian",
        best_move: rr?.best_move ?? "",
        compensation: rr ? String(rr.compensation ?? 0) : "0",
        yellow_cards: rr ? String(rr.yellow_cards ?? 0) : "0",
        removed: rr ? String(rr.removed ?? 0) : "0",
        extra_points: rr ? String(rr.extra_points ?? 0) : "0",
        total_points: rr ? String(rr.total_points ?? 0) : "0",
      };
    });
    setRows(init);
    const loadedJudge = game.data.game_judge_id ?? (game.data.game_judge_confirmed ? "none" : "unset");
    setJudgeSelection(loadedJudge);
    initialSnapshot.current = JSON.stringify({ rows: init, judge: loadedJudge });
  }, [game.data, participants.data]); // eslint-disable-line react-hooks/exhaustive-deps

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
    // send guest_nickname when there's a typed name but no matched profile
    guest_nickname: !r.profile_id && r.nickname.trim() ? r.nickname.trim() : undefined,
    role: r.role,
    best_move: normalizeBestMove(r.best_move),
    compensation: toNum(r.compensation),
    yellow_cards: toNum(r.yellow_cards),
    removed: toNum(r.removed),
    extra_points: toNum(r.extra_points),
    total_points: toNum(r.total_points),
  }));

  // Slots with a nickname but no profile_id are guest slots — allowed for both draft and publish.
  const emptySlots = rows
    .filter((r) => !r.profile_id && !r.nickname.trim())
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
  const duplicateSlotSet = new Set(duplicateParticipantSlots);

  const isDirty = initialSnapshot.current !== "" &&
    JSON.stringify({ rows, judge: judgeSelection }) !== initialSnapshot.current;

  const goBack = () => {
    qc.invalidateQueries({ queryKey: ["game", id, "full"] });
    void navigate({ to: "/game/$id", params: { id } });
  };

  const handleBack = () => {
    if (isDirty) {
      setBackDialogOpen(true);
    } else {
      goBack();
    }
  };

  const saveDraft = async () => {
    if (duplicateParticipantSlots.length > 0) {
      toast.error(`Нельзя сохранить: один и тот же игрок выбран в нескольких слотах (${duplicateParticipantSlots.join(", ")}).`);
      return;
    }
    setSavingDraft(true);
    try {
      const judgeId = judgeSelection === "unset" || judgeSelection === "none" ? null : judgeSelection;
      await gamesApi.saveDraft(id, payloadRows, judgeId, judgeSelection !== "unset");
      initialSnapshot.current = JSON.stringify({ rows, judge: judgeSelection });
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
    if (emptySlots.length > 0) {
      toast.error(`Нельзя опубликовать: слоты ${emptySlots.join(", ")} пустые. Заполните никнейм или выберите игрока.`);
      return;
    }
    setPublishing(true);
    try {
      const judgeId = judgeSelection === "unset" || judgeSelection === "none" ? null : judgeSelection;
      await gamesApi.publish(id, payloadRows, judgeId, judgeSelection !== "unset");
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
      <AlertDialog open={backDialogOpen} onOpenChange={setBackDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Несохранённые изменения</AlertDialogTitle>
            <AlertDialogDescription>
              Вы изменили данные игры. Сохранить черновик перед выходом?
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Отмена</AlertDialogCancel>
            <Button variant="outline" onClick={() => { setBackDialogOpen(false); goBack(); }}>
              Не сохранять
            </Button>
            <AlertDialogAction onClick={async () => {
              setBackDialogOpen(false);
              await saveDraft();
              goBack();
            }}>
              Сохранить черновик
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <PageHeader
        eyebrow="Управление игрой"
        title={game.data.name || `Игра #${game.data.number}`}
        actions={<Button variant="outline" onClick={handleBack}>Назад к игре</Button>}
      />

      {game.data.is_tournament && (
        <section className="rounded-2xl border border-border/60 bg-card/60 p-5">
          <div className="flex flex-wrap items-end gap-4">
            <div className="space-y-1.5">
              <Label className="text-sm font-medium">Главный судья игры</Label>
              <Select value={judgeSelection} onValueChange={setJudgeSelection}>
                <SelectTrigger className="w-64">
                  <SelectValue placeholder="Выбрать судью…" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">Не указан</SelectItem>
                  {(judgesQuery.data?.items ?? [])
                    .filter((j) => j.role === 1)
                    .map((j) => (
                      <SelectItem key={j.profile_id} value={j.profile_id}>
                        {j.nickname || j.name || j.profile_id}
                      </SelectItem>
                    ))}
                </SelectContent>
              </Select>
            </div>
            {judgeSelection === "unset" && (
              <p className="pb-1 text-xs text-amber-600">Выберите судью или «Не указан» перед публикацией.</p>
            )}
          </div>
        </section>
      )}

      <section className="rounded-2xl border border-border/60 bg-card/60 p-6">
        <div className="overflow-x-auto">
          <table className="w-full min-w-[1100px] text-sm">
            <thead>
              <tr className="border-b border-border/60 text-left text-xs uppercase text-muted-foreground">
                <th className="px-2 py-2">Слот</th>
                <th className="px-2 py-2">Никнейм</th>
                <th className="px-2 py-2">Роль</th>
                <th className="w-24 px-2 py-2">Лучший ход</th>
                <th className="px-2 py-2">Компенсация</th>
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
                        style={duplicateSlotSet.has(r.slot) ? {
                          borderColor: "oklch(0.62 0.22 25 / 0.9)",
                          boxShadow: "0 0 0 1px oklch(0.62 0.22 25 / 0.5), 0 0 12px -2px oklch(0.62 0.22 25 / 0.7)",
                        } : undefined}
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
                  <td className="px-2 py-2"><Input value={r.compensation} onChange={(e) => setCell(r.slot, { compensation: e.target.value })} className="h-8" /></td>
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
          <Button
            disabled={
              savingDraft || publishing ||
              emptySlots.length > 0 ||
              duplicateParticipantSlots.length > 0 ||
              (!!game.data.is_tournament && judgeSelection === "unset")
            }
            onClick={() => void publish()}
          >
            {publishing ? "Публикация..." : "Опубликовать"}
          </Button>
        </div>

        <p className="mt-3 text-xs text-muted-foreground">
          Черновики видны только судьям и лидерам клуба. После публикации игра становится доступной всем.
          Незарегистрированных гостей можно добавить просто введя никнейм — без привязки к аккаунту.
        </p>
        {emptySlots.length > 0 && (
          <p className="mt-1 text-xs text-amber-600">
            Публикация недоступна: слоты {emptySlots.join(", ")} пустые.
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

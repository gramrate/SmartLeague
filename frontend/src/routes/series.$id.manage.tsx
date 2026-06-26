import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { PageHeader, PageShell } from "@/components/site/PageShell";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { clubsApi, ApiError, seriesApi } from "@/lib/api";
import { ErrorBlock, LoadingBlock } from "@/components/site/States";
import { useAuthStore } from "@/lib/auth-store";
import { canManageClub, displayUserName } from "@/lib/roles";
import { JudgeRole } from "@/types/api";
import { useEffect, useMemo, useState } from "react";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { fromInputDate, toInputDate } from "@/lib/format";
import { toast } from "sonner";
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from "@/components/ui/alert-dialog";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { X } from "lucide-react";

export const Route = createFileRoute("/series/$id/manage")({ component: SeriesManagePage });

function SeriesManagePage() {
  const { id } = Route.useParams();
  const me = useAuthStore((s) => s.me);
  const status = useAuthStore((s) => s.status);
  const navigate = useNavigate();
  const qc = useQueryClient();

  const seriesQ = useQuery({ queryKey: ["series", id], queryFn: () => seriesApi.get(id) });
  const series = seriesQ.data;
  const canManage = !!series && canManageClub(me, series.club_id);

  useEffect(() => {
    if (status === "ready" && series && !canManage) navigate({ to: "/series/$id", params: { id } });
  }, [status, series, canManage, navigate, id]);

  const club = useQuery({
    queryKey: ["club", series?.club_id],
    queryFn: () => clubsApi.get(series!.club_id),
    enabled: !!series?.club_id,
  });
  const participants = useQuery({
    queryKey: ["series", id, "participants", "manage-payments"],
    queryFn: () => seriesApi.participants(id, 200, 0),
    enabled: !!series && Number(series.price_rub ?? 0) > 0,
    staleTime: 5 * 60 * 1000,
    refetchOnWindowFocus: false,
  });
  const allParticipants = useQuery({
    queryKey: ["series", id, "participants"],
    queryFn: () => seriesApi.participants(id, 200, 0),
    enabled: !!series,
  });
  const judgesQ = useQuery({
    queryKey: ["series", id, "judges"],
    queryFn: () => seriesApi.judges(id),
    enabled: !!series,
  });
  const payments = useQuery({
    queryKey: ["series", id, "payments"],
    queryFn: () => seriesApi.payments(id),
    enabled: !!series && Number(series.price_rub ?? 0) > 0,
    staleTime: 5 * 60 * 1000,
    refetchOnWindowFocus: false,
  });

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [startAt, setStartAt] = useState("");
  const [endAt, setEndAt] = useState("");
  const [priceRub, setPriceRub] = useState("0");
  const [isRating, setIsRating] = useState(false);
  const [isClubOnly, setIsClubOnly] = useState(false);
  const [showToAll, setShowToAll] = useState(true);
  const [isClosed, setIsClosed] = useState(false);
  const [isTournament, setIsTournament] = useState(false);
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [warnDialogOpen, setWarnDialogOpen] = useState(false);
  const [pendingWarnings, setPendingWarnings] = useState<string[]>([]);
  const [paidOverrides, setPaidOverrides] = useState<Record<string, boolean>>({});
  const [settingJudge, setSettingJudge] = useState(false);
  const [removingJudge, setRemovingJudge] = useState<string | null>(null);
  const [judgeDialogOpen, setJudgeDialogOpen] = useState(false);
  const [judgeSearchQ, setJudgeSearchQ] = useState("");
  const nameLimit = 100;
  const descriptionLimit = 1000;

  useEffect(() => {
    if (!series) return;
    setName(series.name ?? "");
    setDescription(series.description ?? "");
    setStartAt(toInputDate(series.start_at));
    setEndAt(toInputDate(series.end_at));
    setPriceRub(String(Number(series.price_rub ?? 0)));
    setIsRating(!!series.is_rating);
    setIsClubOnly(!!series.is_club_only);
    setShowToAll(series.show_to_all !== false);
    setIsClosed(!!series.is_closed);
    setIsTournament(!!series.is_tournament);
  }, [series]);

  const dirty = useMemo(() => {
    if (!series) return false;
    return (
      name.trim() !== (series.name ?? "") ||
      description.trim() !== (series.description ?? "") ||
      startAt !== toInputDate(series.start_at) ||
      endAt !== toInputDate(series.end_at) ||
      Math.max(0, Number(priceRub || 0)) !== Number(series.price_rub ?? 0) ||
      isRating !== !!series.is_rating ||
      isClubOnly !== !!series.is_club_only ||
      showToAll !== (series.show_to_all !== false) ||
      isClosed !== !!series.is_closed ||
      isTournament !== !!series.is_tournament
    );
  }, [series, name, description, startAt, endAt, priceRub, isRating, isClubOnly, showToAll, isClosed, isTournament]);

  useEffect(() => {
    setPaidOverrides({});
  }, [payments.data?.paid_profile_ids]);

  if (seriesQ.isLoading) return <PageShell><LoadingBlock /></PageShell>;
  if (seriesQ.error) return <PageShell><ErrorBlock error={seriesQ.error} /></PageShell>;
  if (!series || !canManage) return null;
  const paidBase = new Set(payments.data?.paid_profile_ids ?? []);
  const isPaid = (profileId: string) => paidOverrides[profileId] ?? paidBase.has(profileId);
  const participantsList = participants.data?.items ?? [];
  const paidParticipants = participantsList.filter((p) => isPaid(p.id));
  const unpaidParticipants = participantsList.filter((p) => !isPaid(p.id));

  const handleSave = () => {
    const warnings: string[] = [];
    if (series.is_tournament && !isTournament) {
      warnings.push("Все судьи серии будут удалены. В каждой игре будет очищена информация о судье.");
    }
    if (Number(series.price_rub ?? 0) > 0 && Math.max(0, Number(priceRub || 0)) === 0) {
      warnings.push("Все данные об оплате участников будут сброшены.");
    }
    if (warnings.length > 0) {
      setPendingWarnings(warnings);
      setWarnDialogOpen(true);
    } else {
      void save();
    }
  };

  const save = async () => {
    setSaving(true);
    try {
      await seriesApi.update(id, {
        name: name.trim(),
        description: description.trim(),
        start_at: fromInputDate(startAt),
        end_at: fromInputDate(endAt),
        price_rub: Math.max(0, Number(priceRub || 0)),
        is_rating: isRating,
        is_club_only: isClubOnly,
        show_to_all: isClubOnly ? showToAll : true,
        is_closed: isClosed,
        is_tournament: isTournament,
      });
      qc.invalidateQueries({ queryKey: ["series", id] });
      qc.invalidateQueries({ queryKey: ["series", id, "full"] });
      qc.invalidateQueries({ queryKey: ["series"] });
      qc.invalidateQueries({ queryKey: ["club", series.club_id, "series"] });
      toast.success("Серия обновлена");
      navigate({ to: "/series/$id", params: { id } });
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : "Ошибка");
    } finally {
      setSaving(false);
    }
  };

  const removeSeries = async () => {
    setDeleting(true);
    try {
      await seriesApi.delete(id);
      qc.invalidateQueries({ queryKey: ["series"] });
      qc.invalidateQueries({ queryKey: ["club", series.club_id, "series"] });
      toast.success("Серия удалена");
      navigate({ to: "/clubs/$id", params: { id: series.club_id } });
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : "Ошибка");
    } finally {
      setDeleting(false);
    }
  };

  const setPaid = async (profileId: string, paid: boolean) => {
    try {
      await seriesApi.setPayment(id, profileId, paid);
      setPaidOverrides((prev) => ({ ...prev, [profileId]: paid }));
      toast.success(paid ? "Оплата отмечена" : "Оплата снята");
    } catch (e) {
      toast.error(e instanceof ApiError ? e.message : "Ошибка");
    }
  };

  const invalidateJudgeRelated = () => {
    qc.invalidateQueries({ queryKey: ["series", id, "judges"] });
    qc.invalidateQueries({ queryKey: ["series", id, "participants"], exact: true });
    qc.invalidateQueries({ queryKey: ["series", id, "participants", "manage-payments"], exact: true });
    qc.invalidateQueries({ queryKey: ["series", id, "payments"], exact: true });
  };

  const judgeSet = (profileId: string, role: JudgeRole) => async () => {
    setSettingJudge(true);
    try {
      await seriesApi.setJudge(id, profileId, role);
      invalidateJudgeRelated();
      toast.success("Судья назначен");
    } catch (e) { toast.error(e instanceof ApiError ? e.message : "Ошибка"); }
    finally { setSettingJudge(false); }
  };

  const judgeRemove = (profileId: string) => async () => {
    setRemovingJudge(profileId);
    try {
      await seriesApi.removeJudge(id, profileId);
      invalidateJudgeRelated();
      toast.success("Судья снят");
    } catch (e) { toast.error(e instanceof ApiError ? e.message : "Ошибка"); }
    finally { setRemovingJudge(null); }
  };

  return (
    <PageShell>
      <AlertDialog open={warnDialogOpen} onOpenChange={setWarnDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Внимание: данные будут сброшены</AlertDialogTitle>
            <AlertDialogDescription asChild>
              <div className="space-y-2">
                {pendingWarnings.map((w, i) => (
                  <p key={i} className="text-sm">— {w}</p>
                ))}
                <p className="mt-2 text-sm font-medium">Это действие нельзя отменить. Продолжить?</p>
              </div>
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Отмена</AlertDialogCancel>
            <AlertDialogAction variant="destructive" onClick={() => { setWarnDialogOpen(false); void save(); }}>
              Сохранить и сбросить
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <PageHeader
        eyebrow={club.data?.name ?? "Серия"}
        title={`Управление: ${series.name}`}
        actions={
          <Button variant="outline" onClick={() => {
            qc.invalidateQueries({ queryKey: ["series", id] });
            void navigate({ to: "/series/$id", params: { id } });
          }}>
            К серии
          </Button>
        }
      />

      <div className="mx-auto max-w-xl">
        <section className="rounded-2xl border border-border/60 bg-card/60 p-6">
          <h2 className="mb-4 font-display text-lg font-semibold">Параметры серии</h2>
          <div className="space-y-4">
            <div className="space-y-1.5">
              <Label>Название</Label>
              <Input value={name} onChange={(e) => setName(e.target.value)} maxLength={nameLimit} />
              <p className="text-xs text-muted-foreground">{name.length}/{nameLimit}</p>
            </div>
            <div className="space-y-1.5">
              <Label>Описание</Label>
              <Textarea rows={4} value={description} onChange={(e) => setDescription(e.target.value)} maxLength={descriptionLimit} />
              <p className="text-xs text-muted-foreground">{description.length}/{descriptionLimit}</p>
            </div>
            <div className="grid gap-3 sm:grid-cols-2">
              <div className="space-y-1.5"><Label>Начало</Label><Input type="date" value={startAt} onChange={(e) => setStartAt(e.target.value)} /></div>
              <div className="space-y-1.5"><Label>Конец</Label><Input type="date" value={endAt} onChange={(e) => setEndAt(e.target.value)} /></div>
            </div>
            <div className="space-y-1.5">
              <Label>Стоимость (₽)</Label>
              <Input type="number" min={0} step={1} value={priceRub} onChange={(e) => setPriceRub(e.target.value)} />
            </div>
            <div className="space-y-1.5">
              <Label>Статус регистрации</Label>
              <Select value={isClosed ? "closed" : "open"} onValueChange={(v) => setIsClosed(v === "closed")}>
                <SelectTrigger><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="open">Открыта</SelectItem>
                  <SelectItem value="closed">Закрыта</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <label className="flex items-center gap-2 text-sm">
              <Checkbox checked={isRating} onCheckedChange={(v) => setIsRating(!!v)} />
              На рейтинг
            </label>
            <label className="flex items-center gap-2 text-sm">
              <Checkbox checked={isTournament} onCheckedChange={(v) => { setIsTournament(!!v); if (!!v) setIsRating(true); }} />
              Турнир
            </label>
            <label className="flex items-center gap-2 text-sm">
              <Checkbox checked={isClubOnly} onCheckedChange={(v) => setIsClubOnly(!!v)} />
              Только для участников клуба
            </label>
            {isClubOnly && (
              <div className="rounded-lg border border-border/60 bg-muted/30 p-3 space-y-2">
                <p className="text-sm font-medium">Показывать серию всем?</p>
                <label className="flex items-center gap-2 text-sm text-foreground">
                  <Checkbox checked={showToAll} onCheckedChange={(v) => setShowToAll(!!v)} />
                  Да — серия видна всем, но вступить могут только участники клуба
                </label>
                {!showToAll && (
                  <p className="text-xs text-muted-foreground">Серия и её игры будут скрыты от людей не из клуба</p>
                )}
              </div>
            )}
            <div className="flex flex-wrap items-center gap-2 pt-2">
              <Button onClick={handleSave} disabled={saving || !dirty || !name.trim() || !startAt || !endAt || name.length > nameLimit || description.length > descriptionLimit}>
                {saving ? "Сохранение..." : "Сохранить"}
              </Button>
            </div>
          </div>
        </section>

        {Number(series.price_rub ?? 0) > 0 && (
          <section className="mt-8 rounded-2xl border border-border/60 bg-card/60 p-6">
            <h2 className="mb-2 font-display text-lg font-semibold">Оплаты участников</h2>
            <p className="mb-4 text-sm text-muted-foreground">Отмечайте тех, кто оплатил участие в серии.</p>
            <div className="mb-4 flex flex-wrap items-center gap-4 text-sm">
              <span>Оплатили: {paidParticipants.length}</span>
              <span>Не оплатили: {unpaidParticipants.length}</span>
            </div>
            {!participantsList.length ? (
              <p className="text-sm text-muted-foreground">Участников пока нет.</p>
            ) : (
              <Tabs defaultValue="paid" className="w-full">
                <TabsList className="grid w-full grid-cols-2">
                  <TabsTrigger value="paid">Оплатили</TabsTrigger>
                  <TabsTrigger value="unpaid">Не оплатили</TabsTrigger>
                </TabsList>
                <TabsContent value="paid" className="mt-3">
                  {!paidParticipants.length ? (
                    <p className="text-sm text-muted-foreground">Пока никто не оплатил.</p>
                  ) : (
                    <ul className="space-y-2">
                      {paidParticipants.map((p) => (
                        <li key={p.id} className="flex items-center justify-between gap-3 rounded-lg border border-border/50 px-3 py-2">
                          <span className="truncate text-sm">{p.nickname || p.name || p.email || p.id}</span>
                          <label className="flex items-center gap-2 text-sm">
                            <Checkbox checked onCheckedChange={(v) => void setPaid(p.id, !!v)} />
                            Оплачено
                          </label>
                        </li>
                      ))}
                    </ul>
                  )}
                </TabsContent>
                <TabsContent value="unpaid" className="mt-3">
                  {!unpaidParticipants.length ? (
                    <p className="text-sm text-muted-foreground">Все участники оплатили.</p>
                  ) : (
                    <ul className="space-y-2">
                      {unpaidParticipants.map((p) => (
                        <li key={p.id} className="flex items-center justify-between gap-3 rounded-lg border border-border/50 px-3 py-2">
                          <span className="truncate text-sm">{p.nickname || p.name || p.email || p.id}</span>
                          <label className="flex items-center gap-2 text-sm">
                            <Checkbox checked={false} onCheckedChange={(v) => void setPaid(p.id, !!v)} />
                            Оплачено
                          </label>
                        </li>
                      ))}
                    </ul>
                  )}
                </TabsContent>
              </Tabs>
            )}
          </section>
        )}

        <section className="mt-8 rounded-2xl border border-border/60 bg-card/60 p-6">
          <div className="mb-3 flex items-center justify-between">
            <h2 className="font-display text-lg font-semibold">Судьи</h2>
            <Button size="sm" variant="outline" onClick={() => setJudgeDialogOpen(true)}>
              Изменить судейский состав
            </Button>
          </div>
          <p className="mb-4 text-sm text-muted-foreground">
            Судьи могут редактировать игры серии, но не могут создавать и удалять их.
          </p>
          {!judgesQ.data?.items?.length ? (
            <p className="text-sm text-muted-foreground">Судей пока нет.</p>
          ) : (
            <ul className="divide-y divide-border/40">
              {judgesQ.data.items.map((j) => (
                <li key={j.profile_id} className="flex items-center justify-between gap-2 py-2 text-sm">
                  <div className="flex min-w-0 items-center gap-2">
                    <Link to="/user/$id" params={{ id: j.profile_id }} className="truncate hover:text-primary">
                      {displayUserName(j)}
                    </Link>
                    <span className={`shrink-0 rounded-full px-2 py-0.5 text-xs font-medium ${j.role === JudgeRole.Main ? "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/40 dark:text-yellow-300" : "bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300"}`}>
                      {j.role === JudgeRole.Main ? "Главный" : "Боковой"}
                    </span>
                  </div>
                  <Button
                    size="sm"
                    variant="ghost"
                    className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
                    disabled={removingJudge === j.profile_id}
                    onClick={judgeRemove(j.profile_id)}
                  >
                    <X className="h-4 w-4" />
                  </Button>
                </li>
              ))}
            </ul>
          )}
        </section>

        <section className="mt-8 rounded-2xl border border-destructive/40 bg-card/60 p-6">
          <h2 className="mb-2 font-display text-lg font-semibold text-destructive">Опасные действия</h2>
          <p className="mb-4 text-sm text-muted-foreground">Удаление серии необратимо.</p>
          <AlertDialog>
            <AlertDialogTrigger asChild>
              <Button variant="destructive">Удалить серию</Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Удалить серию?</AlertDialogTitle>
                <AlertDialogDescription>
                  Это действие нельзя отменить. Будут удалены серия и связанные с ней данные.
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>Отмена</AlertDialogCancel>
                <AlertDialogAction
                  variant="destructive"
                  onClick={() => void removeSeries()}
                  disabled={deleting}
                >
                  {deleting ? "Удаление..." : "Удалить"}
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </section>
      </div>

      <Dialog open={judgeDialogOpen} onOpenChange={(open) => { setJudgeDialogOpen(open); if (!open) setJudgeSearchQ(""); }}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Управление судьями</DialogTitle>
            <DialogDescription>Назначайте судей из участников серии.</DialogDescription>
          </DialogHeader>
          <Tabs defaultValue="judges" className="w-full">
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="judges">Текущие судьи</TabsTrigger>
              <TabsTrigger value="participants">Участники</TabsTrigger>
            </TabsList>
            <TabsContent value="judges" className="mt-3">
              <div className="rounded-lg border border-border/40 bg-background/40 p-3">
                {!judgesQ.data?.items?.length ? (
                  <p className="text-sm text-muted-foreground">Судей пока нет.</p>
                ) : (
                  <ul className="divide-y divide-border/40">
                    {judgesQ.data.items.map((j) => (
                      <li key={j.profile_id} className="flex items-center justify-between gap-2 py-2 text-sm">
                        <div className="flex min-w-0 items-center gap-2">
                          <Link to="/user/$id" params={{ id: j.profile_id }} className="truncate hover:text-primary">
                            {displayUserName(j)}
                          </Link>
                          <span className={`shrink-0 rounded-full px-2 py-0.5 text-xs font-medium ${j.role === JudgeRole.Main ? "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/40 dark:text-yellow-300" : "bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300"}`}>
                            {j.role === JudgeRole.Main ? "Главный" : "Боковой"}
                          </span>
                        </div>
                        <div className="flex shrink-0 items-center gap-1">
                          {j.role !== JudgeRole.Main && (
                            <Button size="sm" variant="outline" disabled={settingJudge} onClick={judgeSet(j.profile_id, JudgeRole.Main)}>
                              → Главный
                            </Button>
                          )}
                          {j.role !== JudgeRole.Side && (
                            <Button size="sm" variant="outline" disabled={settingJudge} onClick={judgeSet(j.profile_id, JudgeRole.Side)}>
                              → Боковой
                            </Button>
                          )}
                          <Button
                            size="sm"
                            variant="ghost"
                            className="h-7 w-7 p-0 text-muted-foreground hover:text-destructive"
                            disabled={removingJudge === j.profile_id}
                            onClick={judgeRemove(j.profile_id)}
                          >
                            <X className="h-4 w-4" />
                          </Button>
                        </div>
                      </li>
                    ))}
                  </ul>
                )}
              </div>
            </TabsContent>
            <TabsContent value="participants" className="mt-3 space-y-3">
              <div className="space-y-1.5">
                <Label>Поиск участника</Label>
                <Input
                  value={judgeSearchQ}
                  maxLength={50}
                  onChange={(e) => setJudgeSearchQ(e.target.value)}
                  placeholder="Никнейм…"
                />
              </div>
              <div className="rounded-lg border border-border/40 bg-background/40 p-3">
                {(() => {
                  const nonJudges = (allParticipants.data?.items ?? [])
                    .filter((p) => !(judgesQ.data?.items ?? []).some((j) => j.profile_id === p.id))
                    .filter((p) => !judgeSearchQ.trim() || displayUserName(p).toLowerCase().includes(judgeSearchQ.trim().toLowerCase()));
                  if (!nonJudges.length) {
                    return (
                      <p className="text-sm text-muted-foreground">
                        {judgeSearchQ.trim() ? "Никого не найдено." : "Все участники уже являются судьями."}
                      </p>
                    );
                  }
                  return (
                    <ul className="divide-y divide-border/40">
                      {nonJudges.map((p) => (
                        <li key={p.id} className="flex items-center justify-between gap-2 py-2 text-sm">
                          <Link to="/user/$id" params={{ id: p.id }} className="min-w-0 flex-1 truncate hover:text-primary">
                            {displayUserName(p)}
                          </Link>
                          <div className="flex shrink-0 gap-1">
                            <Button size="sm" variant="outline" disabled={settingJudge} onClick={judgeSet(p.id, JudgeRole.Main)}>
                              Главный
                            </Button>
                            <Button size="sm" variant="outline" disabled={settingJudge} onClick={judgeSet(p.id, JudgeRole.Side)}>
                              Боковой
                            </Button>
                          </div>
                        </li>
                      ))}
                    </ul>
                  );
                })()}
              </div>
            </TabsContent>
          </Tabs>
        </DialogContent>
      </Dialog>
    </PageShell>
  );
}

import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { PageHeader, PageShell } from "@/components/site/PageShell";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { clubsApi, ApiError, seriesApi } from "@/lib/api";
import { ErrorBlock, LoadingBlock } from "@/components/site/States";
import { useAuthStore } from "@/lib/auth-store";
import { canManageClub } from "@/lib/roles";
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
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

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
  });
  const payments = useQuery({
    queryKey: ["series", id, "payments"],
    queryFn: () => seriesApi.payments(id),
    enabled: !!series && Number(series.price_rub ?? 0) > 0,
  });

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [startAt, setStartAt] = useState("");
  const [endAt, setEndAt] = useState("");
  const [priceRub, setPriceRub] = useState("0");
  const [isRating, setIsRating] = useState(false);
  const [isClubOnly, setIsClubOnly] = useState(false);
  const [isClosed, setIsClosed] = useState(false);
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [paidOverrides, setPaidOverrides] = useState<Record<string, boolean>>({});

  useEffect(() => {
    if (!series) return;
    setName(series.name ?? "");
    setDescription(series.description ?? "");
    setStartAt(toInputDate(series.start_at));
    setEndAt(toInputDate(series.end_at));
    setPriceRub(String(Number(series.price_rub ?? 0)));
    setIsRating(!!series.is_rating);
    setIsClubOnly(!!series.is_club_only);
    setIsClosed(!!series.is_closed);
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
      isClosed !== !!series.is_closed
    );
  }, [series, name, description, startAt, endAt, priceRub, isRating, isClubOnly, isClosed]);

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
        is_closed: isClosed,
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

  return (
    <PageShell>
      <PageHeader
        eyebrow={club.data?.name ?? "Серия"}
        title={`Управление: ${series.name}`}
        actions={<Button variant="outline" asChild><Link to="/series/$id" params={{ id }}>К серии</Link></Button>}
      />

      <div className="mx-auto max-w-xl">
        <section className="rounded-2xl border border-border/60 bg-card/60 p-6">
          <h2 className="mb-4 font-display text-lg font-semibold">Параметры серии</h2>
          <div className="space-y-4">
            <div className="space-y-1.5"><Label>Название</Label><Input value={name} onChange={(e) => setName(e.target.value)} maxLength={200} /></div>
            <div className="space-y-1.5"><Label>Описание</Label><Textarea rows={4} value={description} onChange={(e) => setDescription(e.target.value)} maxLength={10000} /></div>
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
              <Checkbox checked={isClubOnly} onCheckedChange={(v) => setIsClubOnly(!!v)} />
              Только для участников клуба
            </label>
            <div className="flex flex-wrap items-center gap-2 pt-2">
              <Button onClick={() => void save()} disabled={saving || !dirty || !name.trim() || !description.trim() || !startAt || !endAt}>
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
                  className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
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
    </PageShell>
  );
}

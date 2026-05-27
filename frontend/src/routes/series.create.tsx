import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useEffect, useState } from "react";
import { useAuthStore } from "@/lib/auth-store";
import { isClubManager } from "@/lib/roles";
import { seriesApi, ApiError } from "@/lib/api";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { fromInputDate } from "@/lib/format";
import { toast } from "sonner";

export const Route = createFileRoute("/series/create")({ component: CreateSeriesPage });

function CreateSeriesPage() {
  const me = useAuthStore((s) => s.me);
  const status = useAuthStore((s) => s.status);
  const navigate = useNavigate();
  const canCreate = !!me?.club_id && isClubManager(me.club_state);

  useEffect(() => {
    if (status === "ready" && !canCreate) navigate({ to: "/series" });
  }, [canCreate, status, navigate]);

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [startAt, setStartAt] = useState("");
  const [endAt, setEndAt] = useState("");
  const [priceRub, setPriceRub] = useState("0");
  const [isRating, setIsRating] = useState(false);
  const [isClubOnly, setIsClubOnly] = useState(false);
  const [busy, setBusy] = useState(false);
  const nameLimit = 200;
  const descriptionLimit = 10000;

  if (!canCreate) return null;

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setBusy(true);
    try {
      const s = await seriesApi.create(me!.club_id!, {
        name: name.trim(),
        description: description.trim(),
        start_at: fromInputDate(startAt),
        end_at: fromInputDate(endAt),
        price_rub: Math.max(0, Number(priceRub || 0)),
        is_rating: isRating,
        is_club_only: isClubOnly,
      });
      toast.success("Серия создана");
      navigate({ to: "/series/$id", params: { id: s.id } });
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : "Ошибка");
    } finally { setBusy(false); }
  };

  return (
    <PageShell>
      <div className="mx-auto max-w-xl">
        <PageHeader eyebrow="Серии" title="Создание серии" />
        <form onSubmit={onSubmit} className="space-y-4 rounded-2xl border border-border/60 bg-card/60 p-6">
          <div className="space-y-1.5">
            <Label>Название</Label>
            <Input value={name} onChange={(e) => setName(e.target.value)} required maxLength={nameLimit} />
            <p className="text-xs text-muted-foreground">{name.length}/{nameLimit} · осталось {nameLimit - name.length}</p>
          </div>
          <div className="space-y-1.5">
            <Label>Описание</Label>
            <Textarea rows={4} value={description} onChange={(e) => setDescription(e.target.value)} required maxLength={descriptionLimit} />
            <p className="text-xs text-muted-foreground">{description.length}/{descriptionLimit} · осталось {descriptionLimit - description.length}</p>
          </div>
          <div className="grid gap-3 sm:grid-cols-2">
            <div className="space-y-1.5"><Label>Начало</Label><Input type="date" value={startAt} onChange={(e) => setStartAt(e.target.value)} required /></div>
            <div className="space-y-1.5"><Label>Конец</Label><Input type="date" value={endAt} onChange={(e) => setEndAt(e.target.value)} required /></div>
          </div>
          <div className="space-y-1.5">
            <Label>Стоимость (₽)</Label>
            <Input type="number" min={0} step={1} value={priceRub} onChange={(e) => setPriceRub(e.target.value)} />
          </div>
          <label className="flex items-center gap-2 text-sm text-foreground">
            <Checkbox checked={isRating} onCheckedChange={(v) => setIsRating(!!v)} />
            На рейтинг
          </label>
          <label className="flex items-center gap-2 text-sm text-foreground">
            <Checkbox checked={isClubOnly} onCheckedChange={(v) => setIsClubOnly(!!v)} />
            Только для участников клуба
          </label>
          <Button type="submit" disabled={busy || !name || !description || !startAt || !endAt || name.length > nameLimit || description.length > descriptionLimit}>
            {busy ? "Создание…" : "Создать серию"}
          </Button>
        </form>
      </div>
    </PageShell>
  );
}

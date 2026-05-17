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
  const [busy, setBusy] = useState(false);

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
          <div className="space-y-1.5"><Label>Название</Label><Input value={name} onChange={(e) => setName(e.target.value)} required maxLength={200} /></div>
          <div className="space-y-1.5"><Label>Описание</Label><Textarea rows={4} value={description} onChange={(e) => setDescription(e.target.value)} required maxLength={10000} /></div>
          <div className="grid gap-3 sm:grid-cols-2">
            <div className="space-y-1.5"><Label>Начало</Label><Input type="date" value={startAt} onChange={(e) => setStartAt(e.target.value)} required /></div>
            <div className="space-y-1.5"><Label>Конец</Label><Input type="date" value={endAt} onChange={(e) => setEndAt(e.target.value)} required /></div>
          </div>
          <Button type="submit" disabled={busy || !name || !description || !startAt || !endAt}>
            {busy ? "Создание…" : "Создать серию"}
          </Button>
        </form>
      </div>
    </PageShell>
  );
}

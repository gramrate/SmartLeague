import { createFileRoute, Link } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { clubsApi, seriesApi } from "@/lib/api";
import { LoadingBlock, EmptyBlock } from "@/components/site/States";
import { fmtDateRange, fmtRub } from "@/lib/format";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useDebouncedValue } from "@/lib/useDebouncedValue";

export const Route = createFileRoute("/clubs/$id/series")({ component: ClubSeriesPage });

function ClubSeriesPage() {
  const { id } = Route.useParams();
  const [q, setQ] = useState("");
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");
  const [showPast, setShowPast] = useState(false);
  const [showClosed, setShowClosed] = useState(false);
  const [ratingFilter, setRatingFilter] = useState<"all" | "rating" | "non_rating">("all");
  const qLimit = 100;
  const [page, setPage] = useState(1);
  const limit = 10;
  const offset = (page - 1) * limit;
  const debouncedQ = useDebouncedValue(q, 150);
  const club = useQuery({ queryKey: ["club", id], queryFn: () => clubsApi.get(id) });
  const series = useQuery({
    queryKey: ["club", id, "series", debouncedQ, from, to, ratingFilter, showPast, showClosed, page],
    queryFn: () => seriesApi.all({
      q: debouncedQ || undefined,
      club: id,
      from: from || undefined,
      to: to || undefined,
      is_rating: ratingFilter === "all" ? undefined : ratingFilter === "rating",
      show_past: showPast,
      show_closed: showClosed,
      limit,
      offset,
    }),
  });

  return (
    <PageShell>
      <PageHeader eyebrow={club.data?.name ?? "Клуб"} title="Все серии" />
      <div className="mb-6 grid gap-3 rounded-xl border border-border/60 bg-card/40 p-4 sm:grid-cols-2 lg:grid-cols-6">
        <div className="space-y-1"><Label className="text-xs">Название</Label><Input value={q} maxLength={qLimit} onChange={(e) => { setQ(e.target.value); setPage(1); }} placeholder="Поиск…" /><p className="text-xs text-muted-foreground">{q.length}/{qLimit}</p></div>
        <div className="space-y-1"><Label className="text-xs">С</Label><Input type="date" value={from} onChange={(e) => { setFrom(e.target.value); setPage(1); }} /></div>
        <div className="space-y-1"><Label className="text-xs">По</Label><Input type="date" value={to} onChange={(e) => { setTo(e.target.value); setPage(1); }} /></div>
        <div className="space-y-1">
          <Label className="text-xs">Тип игр</Label>
          <Select value={ratingFilter} onValueChange={(v: "all" | "rating" | "non_rating") => { setRatingFilter(v); setPage(1); }}>
            <SelectTrigger><SelectValue /></SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Все игры</SelectItem>
              <SelectItem value="rating">На рейтинг</SelectItem>
              <SelectItem value="non_rating">Без рейтинга</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <label className="flex items-end gap-2 pb-2 text-sm text-muted-foreground">
          <Checkbox checked={showPast} onCheckedChange={(v) => { setShowPast(!!v); setPage(1); }} />
          Показать прошедшие
        </label>
        <label className="flex items-end gap-2 pb-2 text-sm text-muted-foreground">
          <Checkbox checked={showClosed} onCheckedChange={(v) => { setShowClosed(!!v); setPage(1); }} />
          Показать закрытые
        </label>
      </div>
      {series.isLoading ? <LoadingBlock /> :
        !series.data?.items?.length ? <EmptyBlock title="Серий нет" /> : (
        <ul className="grid gap-3 sm:grid-cols-2">
          {series.data.items.map((s) => (
            <Link key={s.id} to="/series/$id" params={{ id: s.id }}
              className="rounded-xl border border-border/60 bg-card/50 p-4 hover:border-primary/50">
              <p className="font-medium">{s.name}</p>
              <p className="mt-1 text-xs text-muted-foreground">{fmtDateRange(s.start_at, s.end_at)}</p>
              <div className="mt-2 flex flex-wrap items-center gap-2">
                <span className="inline-flex rounded-full bg-sky-100 px-2 py-0.5 text-xs text-sky-800">
                  {s.is_rating ? "На рейтинг" : "Без рейтинга"}
                </span>
                {Number(s.price_rub ?? 0) > 0 && (
                  <span className="inline-flex rounded-full bg-amber-100 px-2 py-0.5 text-xs text-amber-800">
                    Платно · {fmtRub(s.price_rub)}
                  </span>
                )}
              </div>
            </Link>
          ))}
        </ul>
      )}
      {series.data?.pagination && (
        <div className="mt-6 flex items-center justify-center gap-3">
          <Button
            variant="outline"
            size="sm"
            disabled={!series.data.pagination.has_previous}
            onClick={() => setPage((p) => Math.max(1, p - 1))}
          >
            Назад
          </Button>
          <span className="text-sm text-muted-foreground">
            Страница {series.data.pagination.current_page} из {series.data.pagination.total_pages}
          </span>
          <Button
            variant="outline"
            size="sm"
            disabled={!series.data.pagination.has_next}
            onClick={() => setPage((p) => p + 1)}
          >
            Далее
          </Button>
        </div>
      )}
    </PageShell>
  );
}

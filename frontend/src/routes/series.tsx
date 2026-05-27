import { createFileRoute, Link, Outlet, useLocation, useNavigate } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { seriesApi } from "@/lib/api";
import { LoadingBlock, ErrorBlock, EmptyBlock } from "@/components/site/States";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useState } from "react";
import { fmtDateRange, fmtRub } from "@/lib/format";
import { useAuthStore } from "@/lib/auth-store";
import { isClubManager } from "@/lib/roles";
import { useDebouncedValue } from "@/lib/useDebouncedValue";

export const Route = createFileRoute("/series")({ component: AllSeriesPage });

function AllSeriesPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const me = useAuthStore((s) => s.me);
  const canCreate = !!me?.club_id && isClubManager(me.club_state);

  const [q, setQ] = useState("");
  const [club, setClub] = useState("");
  const qLimit = 100;
  const clubLimit = 100;
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");
  const [showPast, setShowPast] = useState(false);
  const [showClosed, setShowClosed] = useState(false);
  const [ratingFilter, setRatingFilter] = useState<"all" | "rating" | "non_rating">("all");
  const [page, setPage] = useState(1);
  const pageSize = 10;
  const debouncedQ = useDebouncedValue(q, 150);
  const debouncedClub = useDebouncedValue(club, 150);
  const { data, isLoading, error } = useQuery({
    queryKey: ["series", "all", debouncedQ, debouncedClub, from, to, ratingFilter, showPast, showClosed, page],
    queryFn: () => seriesApi.all({
      q: debouncedQ || undefined,
      club: debouncedClub || undefined,
      from: from || undefined,
      to: to || undefined,
      is_rating: ratingFilter === "all" ? undefined : ratingFilter === "rating",
      show_past: showPast,
      show_closed: showClosed,
      limit: pageSize,
      offset: (page - 1) * pageSize,
    }),
  });

  if (location.pathname !== "/series") {
    return <Outlet />;
  }

  return (
    <PageShell>
      <PageHeader
        eyebrow="Серии" title="Все серии"
        actions={canCreate ? <Button asChild><Link to="/series/create">Создать серию</Link></Button> : null}
      />
      <div className="mb-6 grid gap-3 rounded-xl border border-border/60 bg-card/40 p-4 sm:grid-cols-2 lg:grid-cols-7">
        <div className="space-y-1"><Label className="text-xs">Название</Label><Input value={q} maxLength={qLimit} onChange={(e) => { setQ(e.target.value); setPage(1); }} placeholder="Поиск…" /><p className="text-xs text-muted-foreground">{q.length}/{qLimit}</p></div>
        <div className="space-y-1"><Label className="text-xs">Клуб</Label><Input value={club} maxLength={clubLimit} onChange={(e) => { setClub(e.target.value); setPage(1); }} placeholder="Название клуба…" /><p className="text-xs text-muted-foreground">{club.length}/{clubLimit}</p></div>
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

      {isLoading ? <LoadingBlock /> : error ? <ErrorBlock error={error} /> :
        !data?.items?.length ? <EmptyBlock title="Серии не найдены" /> : (
        <>
        <div className="grid gap-3 sm:grid-cols-2">
          {data.items.map((s) => (
            <div
              key={s.id}
              role="button"
              tabIndex={0}
              onClick={() => navigate({ to: "/series/$id", params: { id: s.id } })}
              onKeyDown={(e) => {
                if (e.key === "Enter" || e.key === " ") {
                  e.preventDefault();
                  navigate({ to: "/series/$id", params: { id: s.id } });
                }
              }}
              className="group flex cursor-pointer items-start justify-between gap-4 rounded-xl border border-border/60 bg-card/50 p-5 hover:border-primary/50 hover:bg-card"
            >
              <div>
                <h3 className="font-display text-lg font-semibold group-hover:text-primary">{s.name}</h3>
                <p className="mt-1 text-xs text-muted-foreground">{fmtDateRange(s.start_at, s.end_at)}</p>
                {s.club_name && (
                  <Link
                    to="/clubs/$id"
                    params={{ id: s.club_id }}
                    onClick={(e) => e.stopPropagation()}
                    className="mt-1 inline-block text-xs text-accent hover:underline">
                    {s.club_name}
                  </Link>
                )}
              </div>
              <div className="flex shrink-0 flex-col items-end gap-1">
                {s.is_rating && (
                  <span className="rounded-full bg-sky-100 px-2 py-0.5 text-xs text-sky-800">
                    На рейтинг
                  </span>
                )}
                {s.price_rub > 0 && (
                  <span className="rounded-full bg-amber-100 px-2 py-0.5 text-xs text-amber-800">Платно · {fmtRub(s.price_rub)}</span>
                )}
              </div>
            </div>
          ))}
        </div>
        <div className="mt-5 flex items-center justify-between text-sm text-muted-foreground">
          <span>Страница {data.pagination.current_page} из {data.pagination.total_pages}</span>
          <div className="flex gap-2">
            <Button variant="outline" size="sm" disabled={!data.pagination.has_previous || isLoading} onClick={() => setPage((p) => Math.max(1, p - 1))}>Назад</Button>
            <Button variant="outline" size="sm" disabled={!data.pagination.has_next || isLoading} onClick={() => setPage((p) => p + 1)}>Далее</Button>
          </div>
        </div>
        </>
      )}
    </PageShell>
  );
}

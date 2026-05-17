import { createFileRoute, Link, Outlet, useLocation, useNavigate } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { seriesApi } from "@/lib/api";
import { LoadingBlock, ErrorBlock, EmptyBlock } from "@/components/site/States";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { useState, useMemo } from "react";
import { fmtDateRange } from "@/lib/format";
import { useAuthStore } from "@/lib/auth-store";
import { isClubManager } from "@/lib/roles";
import { includesCaseInsensitive } from "@/lib/search";

export const Route = createFileRoute("/series")({ component: AllSeriesPage });

function AllSeriesPage() {
  const location = useLocation();
  const navigate = useNavigate();
  const me = useAuthStore((s) => s.me);
  const canCreate = !!me?.club_id && isClubManager(me.club_state);

  const [q, setQ] = useState("");
  const [club, setClub] = useState("");
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");
  const [showPast, setShowPast] = useState(false);
  const [showClosed, setShowClosed] = useState(false);
  const { data, isLoading, error } = useQuery({
    queryKey: ["series", "all", showPast, showClosed],
    queryFn: () => seriesApi.all(200, 0, showPast, showClosed),
  });

  const items = useMemo(() => {
    let xs = data?.items ?? [];
    if (q.trim()) {
      xs = xs.filter((s) => includesCaseInsensitive(s.name, q));
    }
    if (club.trim()) {
      xs = xs.filter((s) => includesCaseInsensitive(s.club_name, club));
    }
    if (from) xs = xs.filter((s) => new Date(s.end_at) >= new Date(from));
    if (to) xs = xs.filter((s) => new Date(s.start_at) <= new Date(to));
    return xs;
  }, [data, q, club, from, to]);

  if (location.pathname !== "/series") {
    return <Outlet />;
  }

  return (
    <PageShell>
      <PageHeader
        eyebrow="Серии" title="Все серии"
        actions={canCreate ? <Button asChild><Link to="/series/create">Создать серию</Link></Button> : null}
      />
      <div className="mb-6 grid gap-3 rounded-xl border border-border/60 bg-card/40 p-4 sm:grid-cols-2 lg:grid-cols-6">
        <div className="space-y-1"><Label className="text-xs">Название</Label><Input value={q} onChange={(e) => setQ(e.target.value)} placeholder="Поиск…" /></div>
        <div className="space-y-1"><Label className="text-xs">Клуб</Label><Input value={club} onChange={(e) => setClub(e.target.value)} placeholder="Название клуба…" /></div>
        <div className="space-y-1"><Label className="text-xs">С</Label><Input type="date" value={from} onChange={(e) => setFrom(e.target.value)} /></div>
        <div className="space-y-1"><Label className="text-xs">По</Label><Input type="date" value={to} onChange={(e) => setTo(e.target.value)} /></div>
        <label className="flex items-end gap-2 pb-2 text-sm text-muted-foreground">
          <Checkbox checked={showPast} onCheckedChange={(v) => setShowPast(!!v)} />
          Показать прошедшие
        </label>
        <label className="flex items-end gap-2 pb-2 text-sm text-muted-foreground">
          <Checkbox checked={showClosed} onCheckedChange={(v) => setShowClosed(!!v)} />
          Показать закрытые
        </label>
      </div>

      {isLoading ? <LoadingBlock /> : error ? <ErrorBlock error={error} /> :
        items.length === 0 ? <EmptyBlock title="Серии не найдены" /> : (
        <div className="grid gap-3 sm:grid-cols-2">
          {items.map((s) => (
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
              <span className="shrink-0 rounded-full bg-secondary px-2 py-0.5 text-xs">{s.games_count} игр</span>
            </div>
          ))}
        </div>
      )}
    </PageShell>
  );
}

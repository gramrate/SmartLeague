import { createFileRoute, Link } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { clubsApi, seriesApi } from "@/lib/api";
import { LoadingBlock, EmptyBlock } from "@/components/site/States";
import { useMemo, useState } from "react";
import { Button } from "@/components/ui/button";

export const Route = createFileRoute("/clubs/$id/games")({ component: ClubGamesPage });

function ClubGamesPage() {
  const { id } = Route.useParams();
  const [page, setPage] = useState(1);
  const pageSize = 20;
  const club = useQuery({ queryKey: ["club", id], queryFn: () => clubsApi.get(id) });
  const series = useQuery({ queryKey: ["club", id, "series"], queryFn: () => clubsApi.series(id) });

  // Aggregate games across series
  const seriesIds = useMemo(() => (series.data?.items ?? []).map((s) => s.id), [series.data]);
  const games = useQuery({
    queryKey: ["club", id, "all-games", seriesIds],
    enabled: seriesIds.length > 0,
    queryFn: async () => {
      const all = await Promise.all(seriesIds.map((sid) => seriesApi.games(sid).catch(() => null)));
      return all
        .flatMap((p, i) => (p?.items ?? []).map((g) => ({
          ...g,
          _seriesId: series.data!.items![i].id,
          _seriesName: series.data!.items![i].name,
        })))
        .sort((a, b) => b.number - a.number);
    },
  });

  const totalPages = Math.max(1, Math.ceil((games.data?.length ?? 0) / pageSize));
  const visibleGames = (games.data ?? []).slice((page - 1) * pageSize, page * pageSize);

  return (
    <PageShell>
      <PageHeader eyebrow={club.data?.name ?? "Клуб"} title="Все игры" />
      {series.isLoading || games.isLoading ? <LoadingBlock /> :
        !games.data?.length ? <EmptyBlock title="Игр нет" /> : (
        <ul className="grid gap-2 sm:grid-cols-2">
          {visibleGames.map((g) => (
            <Link key={g.id} to="/game/$id" params={{ id: g.id }}
              className="rounded-xl border border-border/60 bg-card/50 p-4 hover:border-primary/50">
              <p className="font-medium">{g.name || `Игра #${g.number}`}</p>
              <p className="mt-1 text-xs text-muted-foreground">
                <Link to="/series/$id" params={{ id: g._seriesId }} className="hover:underline" onClick={(e) => e.stopPropagation()}>
                  {g._seriesName}
                </Link>
                {" · #"}
                {g.number}
              </p>
            </Link>
          ))}
        </ul>
      )}
      {!!games.data?.length && (
        <div className="mt-6 flex items-center justify-center gap-3">
          <Button
            variant="outline"
            size="sm"
            disabled={page <= 1}
            onClick={() => setPage((p) => Math.max(1, p - 1))}
          >
            Назад
          </Button>
          <span className="text-sm text-muted-foreground">
            Страница {page} из {totalPages}
          </span>
          <Button
            variant="outline"
            size="sm"
            disabled={page >= totalPages}
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
          >
            Далее
          </Button>
        </div>
      )}
    </PageShell>
  );
}

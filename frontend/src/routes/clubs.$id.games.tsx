import { createFileRoute, Link } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { clubsApi } from "@/lib/api";
import { LoadingBlock, EmptyBlock } from "@/components/site/States";
import { useState } from "react";
import { Button } from "@/components/ui/button";

export const Route = createFileRoute("/clubs/$id/games")({ component: ClubGamesPage });

function ClubGamesPage() {
  const { id } = Route.useParams();
  const [page, setPage] = useState(1);
  const pageSize = 20;
  const offset = (page - 1) * pageSize;
  const club = useQuery({ queryKey: ["club", id], queryFn: () => clubsApi.get(id) });
  const games = useQuery({
    queryKey: ["club", id, "all-games", page],
    queryFn: () => clubsApi.games(id, { limit: pageSize, offset }),
  });

  const totalPages = games.data?.pagination.total_pages ?? 1;
  const visibleGames = games.data?.items ?? [];

  return (
    <PageShell>
      <PageHeader eyebrow={club.data?.name ?? "Клуб"} title="Все игры" />
      {games.isLoading ? <LoadingBlock /> :
        !visibleGames.length ? <EmptyBlock title="Игр нет" /> : (
        <ul className="grid gap-2 sm:grid-cols-2">
          {visibleGames.map((g) => (
            <Link key={g.id} to="/game/$id" params={{ id: g.id }}
              className="rounded-xl border border-border/60 bg-card/50 p-4 hover:border-primary/50">
              <p className="font-medium">{g.name || `Игра #${g.number}`}</p>
              <p className="mt-1 text-xs text-muted-foreground">
                <Link to="/series/$id" params={{ id: g.series_id }} className="hover:underline" onClick={(e) => e.stopPropagation()}>
                  {g.series_name}
                </Link>
                {" · #"}
                {g.number}
              </p>
            </Link>
          ))}
        </ul>
      )}
      {!!visibleGames.length && (
        <div className="mt-6 flex items-center justify-center gap-3">
          <Button
            variant="outline"
            size="sm"
            disabled={!games.data?.pagination.has_previous}
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
            disabled={!games.data?.pagination.has_next}
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
          >
            Далее
          </Button>
        </div>
      )}
    </PageShell>
  );
}

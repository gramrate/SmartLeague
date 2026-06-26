import { createFileRoute, Link } from "@tanstack/react-router";
import { PageShell, PageHeader } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { seriesApi } from "@/lib/api";
import { LoadingBlock, ErrorBlock, EmptyBlock } from "@/components/site/States";
import { displayUserName } from "@/lib/roles";
import { ListPagination } from "@/components/site/ListPagination";
import { Button } from "@/components/ui/button";
import { useState, useMemo } from "react";

export const Route = createFileRoute("/series/$id/leaderboard")({ component: LeaderboardPage });

const PAGE_SIZE = 50;

function LeaderboardPage() {
  const { id } = Route.useParams();
  const [offset, setOffset] = useState(0);

  const series = useQuery({
    queryKey: ["series", id],
    queryFn: () => seriesApi.get(id),
  });

  const leaderboard = useQuery({
    queryKey: ["series", id, "leaderboard", offset],
    queryFn: () => seriesApi.leaderboard(id, PAGE_SIZE, offset),
  });

  const participants = useQuery({
    queryKey: ["series", id, "participants"],
    queryFn: () => seriesApi.participants(id, 200, 0),
  });

  const byId = useMemo(
    () => new Map((participants.data?.items ?? []).map((u) => [u.id, u])),
    [participants.data],
  );

  const globalOffset = offset;

  return (
    <PageShell>
      <PageHeader
        eyebrow={series.data?.name ?? "Серия"}
        title="Таблица лидеров"
        actions={
          <Button variant="outline" asChild>
            <Link to="/series/$id" params={{ id }}>← Назад</Link>
          </Button>
        }
      />

      <section className="rounded-2xl border border-border/60 bg-card/60 p-6">
        {leaderboard.isLoading ? (
          <LoadingBlock />
        ) : leaderboard.error ? (
          <ErrorBlock error={leaderboard.error} />
        ) : !leaderboard.data?.items?.length ? (
          <EmptyBlock title="Пока нет результатов" />
        ) : (
          <>
            <ol className="divide-y divide-border/40">
              {leaderboard.data.items.map((r, i) => {
                const rank = globalOffset + i + 1;
                return (
                  <li key={r.profile_id ?? r.guest_id ?? i} className="flex items-center justify-between py-3">
                    <span className="flex items-center gap-3">
                      <span className="grid h-7 w-7 flex-shrink-0 place-items-center rounded-full bg-primary/15 text-xs font-bold text-primary">
                        {rank}
                      </span>
                      {r.profile_id ? (
                        <Link
                          to="/user/$id"
                          params={{ id: r.profile_id }}
                          className="max-w-[280px] truncate hover:text-primary"
                        >
                          {displayUserName(byId.get(r.profile_id) ?? { id: r.profile_id })}
                        </Link>
                      ) : (
                        <span className="max-w-[280px] truncate text-muted-foreground" title="Гость">
                          {r.guest_nickname ?? "Гость"}
                        </span>
                      )}
                    </span>
                    <span className="font-mono text-sm">{r.points.toFixed(2)}</span>
                  </li>
                );
              })}
            </ol>

            {leaderboard.data.pagination && (
              <ListPagination
                pagination={leaderboard.data.pagination}
                isLoading={leaderboard.isFetching}
                onPrevious={() => setOffset((o) => Math.max(0, o - PAGE_SIZE))}
                onNext={() => setOffset((o) => o + PAGE_SIZE)}
              />
            )}
          </>
        )}
      </section>
    </PageShell>
  );
}

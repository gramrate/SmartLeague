import { createFileRoute, Link } from "@tanstack/react-router";
import { PageHeader, PageShell } from "@/components/site/PageShell";
import { useQuery } from "@tanstack/react-query";
import { clubsApi } from "@/lib/api";
import { EmptyBlock, LoadingBlock } from "@/components/site/States";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useDebouncedValue } from "@/lib/useDebouncedValue";
import { displayUserName } from "@/lib/roles";
import { ClubState } from "@/types/api";
import { RoleBadge } from "@/components/site/RoleBadge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

export const Route = createFileRoute("/clubs/$id/members")({ component: ClubMembersPage });

function ClubMembersPage() {
  const { id } = Route.useParams();
  const [q, setQ] = useState("");
  const qLimit = 50;
  const [clubStateFilter, setClubStateFilter] = useState<"all" | "leader" | "resident" | "member">("all");
  const [page, setPage] = useState(1);
  const limit = 15;
  const offset = (page - 1) * limit;
  const debouncedQ = useDebouncedValue(q, 150);

  const club = useQuery({ queryKey: ["club", id], queryFn: () => clubsApi.get(id) });
  const members = useQuery({
    queryKey: ["club", id, "members", debouncedQ, clubStateFilter, page],
    queryFn: () => clubsApi.members(id, {
      q: debouncedQ || undefined,
      club_state:
        clubStateFilter === "all"
          ? undefined
          : clubStateFilter === "leader"
            ? 3
            : clubStateFilter === "resident"
              ? 2
              : 1,
      limit,
      offset,
    }),
  });

  return (
    <PageShell>
      <PageHeader
        eyebrow={club.data?.name ?? "Клуб"}
        title="Участники клуба"
        actions={<Button variant="outline" asChild><Link to="/clubs/$id" params={{ id }}>К клубу</Link></Button>}
      />

      <div className="mb-4 rounded-xl border border-border/60 bg-card/40 p-4">
        <div className="grid gap-3 sm:grid-cols-2">
          <div className="space-y-1">
          <Label className="text-xs">Поиск</Label>
          <Input
            value={q}
            maxLength={qLimit}
            onChange={(e) => {
              setQ(e.target.value);
              setPage(1);
            }}
            placeholder="Имя или никнейм..."
          />
          <p className="text-xs text-muted-foreground">{q.length}/{qLimit}</p>
          </div>
          <div className="space-y-1">
            <Label className="text-xs">Роль</Label>
            <Select
              value={clubStateFilter}
              onValueChange={(v: "all" | "leader" | "resident" | "member") => {
                setClubStateFilter(v);
                setPage(1);
              }}
            >
              <SelectTrigger><SelectValue /></SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Все роли</SelectItem>
                <SelectItem value="leader">Лидер (и президент)</SelectItem>
                <SelectItem value="resident">Резидент</SelectItem>
                <SelectItem value="member">Участник</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </div>

      {members.isLoading ? <LoadingBlock /> :
        !members.data?.items?.length ? <EmptyBlock title="Участники не найдены" /> : (
          <>
            <ul className="space-y-2 rounded-2xl border border-border/60 bg-card/60 p-4">
              {members.data.items.map((m) => (
                <Link
                  key={m.id}
                  to="/user/$id"
                  params={{ id: m.id }}
                  className="flex items-center justify-between gap-2 rounded-lg border border-border/40 bg-background/40 px-3 py-2 text-sm hover:border-primary/50"
                >
                  <span className="truncate hover:text-primary">{displayUserName(m)}</span>
                  {m.club_id === id && (m.club_state ?? ClubState.None) !== ClubState.None ? (
                    <RoleBadge state={(m.club_state ?? ClubState.None) as ClubState} />
                  ) : (
                    <span />
                  )}
                </Link>
              ))}
            </ul>
            <div className="mt-4 flex items-center justify-between text-sm text-muted-foreground">
              <span>Страница {members.data.pagination.current_page} из {members.data.pagination.total_pages}</span>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  disabled={!members.data.pagination.has_previous}
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                >
                  Назад
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  disabled={!members.data.pagination.has_next}
                  onClick={() => setPage((p) => p + 1)}
                >
                  Далее
                </Button>
              </div>
            </div>
          </>
      )}
    </PageShell>
  );
}

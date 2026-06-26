import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { PageHeader, PageShell } from "@/components/site/PageShell";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { usersApi } from "@/lib/api";
import { EmptyBlock, LoadingBlock } from "@/components/site/States";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { fmtDateRange, fmtRub } from "@/lib/format";
import { JudgeRole } from "@/types/api";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useDebouncedValue } from "@/lib/useDebouncedValue";
import { ListPagination } from "@/components/site/ListPagination";

export const Route = createFileRoute("/user/$id/series")({ component: UserSeriesPage });

function UserSeriesPage() {
  const { id } = Route.useParams();
  const navigate = useNavigate();
  const qc = useQueryClient();
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

  const user = useQuery({ queryKey: ["user", id], queryFn: () => usersApi.get(id) });
  const series = useQuery({
    queryKey: [
      "user",
      id,
      "series",
      debouncedQ,
      from,
      to,
      ratingFilter,
      showPast,
      showClosed,
      page,
    ],
    queryFn: () =>
      usersApi.series(id, {
        q: debouncedQ || undefined,
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
      <PageHeader
        eyebrow={user.data?.nickname ?? "Игрок"}
        title="Все серии игрока"
        actions={
          <Button variant="outline" onClick={() => {
            qc.invalidateQueries({ queryKey: ["user", id] });
            void navigate({ to: "/user/$id", params: { id } });
          }}>
            К профилю
          </Button>
        }
      />
      <div className="mb-6 grid gap-3 rounded-xl border border-border/60 bg-card/40 p-4 sm:grid-cols-2 lg:grid-cols-6">
        <div className="space-y-1">
          <Label className="text-xs">Название</Label>
          <Input
            value={q}
            maxLength={qLimit}
            onChange={(e) => {
              setQ(e.target.value);
              setPage(1);
            }}
            placeholder="Поиск…"
          />
          <p className="text-xs text-muted-foreground">
            {q.length}/{qLimit}
          </p>
        </div>
        <div className="space-y-1">
          <Label className="text-xs">С</Label>
          <Input
            type="date"
            value={from}
            onChange={(e) => {
              setFrom(e.target.value);
              setPage(1);
            }}
          />
        </div>
        <div className="space-y-1">
          <Label className="text-xs">По</Label>
          <Input
            type="date"
            value={to}
            onChange={(e) => {
              setTo(e.target.value);
              setPage(1);
            }}
          />
        </div>
        <div className="space-y-1">
          <Label className="text-xs">Тип игр</Label>
          <Select
            value={ratingFilter}
            onValueChange={(v: "all" | "rating" | "non_rating") => {
              setRatingFilter(v);
              setPage(1);
            }}
          >
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Все игры</SelectItem>
              <SelectItem value="rating">На рейтинг</SelectItem>
              <SelectItem value="non_rating">Без рейтинга</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <label className="flex items-end gap-2 pb-2 text-sm text-muted-foreground">
          <Checkbox
            checked={showPast}
            onCheckedChange={(v) => {
              setShowPast(!!v);
              setPage(1);
            }}
          />
          Показать прошедшие
        </label>
        <label className="flex items-end gap-2 pb-2 text-sm text-muted-foreground">
          <Checkbox
            checked={showClosed}
            onCheckedChange={(v) => {
              setShowClosed(!!v);
              setPage(1);
            }}
          />
          Показать закрытые
        </label>
      </div>

      {series.isLoading ? (
        <LoadingBlock />
      ) : !series.data?.items?.length ? (
        <EmptyBlock title="Серий нет" />
      ) : (
        <ul className="grid gap-3 sm:grid-cols-2">
          {series.data.items.map((s) => (
            <Link
              key={s.id}
              to="/series/$id"
              params={{ id: s.id }}
              className="rounded-xl border border-border/60 bg-card/50 p-4 hover:border-primary/50"
            >
              <p className="break-words font-medium">{s.name}</p>
              <p className="mt-1 text-xs text-muted-foreground">
                {fmtDateRange(s.start_at, s.end_at)}
              </p>
              <div className="mt-2 flex flex-wrap items-center gap-2">
                {s.is_tournament && (
                  <span className="inline-flex rounded-full bg-purple-100 px-2 py-0.5 text-xs text-purple-800 dark:bg-purple-900/40 dark:text-purple-300">
                    Турнир
                  </span>
                )}
                {s.is_judge && (
                  <span className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${s.judge_role === JudgeRole.Main ? "bg-yellow-100 text-yellow-800 dark:bg-yellow-900/40 dark:text-yellow-300" : "bg-slate-100 text-slate-700 dark:bg-slate-800 dark:text-slate-300"}`}>
                    {s.judge_role === JudgeRole.Main ? "Главный судья" : "Боковой судья"}
                  </span>
                )}
                {s.is_rating && !s.is_tournament && (
                  <span className="inline-flex rounded-full bg-sky-100 px-2 py-0.5 text-xs text-sky-800">
                    На рейтинг
                  </span>
                )}
                {s.is_club_only && (
                  <span className="inline-flex rounded-full bg-teal-100 px-2 py-0.5 text-xs text-teal-800 dark:bg-teal-900/40 dark:text-teal-300">
                    Для участников клуба
                  </span>
                )}
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
        <ListPagination
          pagination={series.data.pagination}
          isLoading={series.isLoading}
          onPrevious={() => setPage((p) => Math.max(1, p - 1))}
          onNext={() => setPage((p) => p + 1)}
        />
      )}
    </PageShell>
  );
}

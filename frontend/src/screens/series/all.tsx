import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { listAllSeries } from "../../api/series";

function toDateInputValue(v: string) {
  return new Date(v).toISOString().slice(0, 10);
}

export function AllSeriesPage() {
  const [nameFilter, setNameFilter] = useState("");
  const [clubFilter, setClubFilter] = useState("all");
  const [dateFrom, setDateFrom] = useState("");
  const [dateTo, setDateTo] = useState("");
  const [gamesFilter, setGamesFilter] = useState<"all" | "with_games" | "without_games">("all");

  const q = useQuery({
    queryKey: ["series", "all"],
    queryFn: () => listAllSeries({ limit: 100, offset: 0 })
  });

  const filtered = useMemo(() => {
    if (!q.data) return [];

    const fromDate = dateFrom ? new Date(`${dateFrom}T00:00:00`) : null;
    const toDate = dateTo ? new Date(`${dateTo}T23:59:59`) : null;

    return q.data.items.filter((s) => {
      if (nameFilter.trim() && !s.name.toLowerCase().includes(nameFilter.trim().toLowerCase())) return false;
      if (clubFilter !== "all" && s.club_id !== clubFilter) return false;

      const startAt = new Date(s.start_at);
      const endAt = new Date(s.end_at);
      if (fromDate && endAt < fromDate) return false;
      if (toDate && startAt > toDate) return false;

      if (gamesFilter === "with_games" && s.games_count <= 0) return false;
      if (gamesFilter === "without_games" && s.games_count > 0) return false;

      return true;
    });
  }, [q.data, nameFilter, clubFilter, dateFrom, dateTo, gamesFilter]);

  if (q.isLoading) return <div>Loading...</div>;
  if (q.isError) return <div>Failed to load series</div>;
  if (!q.data) return <div>No data</div>;

  return (
    <div className="space-y-4">
      <h1 className="text-xl font-semibold">All series</h1>

      <div className="grid gap-3 rounded border bg-white p-3 md:grid-cols-5">
        <input
          className="rounded border px-3 py-2 text-sm"
          placeholder="Search by series name"
          value={nameFilter}
          onChange={(e) => setNameFilter(e.target.value)}
        />

        <select className="rounded border px-3 py-2 text-sm" value={clubFilter} onChange={(e) => setClubFilter(e.target.value)}>
          <option value="all">All clubs</option>
          {Array.from(new Map(q.data.items.map((s) => [s.club_id, s.club_name])).entries()).map(([clubId, clubName]) => (
            <option key={clubId} value={clubId}>
              {clubName}
            </option>
          ))}
        </select>

        <input className="rounded border px-3 py-2 text-sm" type="date" value={dateFrom} onChange={(e) => setDateFrom(e.target.value)} />
        <input className="rounded border px-3 py-2 text-sm" type="date" value={dateTo} onChange={(e) => setDateTo(e.target.value)} />

        <select className="rounded border px-3 py-2 text-sm" value={gamesFilter} onChange={(e) => setGamesFilter(e.target.value as "all" | "with_games" | "without_games")}>
          <option value="all">All games</option>
          <option value="with_games">With games</option>
          <option value="without_games">Without games</option>
        </select>
      </div>

      <div className="grid gap-3">
        {filtered.map((s) => (
          <Link key={s.id} to={`/series/${s.id}`} className="rounded border bg-white p-4 hover:border-gray-400">
            <div className="font-medium">{s.name}</div>
            <div className="mt-1 text-xs text-gray-600">Club: {s.club_name}</div>
            <div className="mt-1 text-xs text-gray-600">
              {new Date(s.start_at).toLocaleString()} - {new Date(s.end_at).toLocaleString()}
            </div>
            <div className="mt-1 text-xs text-gray-600">
              {toDateInputValue(s.start_at)} - {toDateInputValue(s.end_at)} | Games: {s.games_count}
            </div>
            {s.description ? <div className="mt-2 text-sm text-gray-700 line-clamp-2">{s.description}</div> : null}
          </Link>
        ))}
      </div>
    </div>
  );
}

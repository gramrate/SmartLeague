import { useQuery } from "@tanstack/react-query";
import { Link, useParams } from "react-router-dom";
import { getGameFull } from "../../api/games";
import { getSeriesParticipants } from "../../api/series";
import { BackButton } from "../../shared/backButton";

export function GameDetailPage() {
  const { id } = useParams();
  const gameId = id!;

  const q = useQuery({ queryKey: ["game", gameId, "full"], queryFn: () => getGameFull(gameId), staleTime: 60_000 });
  const participantsQ = useQuery({
    queryKey: ["series", q.data?.series_id, "participants", { limit: 100, offset: 0 }],
    queryFn: () => getSeriesParticipants(q.data!.series_id, { limit: 100, offset: 0 }),
    enabled: !!q.data?.series_id
  });

  if (q.isLoading) return <div>Loading...</div>;
  if (q.isError) return <div>Failed to load game</div>;
  if (!q.data) return <div>No data</div>;
  const game = (q.data as any).game ?? q.data;
  const participantIDs: string[] = Array.isArray((q.data as any).participant_ids) ? (q.data as any).participant_ids : [];
  const results: Array<{ profile_id: string; place?: number | null; total_points: number; yellow_cards: number; removed: number }> = Array.isArray(
    (q.data as any).results
  )
    ? (q.data as any).results
    : [];
  const nicknameById = new Map((participantsQ.data?.items ?? []).map((p) => [p.id, p.nickname || p.name]));
  const sortedResults = [...results].sort((a, b) => {
    const ap = a.place ?? Number.MAX_SAFE_INTEGER;
    const bp = b.place ?? Number.MAX_SAFE_INTEGER;
    return ap - bp;
  });

  return (
    <div className="space-y-4">
      <BackButton />
      <div className="rounded bg-white p-6 shadow">
        <h1 className="text-xl font-semibold">{game.name ?? "Game"}</h1>
        <div className="mt-2 text-sm text-gray-700">Game #{game.number ?? "-"}</div>
        <Link className="mt-2 inline-block text-sm text-blue-700 hover:underline" to={`/series/${game.series_id}`}>
          Open series
        </Link>
        {game.description ? <div className="mt-2 whitespace-pre-wrap text-sm">{game.description}</div> : null}
      </div>
      <div className="rounded bg-white p-6 shadow">
        <h2 className="text-lg font-semibold">Participants</h2>
        <div className="mt-3 grid gap-2 md:grid-cols-2">
          {participantIDs.map((pid) => (
            <div key={pid} className="rounded border px-3 py-2 text-sm">
              {nicknameById.get(pid) ?? pid}
            </div>
          ))}
        </div>
      </div>
      <div className="rounded bg-white p-6 shadow">
        <h2 className="text-lg font-semibold">Results</h2>
        <div className="mt-3 overflow-auto">
          <table className="min-w-full text-left text-sm">
            <thead>
              <tr className="border-b text-xs text-gray-600">
                <th className="py-2 pr-4">Place</th>
                <th className="py-2 pr-4">Player</th>
                <th className="py-2 pr-4">Total</th>
                <th className="py-2 pr-4">Yellow</th>
                <th className="py-2 pr-4">Removed</th>
              </tr>
            </thead>
            <tbody>
              {sortedResults.map((r) => (
                <tr key={r.profile_id} className="border-b">
                  <td className="py-2 pr-4">{r.place ?? "-"}</td>
                  <td className="py-2 pr-4">{nicknameById.get(r.profile_id) ?? r.profile_id}</td>
                  <td className="py-2 pr-4">{r.total_points}</td>
                  <td className="py-2 pr-4">{r.yellow_cards}</td>
                  <td className="py-2 pr-4">{r.removed}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

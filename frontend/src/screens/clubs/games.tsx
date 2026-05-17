import { useQuery } from "@tanstack/react-query";
import { Link, useParams } from "react-router-dom";
import { listGames } from "../../api/games";
import { getClubSeries } from "../../api/series";
import { BackButton } from "../../shared/backButton";
import { GameStatus } from "../../types/enums";

type ClubGameItem = {
  id: string;
  series_id: string;
  series_name: string;
  name: string;
  number: number;
  status: GameStatus;
};

export function ClubGamesPage() {
  const { id } = useParams();
  const clubId = id!;

  const q = useQuery({
    queryKey: ["club", clubId, "games", "all"],
    queryFn: async () => {
      const seriesResp = await getClubSeries(clubId, { limit: 100, offset: 0 });
      const perSeries = await Promise.all(
        seriesResp.items.map(async (series) => {
          const gamesResp = await listGames(series.id, { limit: 100, offset: 0 });
          return gamesResp.items.map((g) => ({
            id: g.id,
            series_id: series.id,
            series_name: series.name,
            name: g.name,
            number: g.number,
            status: g.status
          }));
        })
      );

      return perSeries
        .flat()
        .sort((a, b) => {
          if (a.status !== b.status) return a.status - b.status;
          return b.number - a.number;
        }) as ClubGameItem[];
    }
  });

  if (q.isLoading) return <div>Loading...</div>;
  if (q.isError) return <div>Failed to load games</div>;
  if (!q.data) return <div>No data</div>;

  return (
    <div className="space-y-4">
      <BackButton />
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-semibold">Club games</h1>
        <Link className="text-sm text-blue-600" to={`/clubs/${clubId}`}>
          Back to club
        </Link>
      </div>

      <div className="grid gap-3">
        {q.data.map((g) => (
          <Link key={g.id} to={`/game/${g.id}`} className="rounded border bg-white p-4 hover:border-gray-400">
            <div className="font-medium">
              #{g.number} {g.name || "Untitled game"}
            </div>
            <div className="mt-1 text-xs text-gray-600">Series: {g.series_name}</div>
          </Link>
        ))}
      </div>

      {q.data.length === 0 ? <div className="text-sm text-gray-600">No games yet</div> : null}
    </div>
  );
}

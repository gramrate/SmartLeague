import { useQuery } from "@tanstack/react-query";
import { Link, useParams } from "react-router-dom";
import { getUserById, getUserGames } from "../../api/users";
import { BackButton } from "../../shared/backButton";

export function UserGamesPage() {
  const { id } = useParams();
  const userId = id!;

  const userQ = useQuery({ queryKey: ["user", userId], queryFn: () => getUserById(userId) });
  const gamesQ = useQuery({ queryKey: ["user", userId, "games", { limit: 100, offset: 0 }], queryFn: () => getUserGames(userId, { limit: 100, offset: 0 }) });

  return (
    <div className="space-y-4">
      <BackButton />
      <h1 className="text-xl font-semibold">{userQ.data ? `${userQ.data.nickname || userQ.data.name}: games` : "Player games"}</h1>

      {gamesQ.isLoading ? <div>Loading...</div> : null}
      {gamesQ.isError ? <div>Failed to load games</div> : null}

      {gamesQ.data ? (
        <div className="grid gap-3">
          {gamesQ.data.items.map((g) => (
            <Link key={g.id} to={`/game/${g.id}`} className="rounded border bg-white p-4 hover:border-gray-400">
              <div className="font-medium">#{g.number} {g.name || "Untitled game"}</div>
              <div className="mt-1 text-xs text-gray-600">Series: {g.series_name}</div>
            </Link>
          ))}
        </div>
      ) : null}
    </div>
  );
}

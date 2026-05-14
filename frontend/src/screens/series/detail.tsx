import { useMutation, useQuery } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router-dom";
import { getLeaderboard, getSeries, getSeriesParticipants, joinSeries, leaveSeries } from "../../api/series";
import { createGame, listGames } from "../../api/games";
import { queryClient } from "../../shared/queryClient";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { ClubState, GameStatus } from "../../types/enums";
import { useAuthStore } from "../../store/authStore";

const createGameSchema = z.object({
  name: z.string().max(200).optional(),
  number: z.coerce.number().int().min(1),
  description: z.string().max(2000).optional(),
  host_id: z.string().uuid().optional().or(z.literal("")),
  status: z.coerce.number().int().min(0).max(2)
});

type CreateGameForm = z.infer<typeof createGameSchema>;

export function SeriesDetailPage() {
  const { id } = useParams();
  const seriesId = id!;
  const navigate = useNavigate();
  const { clubId: myClubID, clubState } = useAuthStore();

  const seriesQ = useQuery({ queryKey: ["series", seriesId], queryFn: () => getSeries(seriesId), staleTime: 60_000 });
  const participantsQ = useQuery({ queryKey: ["series", seriesId, "participants", { limit: 20, offset: 0 }], queryFn: () => getSeriesParticipants(seriesId, { limit: 20, offset: 0 }) });
  const gamesQ = useQuery({ queryKey: ["series", seriesId, "games", { limit: 20, offset: 0 }], queryFn: () => listGames(seriesId, { limit: 20, offset: 0 }) });
  const leaderboardQ = useQuery({ queryKey: ["series", seriesId, "leaderboard", { limit: 20, offset: 0 }], queryFn: () => getLeaderboard(seriesId, { limit: 20, offset: 0 }) });

  const joinM = useMutation({
    mutationFn: () => joinSeries(seriesId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["series", seriesId, "participants"] });
    }
  });
  const leaveM = useMutation({
    mutationFn: () => leaveSeries(seriesId),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["series", seriesId, "participants"] });
    }
  });
  const createForm = useForm<CreateGameForm>({
    resolver: zodResolver(createGameSchema),
    mode: "onChange",
    defaultValues: { name: "", number: 1, description: "", host_id: "", status: GameStatus.Draft }
  });
  const createGameM = useMutation({
    mutationFn: (data: CreateGameForm) =>
      createGame(seriesId, {
        name: data.name?.trim() || undefined,
        number: data.number,
        description: data.description?.trim() || undefined,
        host_id: data.host_id?.trim() ? data.host_id : undefined,
        status: data.status
      }),
    onSuccess: async () => {
      createForm.reset({ name: "", number: 1, description: "", host_id: "", status: GameStatus.Draft });
      await queryClient.invalidateQueries({ queryKey: ["series", seriesId, "games"] });
    }
  });

  if (seriesQ.isLoading) return <div>Loading...</div>;
  if (seriesQ.isError) return <div>Failed to load series</div>;
  if (!seriesQ.data) return <div>No data</div>;
  const canManageSeriesGames =
    myClubID === seriesQ.data.club_id && (clubState === ClubState.Leader || clubState === ClubState.President);

  return (
    <div className="space-y-4">
      <div className="rounded bg-white p-6 shadow">
        <div className="flex items-start justify-between gap-4">
          <div>
            <h1 className="text-xl font-semibold">{seriesQ.data.name}</h1>
            <p className="mt-1 text-sm text-gray-700 whitespace-pre-wrap">{seriesQ.data.scoring_rules}</p>
          </div>
          <div className="flex gap-2">
            <button className="rounded bg-blue-600 px-3 py-2 text-sm text-white disabled:opacity-50" disabled={joinM.isPending} onClick={() => joinM.mutate()}>
              Join
            </button>
            <button className="rounded bg-gray-900 px-3 py-2 text-sm text-white disabled:opacity-50" disabled={leaveM.isPending} onClick={() => leaveM.mutate()}>
              Leave
            </button>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
        <div className="rounded bg-white p-6 shadow">
          <h2 className="text-lg font-semibold">Participants</h2>
          {participantsQ.data ? (
            <div className="mt-3 space-y-2">
              {participantsQ.data.items.map((p) => (
                <div key={p.id} className="rounded border px-3 py-2 text-sm">
                  <button className="text-blue-700 hover:underline" onClick={() => navigate(`/user/${p.id}`)}>
                    {p.name}
                  </button>{" "}
                  <span className="text-xs text-gray-500">{p.email}</span>
                </div>
              ))}
            </div>
          ) : (
            <div className="mt-3 text-sm text-gray-600">Loading...</div>
          )}
        </div>

        <div className="rounded bg-white p-6 shadow">
          <h2 className="text-lg font-semibold">Leaderboard</h2>
          {leaderboardQ.data ? (
            <div className="mt-3 space-y-2">
              {leaderboardQ.data.items.map((r) => (
                <div key={r.profile_id} className="flex items-center justify-between rounded border px-3 py-2 text-sm">
                  <button className="font-medium text-blue-700 hover:underline" onClick={() => navigate(`/user/${r.profile_id}`)}>
                    {r.profile_id}
                  </button>
                  <span>{r.points}</span>
                </div>
              ))}
            </div>
          ) : (
            <div className="mt-3 text-sm text-gray-600">Loading...</div>
          )}
        </div>
      </div>

      <div className="rounded bg-white p-6 shadow">
        <h2 className="text-lg font-semibold">Games</h2>
        {canManageSeriesGames ? (
          <form
            className="mt-3 grid gap-2 rounded border p-3"
            onSubmit={createForm.handleSubmit(async (data) => createGameM.mutateAsync(data))}
          >
            <div className="grid grid-cols-1 gap-2 md:grid-cols-3">
              <input className="rounded border px-3 py-2 text-sm" placeholder="Game name (optional)" {...createForm.register("name")} />
              <input className="rounded border px-3 py-2 text-sm" type="number" min={1} placeholder="Number" {...createForm.register("number")} />
              <select className="rounded border px-3 py-2 text-sm" {...createForm.register("status")}>
                <option value={GameStatus.Draft}>Draft</option>
                <option value={GameStatus.InProgress}>In progress</option>
                <option value={GameStatus.Finished}>Finished</option>
              </select>
            </div>
            <input className="rounded border px-3 py-2 text-sm" placeholder="Host UUID (optional)" {...createForm.register("host_id")} />
            <textarea className="rounded border px-3 py-2 text-sm" rows={2} placeholder="Description (optional)" {...createForm.register("description")} />
            <button className="w-fit rounded bg-gray-900 px-3 py-2 text-sm text-white disabled:opacity-50" disabled={!createForm.formState.isValid || createGameM.isPending}>
              Create game
            </button>
          </form>
        ) : null}
        {gamesQ.data ? (
          <div className="mt-3 space-y-2">
            {gamesQ.data.items.map((g) => (
              <div key={g.id} className="rounded border px-3 py-2 text-sm">
                <div className="font-medium">
                  #{g.number} {g.name}
                </div>
                {g.description ? <div className="text-xs text-gray-600 line-clamp-2">{g.description}</div> : null}
                <div className="mt-2 flex gap-2">
                  <button className="text-xs text-blue-700 hover:underline" onClick={() => navigate(`/game/${g.id}`)}>
                    Open
                  </button>
                  <button className="text-xs text-blue-700 hover:underline" onClick={() => navigate(`/game/${g.id}/manage`)}>
                    Manage
                  </button>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="mt-3 text-sm text-gray-600">Loading...</div>
        )}
      </div>
    </div>
  );
}

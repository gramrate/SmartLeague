import { useMutation, useQuery } from "@tanstack/react-query";
import { useEffect, useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import { getGameFull, setGameParticipants, upsertGameResults } from "../../api/games";
import { getSeriesParticipants } from "../../api/series";
import { queryClient } from "../../shared/queryClient";
import { MafiaRole } from "../../types/enums";

type Row = {
  profile_id: string;
  profile_label: string;
  place?: number;
  role?: MafiaRole;
  best_move?: string;
  first_killed: boolean;
  compensation: number;
  yellow_cards: number;
  removed: number;
  extra_points: number;
  total_points: number;
};

const rowCount = 10;

export function GameManagePage() {
  const { id } = useParams();
  const gameId = id!;
  const q = useQuery({ queryKey: ["game", gameId, "full"], queryFn: () => getGameFull(gameId) });
  const seriesId = q.data?.series_id;
  const participantsQ = useQuery({
    queryKey: ["series", seriesId, "participants", { limit: 200, offset: 0 }],
    queryFn: () => getSeriesParticipants(seriesId!, { limit: 200, offset: 0 }),
    enabled: !!seriesId
  });
  const [rows, setRows] = useState<Row[]>([]);
  const [participantError, setParticipantError] = useState<string | null>(null);
  const [activePickerRow, setActivePickerRow] = useState<number | null>(null);

  const participantIDs = useMemo<string[]>(() => {
    if (!q.data) return [];
    const raw = q.data as any;
    return Array.isArray(raw.participant_ids) ? raw.participant_ids : [];
  }, [q.data]);

  const existingResults = useMemo<any[]>(() => {
    if (!q.data) return [];
    const raw = q.data as any;
    return Array.isArray(raw.results) ? raw.results : [];
  }, [q.data]);

  useEffect(() => {
    const byProfile = new Map(existingResults.map((result) => [result.profile_id, result]));
    const labelsByID = new Map((participantsQ.data?.items ?? []).map((item) => [item.id, item.nickname || item.name]));
    const merged: Row[] = Array.from({ length: rowCount }, (_, index) => {
      const profileID = participantIDs[index] ?? "";
      const current = profileID ? byProfile.get(profileID) : undefined;
      const label = profileID ? labelsByID.get(profileID) ?? profileID : "";
      return {
        profile_id: profileID,
        profile_label: label,
        place: index + 1,
        role: current?.role ?? MafiaRole.Civilian,
        best_move: current?.best_move ?? "",
        first_killed: current?.first_killed ?? false,
        compensation: current?.compensation ?? 0,
        yellow_cards: current?.yellow_cards ?? 0,
        removed: current?.removed ?? 0,
        extra_points: current?.extra_points ?? 0,
        total_points: current?.total_points ?? 0
      };
    });
    setRows(merged);
  }, [participantIDs, existingResults, participantsQ.data?.items]);

  const setParticipantsM = useMutation({
    mutationFn: (ids: string[]) => setGameParticipants(gameId, ids),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["game", gameId, "full"] });
    }
  });

  const upsertResultsM = useMutation({
    mutationFn: (data: Row[]) => upsertGameResults(gameId, data),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["game", gameId, "full"] });
    }
  });

  function updateRow(index: number, patch: Partial<Row>) {
    setRows((prev) => prev.map((row, i) => (i === index ? { ...row, ...patch } : row)));
  }

  async function saveParticipants() {
    const ids = rows.map((row) => row.profile_id.trim());
    if (ids.some((idValue) => idValue === "")) {
      setParticipantError("Fill all 10 participant UUID fields");
      return;
    }
    if (new Set(ids).size !== rowCount) {
      setParticipantError("Participant UUIDs must be unique");
      return;
    }
    setParticipantError(null);
    await setParticipantsM.mutateAsync(ids);
  }

  async function submitResults() {
    const ids = rows.map((row) => row.profile_id.trim());
    if (ids.some((idValue) => idValue === "")) {
      setParticipantError("Fill participants before submitting results");
      return;
    }
    setParticipantError(null);
    await upsertResultsM.mutateAsync(rows.map((row) => ({ ...row, profile_id: row.profile_id.trim() })));
  }

  if (q.isLoading) return <div>Loading...</div>;
  if (q.isError) return <div>Failed to load game</div>;
  if (!q.data) return <div>No data</div>;

  const participantOptions = participantsQ.data?.items ?? [];

  return (
    <div className="space-y-4">
      <div className="rounded bg-white p-6 shadow">
        <h1 className="text-xl font-semibold">Manage game</h1>
        <div className="mt-1 text-sm text-gray-600">{q.data.name}</div>
      </div>

      <div className="rounded bg-white p-6 shadow">
        <h2 className="text-lg font-semibold">Game table (10 players)</h2>
        <div className="mt-3 overflow-auto">
          <table className="min-w-full text-left text-xs">
            <thead>
              <tr className="border-b">
                <th className="px-2 py-2">Place</th>
                <th className="px-2 py-2">User UUID</th>
                <th className="px-2 py-2">Role</th>
                <th className="px-2 py-2">Best move</th>
                <th className="px-2 py-2">First killed</th>
                <th className="px-2 py-2">Comp</th>
                <th className="px-2 py-2">Yellow</th>
                <th className="px-2 py-2">Removed</th>
                <th className="px-2 py-2">Extra</th>
                <th className="px-2 py-2">Total</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((row, index) => (
                <tr key={`row-${index}`} className="border-b">
                  <td className="px-2 py-2">
                    <input type="number" className="w-16 rounded border bg-gray-100 px-2 py-1" value={index + 1} readOnly />
                  </td>
                  <td className="px-2 py-2">
                    <div className="relative w-72">
                      <input
                        className="w-full rounded border px-2 py-1"
                        value={row.profile_label}
                        placeholder="Начни вводить никнейм"
                        onFocus={() => setActivePickerRow(index)}
                        onChange={(event) => updateRow(index, { profile_label: event.target.value, profile_id: "" })}
                      />
                      {activePickerRow === index ? (
                        <div className="absolute z-10 mt-1 max-h-48 w-full overflow-y-auto rounded border bg-white shadow">
                          {participantOptions
                            .filter((item) => {
                              const qv = row.profile_label.trim().toLowerCase();
                              if (!qv) return true;
                              const nickname = (item.nickname ?? "").toLowerCase();
                              const name = (item.name ?? "").toLowerCase();
                              return nickname.includes(qv) || name.includes(qv);
                            })
                            .map((item) => (
                              <button
                                key={item.id}
                                type="button"
                                className="block w-full px-2 py-1 text-left text-xs hover:bg-gray-100"
                                onMouseDown={(event) => {
                                  event.preventDefault();
                                  updateRow(index, { profile_id: item.id, profile_label: item.nickname || item.name });
                                  setActivePickerRow(null);
                                }}
                              >
                                {(item.nickname || item.name) + " — " + item.name}
                              </button>
                            ))}
                        </div>
                      ) : null}
                    </div>
                  </td>
                  <td className="px-2 py-2">
                    <select
                      className="w-28 rounded border px-2 py-1"
                      value={row.role ?? MafiaRole.Civilian}
                      onChange={(event) => updateRow(index, { role: event.target.value as MafiaRole, place: index + 1 })}
                    >
                      <option value={MafiaRole.Civilian}>Мирный</option>
                      <option value={MafiaRole.Mafia}>Мафия</option>
                      <option value={MafiaRole.Don}>Дон</option>
                      <option value={MafiaRole.Sheriff}>Шериф</option>
                    </select>
                  </td>
                  <td className="px-2 py-2">
                    <input
                      className="w-24 rounded border px-2 py-1"
                      placeholder="1,2,3"
                      value={row.best_move ?? ""}
                      onChange={(event) => updateRow(index, { best_move: event.target.value, place: index + 1 })}
                    />
                  </td>
                  <td className="px-2 py-2">
                    <input type="checkbox" checked={row.first_killed} onChange={(event) => updateRow(index, { first_killed: event.target.checked, place: index + 1 })} />
                  </td>
                  <td className="px-2 py-2">
                    <input
                      type="number"
                      step="0.25"
                      className="w-16 rounded border px-2 py-1"
                      value={row.compensation}
                      onChange={(event) => updateRow(index, { compensation: Number(event.target.value || 0), place: index + 1 })}
                    />
                  </td>
                  <td className="px-2 py-2">
                    <input
                      type="number"
                      className="w-16 rounded border px-2 py-1"
                      value={row.yellow_cards}
                      onChange={(event) => updateRow(index, { yellow_cards: Number(event.target.value || 0), place: index + 1 })}
                    />
                  </td>
                  <td className="px-2 py-2">
                    <input
                      type="number"
                      className="w-16 rounded border px-2 py-1"
                      value={row.removed}
                      onChange={(event) => updateRow(index, { removed: Number(event.target.value || 0), place: index + 1 })}
                    />
                  </td>
                  <td className="px-2 py-2">
                    <input
                      type="number"
                      step="0.25"
                      className="w-16 rounded border px-2 py-1"
                      value={row.extra_points}
                      onChange={(event) => updateRow(index, { extra_points: Number(event.target.value || 0), place: index + 1 })}
                    />
                  </td>
                  <td className="px-2 py-2">
                    <input
                      type="number"
                      step="0.25"
                      className="w-16 rounded border px-2 py-1"
                      value={row.total_points}
                      onChange={(event) => updateRow(index, { total_points: Number(event.target.value || 0), place: index + 1 })}
                    />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        {participantError ? <p className="mt-2 text-xs text-red-600">{participantError}</p> : null}
        <div className="mt-3 flex gap-2">
          <button className="rounded bg-gray-900 px-4 py-2 text-sm text-white disabled:opacity-50" disabled={setParticipantsM.isPending} onClick={saveParticipants}>
            Save participants
          </button>
          <button className="rounded bg-blue-600 px-4 py-2 text-sm text-white disabled:opacity-50" disabled={upsertResultsM.isPending} onClick={submitResults}>
            Submit results
          </button>
        </div>
      </div>
    </div>
  );
}

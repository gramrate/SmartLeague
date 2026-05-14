import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { createSeries } from "../../api/series";
import { GameType, SeriesStatus } from "../../types/enums";

const schema = z.object({
  name: z.string().min(1).max(200),
  scoring_rules: z.string().min(1).max(10_000),
  start_at: z.string().min(1),
  end_at: z.string().min(1),
  description: z.string().max(5_000).optional(),
  price_rub: z.coerce.number().min(0).max(100_000_000),
  is_closed: z.boolean(),
  game_type: z.coerce.number(),
  status: z.coerce.number()
});

type FormData = z.infer<typeof schema>;

export function SeriesCreatePage() {
  const navigate = useNavigate();
  const form = useForm<FormData>({
    resolver: zodResolver(schema),
    mode: "onChange",
    defaultValues: {
      name: "",
      scoring_rules: "",
      start_at: "",
      end_at: "",
      description: "",
      price_rub: 0,
      is_closed: false,
      game_type: GameType.SportMafia,
      status: SeriesStatus.Registration
    }
  });

  const mut = useMutation({
    mutationFn: (data: FormData) =>
      createSeries({
        ...data,
        description: data.description ? data.description : null
      }),
    onSuccess: (s) => navigate(`/series/${s.id}`)
  });

  return (
    <div className="max-w-2xl rounded bg-white p-6 shadow">
      <h1 className="text-xl font-semibold">Create series</h1>
      <form className="mt-4 space-y-3" onSubmit={form.handleSubmit(async (d) => mut.mutateAsync(d))}>
        <div>
          <label className="text-sm">Name</label>
          <input className="mt-1 w-full rounded border px-3 py-2" {...form.register("name")} />
        </div>
        <div>
          <label className="text-sm">Rules</label>
          <textarea className="mt-1 w-full rounded border px-3 py-2" rows={6} {...form.register("scoring_rules")} />
        </div>
        <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
          <div>
            <label className="text-sm">Start</label>
            <input className="mt-1 w-full rounded border px-3 py-2" type="datetime-local" {...form.register("start_at")} />
          </div>
          <div>
            <label className="text-sm">End</label>
            <input className="mt-1 w-full rounded border px-3 py-2" type="datetime-local" {...form.register("end_at")} />
          </div>
        </div>
        <div>
          <label className="text-sm">Price (RUB)</label>
          <input className="mt-1 w-full rounded border px-3 py-2" type="number" {...form.register("price_rub")} />
        </div>
        <div className="flex items-center gap-2">
          <input type="checkbox" {...form.register("is_closed")} />
          <span className="text-sm">Closed series (club only)</span>
        </div>
        <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
          <div>
            <label className="text-sm">Game type</label>
            <input className="mt-1 w-full rounded border px-3 py-2 bg-gray-100" value="Sport mafia" readOnly />
          </div>
          <div>
            <label className="text-sm">Status</label>
            <select className="mt-1 w-full rounded border px-3 py-2" {...form.register("status")}>
              <option value={SeriesStatus.Closed}>Closed</option>
              <option value={SeriesStatus.Registration}>Registration</option>
              <option value={SeriesStatus.ClosedRegistration}>Closed registration</option>
              <option value={SeriesStatus.Games}>Games</option>
            </select>
          </div>
        </div>
        <button className="rounded bg-gray-900 px-4 py-2 text-white disabled:opacity-50" disabled={!form.formState.isValid || mut.isPending}>
          Create
        </button>
      </form>
    </div>
  );
}

import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { createSeries } from "../../api/series";
import { GameType } from "../../types/enums";
import { BackButton } from "../../shared/backButton";

const schema = z.object({
  name: z.string().min(1).max(200),
  description: z.string().min(1).max(10_000),
  start_at: z.string().min(1),
  end_at: z.string().min(1),
  price_rub: z.coerce.number().min(0).max(100_000_000),
  is_closed: z.boolean(),
  game_type: z.coerce.number()
});

type FormData = z.infer<typeof schema>;

export function SeriesCreatePage() {
  const navigate = useNavigate();
  const form = useForm<FormData>({
    resolver: zodResolver(schema),
    mode: "onChange",
    defaultValues: {
      name: "",
      description: "",
      start_at: "",
      end_at: "",
      price_rub: 0,
      is_closed: false,
      game_type: GameType.SportMafia
    }
  });

  const mut = useMutation({
    mutationFn: (data: FormData) => createSeries(data),
    onSuccess: (s) => navigate(`/series/${s.id}`)
  });

  return (
    <div className="space-y-3">
      <BackButton />
      <div className="max-w-2xl rounded bg-white p-6 shadow">
        <h1 className="text-xl font-semibold">Create series</h1>
      <form className="mt-4 space-y-3" onSubmit={form.handleSubmit(async (d) => mut.mutateAsync(d))}>
        <div>
          <label className="text-sm">Name</label>
          <input className="mt-1 w-full rounded border px-3 py-2" {...form.register("name")} />
        </div>
        <div>
          <label className="text-sm">Описание</label>
          <textarea className="mt-1 w-full rounded border px-3 py-2" rows={6} {...form.register("description")} />
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
        <div>
          <label className="text-sm">Game type</label>
          <input className="mt-1 w-full rounded border px-3 py-2 bg-gray-100" value="Sport mafia" readOnly />
        </div>
        <button className="rounded bg-gray-900 px-4 py-2 text-white disabled:opacity-50" disabled={!form.formState.isValid || mut.isPending}>
          Create
        </button>
      </form>
      </div>
    </div>
  );
}

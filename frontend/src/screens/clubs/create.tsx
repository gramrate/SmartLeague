import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { createClub } from "../../api/clubs";
import { BackButton } from "../../shared/backButton";
import { useAuthStore } from "../../store/authStore";

const schema = z.object({
  name: z.string().min(1).max(200),
  description: z.string().max(2000).optional()
});

type FormData = z.infer<typeof schema>;

export function ClubCreatePage() {
  const navigate = useNavigate();
  const { clubId, init } = useAuthStore();
  const alreadyInClub = !!clubId;
  const form = useForm<FormData>({ resolver: zodResolver(schema), mode: "onChange", defaultValues: { name: "", description: "" } });
  const mut = useMutation({
    mutationFn: (data: FormData) => createClub({ name: data.name, description: data.description ? data.description : null }),
    onSuccess: async (club) => {
      await init();
      navigate(`/clubs/${club.id}`);
    }
  });

  return (
    <div className="space-y-3">
      <BackButton />
      <div className="max-w-xl rounded bg-white p-6 shadow">
        <h1 className="text-xl font-semibold">Create club</h1>
      {alreadyInClub ? <p className="mt-2 text-sm text-red-600">Нельзя создать клуб, так как вы уже состоите в клубе.</p> : null}
      <form className="mt-4 space-y-3" onSubmit={form.handleSubmit(async (data) => mut.mutateAsync(data))}>
        <div>
          <label className="text-sm">Name</label>
          <input className="mt-1 w-full rounded border px-3 py-2" {...form.register("name")} />
          {form.formState.errors.name ? <p className="mt-1 text-xs text-red-600">{form.formState.errors.name.message}</p> : null}
        </div>
        <div>
          <label className="text-sm">Description</label>
          <textarea className="mt-1 w-full rounded border px-3 py-2" rows={4} {...form.register("description")} />
          {form.formState.errors.description ? <p className="mt-1 text-xs text-red-600">{form.formState.errors.description.message}</p> : null}
        </div>
        <button disabled={alreadyInClub || !form.formState.isValid || mut.isPending} className="rounded bg-gray-900 px-4 py-2 text-white disabled:opacity-50">
          Create
        </button>
      </form>
      </div>
    </div>
  );
}

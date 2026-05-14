import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { registerUser } from "../../api/auth";
import { useAuthStore } from "../../store/authStore";

const schema = z.object({
  nickname: z.string().min(1).max(100).optional(),
  name: z.string().min(2).max(100),
  show_name: z.boolean().default(true),
  description: z.string().max(2000).optional(),
  email: z.string().email().min(6).max(254),
  password: z.string().min(8).max(100)
});

type FormData = z.infer<typeof schema>;

export function RegisterPage() {
  const navigate = useNavigate();
  const { init } = useAuthStore();
  const [serverError, setServerError] = useState<string | null>(null);

  const form = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { nickname: "", name: "", show_name: true, description: "", email: "", password: "" },
    mode: "onChange"
  });

  return (
    <div className="mx-auto max-w-sm rounded bg-white p-6 shadow">
      <h1 className="text-xl font-semibold">Register</h1>
      <form
        className="mt-4 space-y-3"
        onSubmit={form.handleSubmit(async (data) => {
          setServerError(null);
          try {
            await registerUser(data);
            await init();
            navigate("/clubs");
          } catch (e: any) {
            setServerError(e?.message || "registration failed");
          }
        })}
      >
        <div>
          <label className="text-sm">Name</label>
          <input className="mt-1 w-full rounded border px-3 py-2" {...form.register("name")} />
          {form.formState.errors.name ? <p className="mt-1 text-xs text-red-600">{form.formState.errors.name.message}</p> : null}
        </div>
        <div>
          <label className="text-sm">Nickname</label>
          <input className="mt-1 w-full rounded border px-3 py-2" {...form.register("nickname")} />
        </div>
        <div className="flex items-center gap-2">
          <input type="checkbox" {...form.register("show_name")} />
          <span className="text-sm">Show name</span>
        </div>
        <div>
          <label className="text-sm">Description</label>
          <textarea className="mt-1 w-full rounded border px-3 py-2" rows={3} {...form.register("description")} />
        </div>
        <div>
          <label className="text-sm">Email</label>
          <input className="mt-1 w-full rounded border px-3 py-2" {...form.register("email")} />
          {form.formState.errors.email ? <p className="mt-1 text-xs text-red-600">{form.formState.errors.email.message}</p> : null}
        </div>
        <div>
          <label className="text-sm">Password</label>
          <input className="mt-1 w-full rounded border px-3 py-2" type="password" {...form.register("password")} />
          {form.formState.errors.password ? <p className="mt-1 text-xs text-red-600">{form.formState.errors.password.message}</p> : null}
        </div>
        {serverError ? <p className="text-sm text-red-600">{serverError}</p> : null}
        <button
          disabled={!form.formState.isValid || form.formState.isSubmitting}
          className="w-full rounded bg-gray-900 px-3 py-2 text-white disabled:opacity-50"
        >
          Create account
        </button>
      </form>
    </div>
  );
}

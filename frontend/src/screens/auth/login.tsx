import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { loginUser } from "../../api/auth";
import { useAuthStore } from "../../store/authStore";

const schema = z.object({
  email: z.string().email().min(6).max(254),
  password: z.string().min(8).max(100)
});

type FormData = z.infer<typeof schema>;

export function LoginPage() {
  const navigate = useNavigate();
  const { init } = useAuthStore();
  const [serverError, setServerError] = useState<string | null>(null);

  const form = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: { email: "", password: "" }
  });

  return (
    <div className="mx-auto max-w-sm rounded bg-white p-6 shadow">
      <h1 className="text-xl font-semibold">Login</h1>
      <form
        className="mt-4 space-y-3"
        onSubmit={form.handleSubmit(async (data) => {
          setServerError(null);
          try {
            await loginUser(data);
            await init();
            navigate("/clubs");
          } catch (e: any) {
            setServerError(e?.message || "login failed");
          }
        })}
      >
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
          Sign in
        </button>
      </form>
    </div>
  );
}


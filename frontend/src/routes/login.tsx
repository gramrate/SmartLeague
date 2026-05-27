import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { PageShell } from "@/components/site/PageShell";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useAuthStore } from "@/lib/auth-store";
import { authApi, ApiError } from "@/lib/api";
import { useEffect, useState } from "react";

const schema = z.object({
  email: z.string().email("Некорректный email").max(254),
  password: z.string().min(8, "Минимум 8 символов").max(100),
});
type Form = z.infer<typeof schema>;

export const Route = createFileRoute("/login")({ component: LoginPage });

function LoginPage() {
  const navigate = useNavigate();
  const setMe = useAuthStore((s) => s.setMe);
  const me = useAuthStore((s) => s.me);
  const [serverError, setServerError] = useState<string | null>(null);
  const form = useForm<Form>({ resolver: zodResolver(schema), defaultValues: { email: "", password: "" } });

  useEffect(() => { if (me) navigate({ to: "/account" }); }, [me, navigate]);

  const onSubmit = form.handleSubmit(async (values) => {
    setServerError(null);
    try {
      const user = await authApi.login(values);
      setMe(user);
      navigate({ to: "/account" });
    } catch (e) {
      setServerError(e instanceof ApiError ? e.message : "Не удалось войти");
    }
  });

  return (
    <PageShell>
      <div className="mx-auto max-w-md">
        <div className="rounded-2xl border border-border/60 bg-card/60 p-8 shadow-[var(--shadow-card)]">
          <h1 className="font-display text-2xl font-bold">С возвращением</h1>
          <p className="mt-1 text-sm text-muted-foreground">Войдите, чтобы получить доступ к аккаунту.</p>
          <form className="mt-6 space-y-4" onSubmit={onSubmit}>
            <div className="space-y-1.5">
              <Label htmlFor="email">Почта</Label>
              <Input id="email" type="email" autoComplete="email" {...form.register("email")} maxLength={254} />
              <p className="text-xs text-muted-foreground">{(form.watch("email") || "").length}/254 · осталось {254 - (form.watch("email") || "").length}</p>
              {form.formState.errors.email && <p className="text-xs text-destructive">{form.formState.errors.email.message}</p>}
            </div>
            <div className="space-y-1.5">
              <Label htmlFor="password">Пароль</Label>
              <Input id="password" type="password" autoComplete="current-password" {...form.register("password")} maxLength={100} />
              <p className="text-xs text-muted-foreground">{(form.watch("password") || "").length}/100 · осталось {100 - (form.watch("password") || "").length}</p>
              {form.formState.errors.password && <p className="text-xs text-destructive">{form.formState.errors.password.message}</p>}
            </div>
            {serverError && (
              <div className="rounded-md border border-destructive/40 bg-destructive/10 px-3 py-2 text-sm text-destructive">{serverError}</div>
            )}
            <Button type="submit" className="w-full" disabled={form.formState.isSubmitting || (form.watch("email") || "").length > 254 || (form.watch("password") || "").length > 100}>
              {form.formState.isSubmitting ? "Входим…" : "Войти"}
            </Button>
            <p className="text-center text-sm text-muted-foreground">
              Нет аккаунта? <Link to="/register" className="text-primary hover:underline">Создать</Link>
            </p>
          </form>
        </div>
      </div>
    </PageShell>
  );
}

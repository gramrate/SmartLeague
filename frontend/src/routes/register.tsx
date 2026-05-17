import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { PageShell } from "@/components/site/PageShell";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useAuthStore } from "@/lib/auth-store";
import { authApi, ApiError } from "@/lib/api";
import { useEffect, useState } from "react";
import { Checkbox } from "@/components/ui/checkbox";

const schema = z.object({
  email: z.string().email("Некорректный email"),
  password: z.string().min(8, "Минимум 8 символов"),
  name: z.string().min(1, "Обязательно").max(100),
  nickname: z.string().max(100).optional().or(z.literal("")),
  description: z.string().max(2000).optional().or(z.literal("")),
  show_name: z.boolean(),
});
type Form = z.infer<typeof schema>;

export const Route = createFileRoute("/register")({ component: RegisterPage });

function RegisterPage() {
  const navigate = useNavigate();
  const setMe = useAuthStore((s) => s.setMe);
  const me = useAuthStore((s) => s.me);
  const [serverError, setServerError] = useState<string | null>(null);
  const form = useForm<Form>({
    resolver: zodResolver(schema),
    defaultValues: { email: "", password: "", name: "", nickname: "", description: "", show_name: true },
  });

  useEffect(() => { if (me) navigate({ to: "/account" }); }, [me, navigate]);

  const onSubmit = form.handleSubmit(async (values) => {
    setServerError(null);
    try {
      const user = await authApi.register({
        email: values.email,
        password: values.password,
        name: values.name,
        nickname: values.nickname || undefined,
        description: values.description || undefined,
        show_name: values.show_name,
      });
      setMe(user);
      navigate({ to: "/account" });
    } catch (e) {
      setServerError(e instanceof ApiError ? e.message : "Не удалось зарегистрироваться");
    }
  });

  return (
    <PageShell>
      <div className="mx-auto max-w-md">
        <div className="rounded-2xl border border-border/60 bg-card/60 p-8 shadow-[var(--shadow-card)]">
          <h1 className="font-display text-2xl font-bold">Создайте аккаунт</h1>
          <p className="mt-1 text-sm text-muted-foreground">Присоединяйтесь к SmartLeague.</p>
          <form className="mt-6 space-y-4" onSubmit={onSubmit}>
            <div className="grid grid-cols-2 gap-3">
              <Field label="Имя" error={form.formState.errors.name?.message}>
                <Input {...form.register("name")} />
              </Field>
              <Field label="Никнейм" error={form.formState.errors.nickname?.message}>
                <Input {...form.register("nickname")} />
              </Field>
            </div>
            <Field label="Почта" error={form.formState.errors.email?.message}>
              <Input type="email" autoComplete="email" {...form.register("email")} />
            </Field>
            <Field label="Пароль" error={form.formState.errors.password?.message}>
              <Input type="password" autoComplete="new-password" {...form.register("password")} />
            </Field>
            <Field label="О себе" error={form.formState.errors.description?.message}>
              <Textarea rows={3} {...form.register("description")} />
            </Field>
            <label className="flex items-center gap-2 text-sm text-muted-foreground">
              <Checkbox
                checked={form.watch("show_name")}
                onCheckedChange={(v) => form.setValue("show_name", !!v)}
              />
              Показывать мое настоящее имя публично
            </label>
            {serverError && (
              <div className="rounded-md border border-destructive/40 bg-destructive/10 px-3 py-2 text-sm text-destructive">{serverError}</div>
            )}
            <Button type="submit" className="w-full" disabled={form.formState.isSubmitting}>
              {form.formState.isSubmitting ? "Создаем…" : "Создать аккаунт"}
            </Button>
            <p className="text-center text-sm text-muted-foreground">
              Уже зарегистрированы? <Link to="/login" className="text-primary hover:underline">Войти</Link>
            </p>
          </form>
        </div>
      </div>
    </PageShell>
  );
}

function Field({ label, error, children }: { label: string; error?: string; children: React.ReactNode }) {
  return (
    <div className="space-y-1.5">
      <Label>{label}</Label>
      {children}
      {error && <p className="text-xs text-destructive">{error}</p>}
    </div>
  );
}

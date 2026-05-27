import { createFileRoute, Link } from "@tanstack/react-router";
import { PageShell } from "@/components/site/PageShell";
import { Button } from "@/components/ui/button";
import { Trophy, Users, Calendar, Search, ArrowRight } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { clubsApi, seriesApi } from "@/lib/api";
import { fmtDateRange } from "@/lib/format";

export const Route = createFileRoute("/")({
  component: Index,
});

function Index() {
  const clubs = useQuery({ queryKey: ["clubs", "preview"], queryFn: () => clubsApi.all({ limit: 6, offset: 0 }), retry: 0 });
  const series = useQuery({ queryKey: ["series", "preview"], queryFn: () => seriesApi.all({ limit: 6, offset: 0 }), retry: 0 });

  return (
    <PageShell>
      {/* Hero */}
      <section className="relative overflow-hidden rounded-3xl border border-border/60 bg-card/40 px-6 py-16 sm:px-12 sm:py-24">
        <div className="absolute inset-0 -z-10 opacity-60" style={{ background: "var(--gradient-radial)" }} />
        <div className="absolute -top-32 -right-32 -z-10 h-96 w-96 rounded-full bg-primary/20 blur-3xl" />
        <div className="max-w-3xl">
          <p className="mb-4 inline-flex items-center gap-2 rounded-full border border-primary/40 bg-primary/10 px-3 py-1 text-xs font-medium uppercase tracking-widest text-primary">
            <Trophy className="h-3 w-3" /> Лига спортивной мафии
          </p>
          <h1 className="break-words font-display text-4xl font-bold leading-tight sm:text-6xl">
            Клубы. Серии. <span className="bg-gradient-to-r from-primary to-primary-glow bg-clip-text text-transparent">Слава.</span>
          </h1>
          <p className="mt-4 max-w-xl text-base text-muted-foreground sm:text-lg">
            Изучайте активные клубы, следите за турнирными сериями,
            смотрите отдельные игры и открывайте сильнейших игроков.
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <Button size="lg" asChild>
              <Link to="/series">Смотреть серии <ArrowRight className="ml-1 h-4 w-4" /></Link>
            </Button>
            <Button size="lg" variant="outline" asChild>
              <Link to="/players">Найти игроков</Link>
            </Button>
          </div>
        </div>
      </section>

      {/* Highlights */}
      <div className="mt-10 grid gap-4 sm:grid-cols-3">
        <FeatureCard icon={<Users className="h-5 w-5" />} title="Клубы" desc="Вступайте в клубы, получайте роли, проводите турниры." to="/clubs" />
        <FeatureCard icon={<Calendar className="h-5 w-5" />} title="Серии" desc="Текущие и предстоящие соревновательные серии." to="/series" />
        <FeatureCard icon={<Search className="h-5 w-5" />} title="Игроки" desc="Поиск профилей, последних игр и истории." to="/players" />
      </div>

      {/* Clubs preview */}
      <section className="mt-14">
        <SectionHeader title="Популярные клубы" link={{ to: "/clubs", label: "Все клубы" }} />
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {clubs.data?.items?.slice(0, 6).map((c) => (
            <Link key={c.id} to="/clubs/$id" params={{ id: c.id }}
              className="group rounded-xl border border-border/60 bg-card/50 p-5 transition-all hover:border-primary/50 hover:bg-card hover:shadow-[var(--shadow-glow)]">
              <h3 className="break-words font-display text-lg font-semibold group-hover:text-primary">{c.name}</h3>
              {c.description && <p className="mt-2 line-clamp-2 break-words text-sm text-muted-foreground">{c.description}</p>}
            </Link>
          ))}
          {clubs.data?.items?.length === 0 && (
            <p className="col-span-full text-sm text-muted-foreground">Клубов пока нет.</p>
          )}
        </div>
      </section>

      {/* Series preview */}
      <section className="mt-14">
        <SectionHeader title="Последние серии" link={{ to: "/series", label: "Все серии" }} />
        <div className="grid gap-3 sm:grid-cols-2">
          {series.data?.items?.slice(0, 6).map((s) => (
            <Link key={s.id} to="/series/$id" params={{ id: s.id }}
              className="group flex items-start justify-between gap-4 rounded-xl border border-border/60 bg-card/50 p-5 transition-all hover:border-primary/50 hover:bg-card">
              <div>
                <h3 className="break-words font-display text-lg font-semibold group-hover:text-primary">{s.name}</h3>
                <p className="mt-1 text-xs text-muted-foreground">{fmtDateRange(s.start_at, s.end_at)}</p>
                {s.club_name && <p className="mt-1 text-xs text-accent">{s.club_name}</p>}
              </div>
              <span className="rounded-full bg-secondary px-2 py-0.5 text-xs">{s.games_count} игр</span>
            </Link>
          ))}
          {series.data?.items?.length === 0 && (
            <p className="col-span-full text-sm text-muted-foreground">Серий пока нет.</p>
          )}
        </div>
      </section>
    </PageShell>
  );
}

function FeatureCard({ icon, title, desc, to }: { icon: React.ReactNode; title: string; desc: string; to: string }) {
  return (
    <Link to={to} className="group flex items-start gap-4 rounded-xl border border-border/60 bg-card/50 p-5 transition-all hover:border-primary/50 hover:bg-card">
      <div className="grid h-10 w-10 place-items-center rounded-lg bg-primary/15 text-primary">{icon}</div>
      <div>
        <h3 className="font-semibold group-hover:text-primary">{title}</h3>
        <p className="mt-0.5 text-sm text-muted-foreground">{desc}</p>
      </div>
    </Link>
  );
}

function SectionHeader({ title, link }: { title: string; link: { to: string; label: string } }) {
  return (
    <div className="mb-4 flex items-end justify-between">
      <h2 className="font-display text-2xl font-semibold">{title}</h2>
      <Link to={link.to} className="text-sm text-primary hover:underline">{link.label} →</Link>
    </div>
  );
}

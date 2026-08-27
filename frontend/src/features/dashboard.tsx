import { useEffect, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { MonthNav } from "@/components/month-nav";
import { Button } from "@/components/ui/button";
import { Card, CardTitle } from "@/components/ui/card";
import { api, type Dashboard } from "@/lib/api";
import { fmt, formatDashboardTotals } from "@/lib/dashboard-copy";
import { monthTitle } from "@/lib/labels";

function nowPeriod() {
  const d = new Date();
  return { year: d.getFullYear(), month: d.getMonth() + 1 };
}

export function DashboardPage() {
  const [{ year, month }, setPeriod] = useState(nowPeriod);
  const q = useQuery({
    queryKey: ["dashboard", year, month],
    queryFn: () => api<Dashboard>(`/dashboard?year=${year}&month=${month}`),
  });
  const d = q.data;
  const [copyState, setCopyState] = useState<"idle" | "copied" | "error">("idle");

  useEffect(() => {
    if (copyState !== "copied") return;
    const id = window.setTimeout(() => setCopyState("idle"), 2000);
    return () => window.clearTimeout(id);
  }, [copyState]);

  async function copyTotals() {
    if (!d) return;
    try {
      await navigator.clipboard.writeText(formatDashboardTotals(d, year, month));
      setCopyState("copied");
    } catch {
      setCopyState("error");
    }
  }

  return (
    <div className="space-y-6">
      <header className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 className="text-3xl">Inicio</h1>
          <p className="text-sm text-muted-foreground">
            Totales para copiar a mano al informe de la sucursal. Año de servicio {d?.serviceYear ?? "—"}.
          </p>
        </div>
        <div className="flex flex-wrap items-end gap-2">
          <Button type="button" variant="outline" disabled={!d} onClick={() => void copyTotals()}>
            {copyState === "copied" ? "Copiado" : "Copiar totales"}
          </Button>
          <MonthNav year={year} month={month} onChange={(y, m) => setPeriod({ year: y, month: m })} />
        </div>
      </header>
      {copyState === "error" ? (
        <p className="text-sm text-destructive">No se pudo copiar. Copia los números a mano.</p>
      ) : null}

      {q.isError ? <p className="text-destructive">{q.error.message}</p> : null}

      <section className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <Stat title="Publicadores activos" value={fmt(d?.publishers.active)} hint={`${fmt(d?.publishers.regular)} regulares · ${fmt(d?.publishers.irregular)} irregulares`} />
        <Stat title="Participaron" value={d ? `${d.reports.shared}/${d.reports.shouldReport}` : "—"} hint={`${fmt(d?.reports.participation, 0)} % este mes`} />
        <Stat title="Estudios bíblicos" value={fmt(d?.reports.bibleStudies)} hint="Suma del mes" />
        <Stat title="Horas de precursores" value={fmt(d?.reports.pioneerHours, 1)} hint="RP, SP y PA" />
      </section>

      <section className="grid gap-4 lg:grid-cols-2">
        <Card>
          <CardTitle>Asistencia · {monthTitle(year, month)}</CardTitle>
          <dl className="mt-4 grid grid-cols-2 gap-3 text-sm">
            <div>
              <dt className="text-muted-foreground">Entre semana (promedio)</dt>
              <dd className="text-xl">{fmt(d?.attendance.midweekAvg, 1)}</dd>
            </div>
            <div>
              <dt className="text-muted-foreground">Fin de semana (promedio)</dt>
              <dd className="text-xl">{fmt(d?.attendance.weekendAvg, 1)}</dd>
            </div>
          </dl>
          <Link to="/attendance" className="mt-4 inline-block text-sm text-primary underline">
            Registrar asistencia
          </Link>
        </Card>
        <Card>
          <CardTitle>Informes pendientes</CardTitle>
          <p className="mt-1 text-sm text-muted-foreground">
            {d ? `${d.reports.reported} de ${d.reports.shouldReport} registrados` : "Cargando…"}
          </p>
          <ul className="mt-3 max-h-56 space-y-1 overflow-auto text-sm">
            {(d?.reports.missing ?? []).length === 0 ? (
              <li className="text-muted-foreground">Nadie pendiente.</li>
            ) : (
              d?.reports.missing.map((p) => (
                <li key={p.publisherId}>
                  {p.lastName}, {p.firstName}
                </li>
              ))
            )}
          </ul>
          <Link to="/reports" className="mt-4 inline-block text-sm text-primary underline">
            Abrir la cuadrícula del mes
          </Link>
        </Card>
      </section>
    </div>
  );
}

function Stat({ title, value, hint }: { title: string; value: string; hint: string }) {
  return (
    <Card>
      <div className="text-sm text-muted-foreground">{title}</div>
      <div className="mt-1 font-sans text-3xl">{value}</div>
      <div className="mt-1 text-xs text-muted-foreground">{hint}</div>
    </Card>
  );
}

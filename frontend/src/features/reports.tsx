import { useEffect, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { MonthNav } from "@/components/month-nav";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { api, type ReportRow, type ReportsResponse } from "@/lib/api";

function nowPeriod() {
  const d = new Date();
  return { year: d.getFullYear(), month: d.getMonth() + 1 };
}

export function ReportsPage() {
  const qc = useQueryClient();
  const [{ year, month }, setPeriod] = useState(nowPeriod);
  const q = useQuery({
    queryKey: ["reports", year, month],
    queryFn: () => api<ReportsResponse>(`/reports?year=${year}&month=${month}`),
  });
  const [rows, setRows] = useState<ReportRow[]>([]);

  useEffect(() => {
    if (q.data) setRows(q.data.reports);
  }, [q.data]);

  function update(id: string, patch: Partial<ReportRow>) {
    setRows((cur) => cur.map((r) => (r.publisherId === id ? { ...r, ...patch } : r)));
  }

  const save = useMutation({
    mutationFn: () =>
      api("/reports", {
        method: "PUT",
        body: JSON.stringify({
          year,
          month,
          reports: rows.map((r) => ({
            publisherId: r.publisherId,
            sharedInMinistry: r.sharedInMinistry,
            bibleStudies: Number(r.bibleStudies) || 0,
            hours: r.hourReporter || r.auxiliaryPioneer ? r.hours : null,
            auxiliaryPioneer: r.auxiliaryPioneer,
            late: r.late,
            remarks: r.remarks,
          })),
        }),
      }),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["reports"] });
      await qc.invalidateQueries({ queryKey: ["dashboard"] });
      await qc.invalidateQueries({ queryKey: ["publishers"] });
    },
  });

  function markAll(shared: boolean) {
    setRows((cur) => cur.map((r) => ({ ...r, sharedInMinistry: shared })));
  }

  return (
    <div className="space-y-6">
      <header className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 className="text-3xl">Informes</h1>
          <p className="text-sm text-muted-foreground">
            Año de servicio {q.data?.serviceYear ?? "—"}. Los estudiantes no aparecen aquí.
          </p>
        </div>
        <MonthNav year={year} month={month} onChange={(y, m) => setPeriod({ year: y, month: m })} />
      </header>

      <div className="flex flex-wrap gap-2">
        <Button variant="outline" onClick={() => markAll(true)}>
          Todos participaron
        </Button>
        <Button variant="outline" onClick={() => markAll(false)}>
          Nadie participó
        </Button>
        <Button onClick={() => save.mutate()} disabled={save.isPending}>
          Guardar mes
        </Button>
        {save.isError ? <span className="self-center text-sm text-destructive">{save.error.message}</span> : null}
        {save.isSuccess ? <span className="self-center text-sm text-primary">Guardado.</span> : null}
      </div>

      {q.data && q.data.missing.length > 0 ? (
        <p className="text-sm text-accent">
          Sin registro: {q.data.missing.map((m) => `${m.lastName}, ${m.firstName}`).join(" · ")}
        </p>
      ) : null}

      <div className="overflow-x-auto rounded-lg border border-border bg-card">
        <table className="w-full text-left text-sm">
          <thead className="bg-muted">
            <tr>
              <th className="px-3 py-2">Publicador</th>
              <th className="px-3 py-2">Participó</th>
              <th className="px-3 py-2">Estudios</th>
              <th className="px-3 py-2">PA</th>
              <th className="px-3 py-2">Horas</th>
              <th className="px-3 py-2">Tarde</th>
              <th className="px-3 py-2">Notas</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((r) => {
              const hourReporter = r.isRegularPioneer || r.isSpecialPioneer || r.auxiliaryPioneer;
              return (
                <tr key={r.publisherId} className="border-t border-border">
                  <td className="px-3 py-2">
                    {r.lastName}, {r.firstName}
                    {!r.hasReport ? <span className="ml-2 text-xs text-accent">sin ficha</span> : null}
                  </td>
                  <td className="px-3 py-2">
                    <Checkbox
                      checked={r.sharedInMinistry}
                      onCheckedChange={(checked) => update(r.publisherId, { sharedInMinistry: checked === true })}
                    />
                  </td>
                  <td className="px-3 py-2">
                    <Input
                      className="w-20"
                      type="number"
                      min={0}
                      value={r.bibleStudies}
                      onChange={(e) => update(r.publisherId, { bibleStudies: Number(e.target.value) })}
                    />
                  </td>
                  <td className="px-3 py-2">
                    <Checkbox
                      checked={r.auxiliaryPioneer}
                      onCheckedChange={(checked) =>
                        update(r.publisherId, {
                          auxiliaryPioneer: checked === true,
                          hours: checked === true ? r.hours : null,
                        })
                      }
                    />
                  </td>
                  <td className="px-3 py-2">
                    {hourReporter ? (
                      <Input
                        className="w-24"
                        type="number"
                        min={0}
                        step="0.5"
                        value={r.hours ?? ""}
                        onChange={(e) =>
                          update(r.publisherId, { hours: e.target.value === "" ? null : Number(e.target.value) })
                        }
                      />
                    ) : (
                      <span className="text-muted-foreground">—</span>
                    )}
                  </td>
                  <td className="px-3 py-2">
                    <Checkbox checked={r.late} onCheckedChange={(checked) => update(r.publisherId, { late: checked === true })} />
                  </td>
                  <td className="px-3 py-2">
                    <Input value={r.remarks} onChange={(e) => update(r.publisherId, { remarks: e.target.value })} />
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}

import { useMemo, useState, type FormEvent } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { MonthNav } from "@/components/month-nav";
import { Button } from "@/components/ui/button";
import { Card, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { NativeSelect, NativeSelectOption } from "@/components/ui/native-select";
import { api, type Attendance } from "@/lib/api";

function nowPeriod() {
  const d = new Date();
  return { year: d.getFullYear(), month: d.getMonth() + 1 };
}

function monthBounds(year: number, month: number) {
  const from = `${year}-${String(month).padStart(2, "0")}-01`;
  const last = new Date(year, month, 0).getDate();
  const to = `${year}-${String(month).padStart(2, "0")}-${String(last).padStart(2, "0")}`;
  return { from, to };
}

export function AttendancePage() {
  const qc = useQueryClient();
  const [{ year, month }, setPeriod] = useState(nowPeriod);
  const { from, to } = useMemo(() => monthBounds(year, month), [year, month]);
  const q = useQuery({
    queryKey: ["attendance", from, to],
    queryFn: () => api<Attendance[]>(`/attendance?from=${from}&to=${to}`),
  });

  const [date, setDate] = useState(from);
  const [kind, setKind] = useState<"midweek" | "weekend">("weekend");
  const [inPerson, setInPerson] = useState("0");
  const [online, setOnline] = useState("");

  const save = useMutation({
    mutationFn: () =>
      api("/attendance", {
        method: "PUT",
        body: JSON.stringify({
          date,
          kind,
          inPerson: Number(inPerson) || 0,
          online: online === "" ? null : Number(online),
        }),
      }),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["attendance"] });
      await qc.invalidateQueries({ queryKey: ["dashboard"] });
    },
  });

  const del = useMutation({
    mutationFn: (id: string) => api(`/attendance/${id}`, { method: "DELETE" }),
    onSuccess: async () => {
      await qc.invalidateQueries({ queryKey: ["attendance"] });
      await qc.invalidateQueries({ queryKey: ["dashboard"] });
    },
  });

  function onSubmit(e: FormEvent) {
    e.preventDefault();
    save.mutate();
  }

  return (
    <div className="space-y-6">
      <header className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <h1 className="text-3xl">Asistencia</h1>
          <p className="text-sm text-muted-foreground">Conteo por reunión, no por persona.</p>
        </div>
        <MonthNav year={year} month={month} onChange={(y, m) => setPeriod({ year: y, month: m })} />
      </header>

      <div className="grid gap-6 lg:grid-cols-[1fr_320px]">
        <div className="overflow-x-auto rounded-lg border border-border bg-card">
          <table className="w-full text-left text-sm">
            <thead className="bg-muted">
              <tr>
                <th className="px-3 py-2">Fecha</th>
                <th className="px-3 py-2">Reunión</th>
                <th className="px-3 py-2">Presencial</th>
                <th className="px-3 py-2">En línea</th>
                <th className="px-3 py-2" />
              </tr>
            </thead>
            <tbody>
              {(q.data ?? []).map((a) => (
                <tr key={a.id} className="border-t border-border">
                  <td className="px-3 py-2">{a.date}</td>
                  <td className="px-3 py-2">{a.kind === "midweek" ? "Entre semana" : "Fin de semana"}</td>
                  <td className="px-3 py-2">{a.inPerson}</td>
                  <td className="px-3 py-2">{a.online ?? "—"}</td>
                  <td className="px-3 py-2">
                    <Button variant="ghost" size="sm" onClick={() => del.mutate(a.id)}>
                      Borrar
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <Card>
          <CardTitle>Registrar reunión</CardTitle>
          <form className="mt-4 space-y-3" onSubmit={onSubmit}>
            <div>
              <Label>Fecha</Label>
              <Input className="mt-1" type="date" value={date} onChange={(e) => setDate(e.target.value)} required />
            </div>
            <div>
              <Label>Tipo</Label>
              <NativeSelect className="mt-1 w-full" value={kind} onChange={(e) => setKind(e.target.value as "midweek" | "weekend")}>
                <NativeSelectOption value="midweek">Entre semana</NativeSelectOption>
                <NativeSelectOption value="weekend">Fin de semana</NativeSelectOption>
              </NativeSelect>
            </div>
            <div>
              <Label>Presencial</Label>
              <Input className="mt-1" type="number" min={0} value={inPerson} onChange={(e) => setInPerson(e.target.value)} />
            </div>
            <div>
              <Label>En línea (opcional)</Label>
              <Input className="mt-1" type="number" min={0} value={online} onChange={(e) => setOnline(e.target.value)} />
            </div>
            {save.isError ? <p className="text-sm text-destructive">{save.error.message}</p> : null}
            <Button type="submit" disabled={save.isPending}>
              Guardar
            </Button>
          </form>
        </Card>
      </div>
    </div>
  );
}

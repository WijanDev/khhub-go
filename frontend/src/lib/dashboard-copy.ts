import type { Dashboard } from "@/lib/api";
import { monthTitle } from "@/lib/labels";

export function fmt(n: number | null | undefined, digits = 0) {
  if (n == null) return "—";
  return n.toLocaleString("es", { maximumFractionDigits: digits, minimumFractionDigits: digits });
}

/** Plain-text branch-report figures; same values the dashboard cards show. */
export function formatDashboardTotals(d: Dashboard, year: number, month: number): string {
  return [
    `Totales · ${monthTitle(year, month)}`,
    `Año de servicio ${d.serviceYear}`,
    `Publicadores activos: ${fmt(d.publishers.active)}`,
    `Regulares: ${fmt(d.publishers.regular)} · Irregulares: ${fmt(d.publishers.irregular)}`,
    `Participaron: ${d.reports.shared}/${d.reports.shouldReport}`,
    `Participación: ${fmt(d.reports.participation, 0)} %`,
    `Estudios bíblicos: ${fmt(d.reports.bibleStudies)}`,
    `Horas de precursores: ${fmt(d.reports.pioneerHours, 1)}`,
    `Asistencia entre semana (promedio): ${fmt(d.attendance.midweekAvg, 1)}`,
    `Asistencia fin de semana (promedio): ${fmt(d.attendance.weekendAvg, 1)}`,
  ].join("\n");
}

export const weekdayLabels = [
  "Domingo",
  "Lunes",
  "Martes",
  "Miércoles",
  "Jueves",
  "Viernes",
  "Sábado",
];

export const monthLabels = [
  "enero",
  "febrero",
  "marzo",
  "abril",
  "mayo",
  "junio",
  "julio",
  "agosto",
  "septiembre",
  "octubre",
  "noviembre",
  "diciembre",
];

export const spiritualLabels: Record<string, string> = {
  student: "Estudiante",
  unbaptized_publisher: "Publicador no bautizado",
  publisher: "Publicador",
};

export const activityLabels: Record<string, string> = {
  regular: "Regular",
  irregular: "Irregular",
  inactive: "Inactivo",
};

export const genderLabels: Record<string, string> = {
  male: "Varón",
  female: "Mujer",
};

export function monthTitle(year: number, month: number) {
  return `${monthLabels[month - 1] ?? month} ${year}`;
}

import { Button } from "@/components/ui/button";
import { monthTitle } from "@/lib/labels";

type Props = {
  year: number;
  month: number;
  onChange: (year: number, month: number) => void;
};

export function MonthNav({ year, month, onChange }: Props) {
  function shift(delta: number) {
    const d = new Date(year, month - 1 + delta, 1);
    onChange(d.getFullYear(), d.getMonth() + 1);
  }
  return (
    <div className="flex items-center gap-2">
      <Button variant="outline" size="sm" onClick={() => shift(-1)}>
        ←
      </Button>
      <div className="min-w-40 text-center capitalize">{monthTitle(year, month)}</div>
      <Button variant="outline" size="sm" onClick={() => shift(1)}>
        →
      </Button>
    </div>
  );
}

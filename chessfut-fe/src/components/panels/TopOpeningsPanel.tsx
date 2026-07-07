import type { Card } from "@/types/card.types";

export function TopOpeningsPanel({ card }: { card: Card }) {
  if (!card.top_openings || card.top_openings.length === 0) {
    return null;
  }

  return (
    <div className="w-full max-w-xs rounded-xl border border-white/10 bg-neutral-900 p-5 text-white">
      <h2 className="mb-1 text-xs font-semibold uppercase tracking-widest text-emerald-400">
        Top Openings
      </h2>
      <p className="mb-4 text-xs text-white/40">Last 30 days</p>

      <div className="flex flex-col gap-3">
        {card.top_openings.map((o) => (
          <div key={o.eco + o.name} className="flex items-center justify-between gap-3">
            <div className="min-w-0">
              <p className="truncate text-sm font-semibold">{o.name}</p>
              <p className="text-xs text-white/40">
                {o.eco} · {o.count} game{o.count === 1 ? "" : "s"}
              </p>
            </div>
            <div className="shrink-0 text-right">
              <p className="text-sm font-bold text-emerald-400">{o.win_rate.toFixed(0)}%</p>
              <p className="text-[10px] uppercase tracking-wide text-white/30">Win rate</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
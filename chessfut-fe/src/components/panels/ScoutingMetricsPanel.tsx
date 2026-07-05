import type { Card } from "@/types/card.types";
import { MetricBar } from "./MetricBar";

export function ScoutingMetricsPanel({ card }: { card: Card }) {
  return (
    <div className="w-full max-w-xs rounded-xl border border-white/10 bg-neutral-900 p-5 text-white">
      <h2 className="mb-4 text-xs font-semibold uppercase tracking-widest text-emerald-400">
        Scouting Metrics
      </h2>

      <div className="mb-4 flex justify-between text-sm">
        <span className="text-white/70">Work Rate</span>
        <span className="font-semibold">
          {card.work_rate.attack} / {card.work_rate.defense}
        </span>
      </div>

      <div className="flex flex-col gap-3">
        <MetricBar label="Bullet Rating" value={card.bullet.rating} max={3500} display={String(card.bullet.rating)} />
        <MetricBar label="Blitz Rating" value={card.blitz.rating} max={3500} display={String(card.blitz.rating)} />
        <MetricBar label="Rapid Rating" value={card.rapid.rating} max={3000} display={String(card.rapid.rating)} />
        <MetricBar label="Daily Rating" value={card.daily.rating} max={2600} display={String(card.daily.rating)} />
        {card.tactics_rating ? (
          <MetricBar
            label="Tactics Rating"
            value={card.tactics_rating}
            max={3500}
            display={String(card.tactics_rating)}
          />
        ) : null}
        {card.puzzle_rush_accuracy ? (
          <MetricBar
            label="Puzzle Rush Accuracy"
            value={card.puzzle_rush_accuracy}
            max={100}
            display={`${card.puzzle_rush_accuracy.toFixed(1)}%`}
          />
        ) : null}
        <MetricBar label="Followers" value={card.followers} max={2000000} display={card.followers.toLocaleString()} />
      </div>
    </div>
  );
}
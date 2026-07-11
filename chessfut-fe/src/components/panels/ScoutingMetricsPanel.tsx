import type { Card, Position } from "@/types/card.types";
import { MetricBar } from "./MetricBar";

const RD_UNRELIABLE_THRESHOLD = 50;

const POSITION_EXPLANATIONS: Record<Position, string> = {
  ST: "Forward — assigned for high SHO/PAC and fast, decisive games.",
  CAM: "Attacking Midfielder — assigned for creative, DRI-led playmaking.",
  CM: "Midfielder — assigned for balanced PAS/DRI playmaking.",
  CB: "Defender — assigned for high DEF, longer games and a higher draw rate.",
};

function fideRatingLabel(card: Card): string | null {
  if (card.fide_source === "verified" && card.fide_rating) {
    return String(card.fide_rating);
  }
  if (card.fide_source === "title_default" && card.effective_fide_rating) {
    return `Not linked (assumed ${card.effective_fide_rating})`;
  }
  return null;
}

function hasNoRecordedGames(card: Card): boolean {
  return card.bullet.rating === 0 && card.blitz.rating === 0 && card.rapid.rating === 0;
}

export function ScoutingMetricsPanel({ card }: { card: Card }) {
  const fideLabel = fideRatingLabel(card);
  const hasUnreliableRating =
    (card.bullet.rd > RD_UNRELIABLE_THRESHOLD && card.bullet.rating > 0) ||
    (card.blitz.rd > RD_UNRELIABLE_THRESHOLD && card.blitz.rating > 0) ||
    (card.rapid.rd > RD_UNRELIABLE_THRESHOLD && card.rapid.rating > 0);
  const positionExplanation = card.position ? POSITION_EXPLANATIONS[card.position] : null;

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

      {fideLabel && (
        <div className="mb-4 flex justify-between text-sm">
          <span className="text-white/70">FIDE Rating</span>
          <span className="font-semibold">{fideLabel}</span>
        </div>
      )}

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
        <MetricBar label="Games Played" value={card.games_snapshot} max={100000} display={card.games_snapshot.toLocaleString()} />
      </div>

      {positionExplanation && <p className="mt-4 text-xs text-white/40">{positionExplanation}</p>}
      {fideLabel && card.fide_source === "title_default" && (
        <p className="mt-2 text-xs text-white/40">
          FIDE rating isn&apos;t linked on Chess.com — the {card.title} title&apos;s minimum norm rating was assumed instead.
        </p>
      )}
      {card.title && hasNoRecordedGames(card) && (
        <p className="mt-2 text-xs text-white/40">
          No recorded Chess.com games in any format — this OVR is based on FIDE strength alone, discounted since
          it&apos;s unproven on this platform.
        </p>
      )}
      {hasUnreliableRating && (
        <p className="mt-2 text-xs text-white/40">
          Some ratings are based on limited recent games and may shift as more are played.
        </p>
      )}
      {card.chesscom_url && (
        <a
          href={card.chesscom_url}
          target="_blank"
          rel="noopener noreferrer"
          className="mt-4 block text-center text-xs text-emerald-400 underline hover:text-emerald-300"
        >
          View on Chess.com
        </a>
      )}
    </div>
  );
}
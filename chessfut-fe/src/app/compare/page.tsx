import type { Metadata } from "next";
import { getDetailedCard, getLeaderboard } from "@/lib/api";
import { ComparisonResult } from "@/components/compare/ComparisonResult";
import { CompareSearchForm } from "@/components/compare/CompareSearchForm";
import { TopBar } from "@/components/layout/TopBar";
import { Footer } from "@/components/layout/Footer";
import type { Card } from "@/types/card.types";

export const metadata: Metadata = {
  title: "Compare Players",
  description: "Compare two Chess.com players side by side — OVR, attributes, and stats.",
};

async function loadExamplePair(): Promise<[Card, Card] | null> {
  try {
    const { leaderboard } = await getLeaderboard(50);
    const titled = leaderboard.filter((c) => c.tier === "titled");
    const shuffled = [...titled].sort(() => Math.random() - 0.5);
    if (shuffled.length < 2) return null;
    return [shuffled[0], shuffled[1]];
  } catch {
    return null;
  }
}

export default async function ComparePage({
  searchParams,
}: {
  searchParams: Promise<{ a?: string; b?: string }>;
}) {
  const { a, b } = await searchParams;

  let cardA: Card | null = null;
  let cardB: Card | null = null;
  let errored = false;

  if (a && b) {
    try {
      [cardA, cardB] = await Promise.all([getDetailedCard(a), getDetailedCard(b)]);
    } catch {
      errored = true;
    }
  } else {
    const pair = await loadExamplePair();
    if (pair) [cardA, cardB] = pair;
  }

  return (
    <div className="min-h-screen bg-neutral-900">
      <TopBar showBack />
      <main className="flex flex-col items-center gap-8 p-8 text-white">
        {errored && <p className="text-lg text-white/80">One or both players could not be found.</p>}

        <CompareSearchForm />

        {cardA && cardB && <ComparisonResult cardA={cardA} cardB={cardB} />}
      </main>
      <Footer />
    </div>
  );
}
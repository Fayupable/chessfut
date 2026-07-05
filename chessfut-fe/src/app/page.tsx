import { getFastCard, getStats } from "@/lib/api";
import { SearchForm } from "@/components/home/SearchForm";
import { CardFan } from "@/components/home/CardFan";
import { ChessMascot } from "@/components/home/ChessMascot";
import { TopBar } from "@/components/layout/TopBar";
import { Footer } from "@/components/layout/Footer";
import type { Card } from "@/types/card.types";

const SAMPLE_USERNAMES = ["hikaru", "fabianocaruana", "ckgchess", "magnuscarlsen"];

async function loadSampleCards(): Promise<Card[]> {
  const results = await Promise.allSettled(SAMPLE_USERNAMES.map((u) => getFastCard(u)));
  return results
    .filter((r): r is PromiseFulfilledResult<Card> => r.status === "fulfilled")
    .map((r) => r.value);
}

async function loadTotalCards(): Promise<number | null> {
  try {
    const stats = await getStats();
    return stats.total_cards;
  } catch {
    return null;
  }
}

export default async function HomePage() {
  const [sampleCards, totalCards] = await Promise.all([loadSampleCards(), loadTotalCards()]);

  return (
    <div className="min-h-screen bg-neutral-900">
      <TopBar />
      <main className="mx-auto flex max-w-6xl flex-col items-center justify-center gap-12 px-6 py-16 text-white md:flex-row md:items-center md:justify-between md:gap-12 lg:px-12">
        <div className="flex max-w-md flex-col gap-6">
          <ChessMascot />
          <span className="w-fit rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs tracking-wide text-white/60">
            CHESS.COM × CHESSFUT
          </span>
          <h1 className="text-5xl font-extrabold leading-tight">
            GET SCOUTED<span className="text-emerald-400">.</span>
          </h1>
          <p className="text-white/60">
            Your Chess.com stats, turned into a World-Cup-style player card rated out of 99.
          </p>
          <SearchForm />
          <p className="text-sm text-white/40">
            Example:{" "}
            <a href="/hikaru" className="underline hover:text-white">
              hikaru
            </a>{" "}
            ·{" "}
            <a href="/fabianocaruana" className="underline hover:text-white">
              fabianocaruana
            </a>{" "}
            · or your own
          </p>
          {totalCards !== null && (
            <p className="flex items-center gap-2 text-sm text-white/50">
              <span className="h-2 w-2 rounded-full bg-emerald-400" />
              <span className="font-semibold text-white">{totalCards.toLocaleString()}</span> cards rated
            </p>
          )}
        </div>

        {sampleCards.length > 0 && (
          <div className="hidden sm:block">
            <CardFan cards={sampleCards} />
          </div>
        )}
      </main>
      <Footer />
    </div>
  );
}
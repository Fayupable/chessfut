import Link from "next/link";
import type { Metadata } from "next";
import { getLeaderboard } from "@/lib/api";
import { TopBar } from "@/components/layout/TopBar";
import { Footer } from "@/components/layout/Footer";

export const metadata: Metadata = {
    title: "Leaderboard",
    description: "Top-rated Chess.com players, ranked by OVR on Chessfut.",
};

export default async function LeaderboardPage() {
    const { leaderboard } = await getLeaderboard(50);

    return (
        <div className="min-h-screen bg-neutral-900">
            <TopBar showBack />
            <main className="mx-auto flex max-w-2xl flex-col gap-2 p-8 text-white">
                <h1 className="mb-4 text-2xl font-bold">Leaderboard</h1>

                {leaderboard.length === 0 ? (
                    <p className="text-white/50">No cards yet.</p>
                ) : (
                    leaderboard.map((card, i) => (
                        <Link
                            key={card.username}
                            href={`/${card.username}`}
                            className="flex items-center gap-4 rounded-lg border border-white/10 bg-neutral-800 px-4 py-3 hover:bg-neutral-700"
                        >
                            <span className="w-8 text-right text-white/40">{i + 1}</span>
                            <img src={card.avatar} alt={card.username} className="h-10 w-10 rounded-full object-cover" />
                            <div className="flex-1">
                                <p className="font-semibold">{card.name || card.username}</p>
                                <p className="text-xs text-white/50">
                                    {card.title ? `${card.title} · ` : ""}
                                    {card.position ?? ""}
                                </p>
                            </div>
                            <span className="text-xl font-bold text-emerald-400">{card.ovr}</span>
                        </Link>
                    ))
                )}
            </main>
            <Footer />
        </div>
    );
}
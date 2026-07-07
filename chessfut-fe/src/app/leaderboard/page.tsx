import Link from "next/link";
import type { Metadata } from "next";
import { getLeaderboard } from "@/lib/api";
import { TopBar } from "@/components/layout/TopBar";
import { Footer } from "@/components/layout/Footer";

export const metadata: Metadata = {
    title: "Leaderboard",
    description: "Top-rated Chess.com players, ranked by OVR on Chessfut.",
};

const PAGE_SIZE = 20;

export default async function LeaderboardPage({
    searchParams,
}: {
    searchParams: Promise<{ page?: string }>;
}) {
    const { page: pageParam } = await searchParams;
    const page = Math.max(1, parseInt(pageParam ?? "1", 10) || 1);
    const offset = (page - 1) * PAGE_SIZE;

    const { leaderboard } = await getLeaderboard(PAGE_SIZE, offset);
    const hasNextPage = leaderboard.length === PAGE_SIZE;

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
                            <span className="w-8 text-right text-white/40">{offset + i + 1}</span>
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

                <div className="mt-4 flex items-center justify-between">
                    {page > 1 ? (
                        <Link
                            href={`/leaderboard?page=${page - 1}`}
                            className="rounded-md bg-neutral-800 px-4 py-2 text-sm font-semibold hover:bg-neutral-700"
                        >
                            &larr; Previous
                        </Link>
                    ) : (
                        <span />
                    )}
                    <span className="text-sm text-white/40">Page {page}</span>
                    {hasNextPage ? (
                        <Link
                            href={`/leaderboard?page=${page + 1}`}
                            className="rounded-md bg-neutral-800 px-4 py-2 text-sm font-semibold hover:bg-neutral-700"
                        >
                            Next &rarr;
                        </Link>
                    ) : (
                        <span />
                    )}
                </div>

                <form action="/leaderboard" method="GET" className="mt-2 flex items-center justify-center gap-2">
                    <label htmlFor="page-jump" className="text-sm text-white/40">
                        Go to page
                    </label>
                    <input
                        id="page-jump"
                        name="page"
                        type="number"
                        min={1}
                        defaultValue={page}
                        className="w-20 rounded-md border border-neutral-700 bg-neutral-800 px-2 py-1 text-center text-sm outline-none"
                    />
                    <button
                        type="submit"
                        className="rounded-md bg-emerald-500 px-3 py-1 text-sm font-semibold text-black hover:bg-emerald-400"
                    >
                        Go
                    </button>
                </form>
            </main>
            <Footer />
        </div>
    );
}
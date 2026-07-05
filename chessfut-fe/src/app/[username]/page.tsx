import { getDetailedCard } from "@/lib/api";
import { PlayerCard } from "@/components/card/PlayerCard";
import { CardActions } from "@/components/card/CardActions";
import { ScoutingMetricsPanel } from "@/components/panels/ScoutingMetricsPanel";
import { TopBar } from "@/components/layout/TopBar";
import { Footer } from "@/components/layout/Footer";

export default async function PlayerPage({
    params,
}: {
    params: Promise<{ username: string }>;
}) {
    const { username } = await params;

    try {
        const card = await getDetailedCard(username);
        return (
            <div className="min-h-screen bg-neutral-900">
                <TopBar showBack />
                <main className="flex flex-col items-center justify-center gap-8 p-8 md:flex-row md:items-start md:justify-center">
                    <div className="flex flex-col items-center gap-4">
                        <div className="w-[300px]">
                            <PlayerCard card={card} />
                        </div>
                        <CardActions username={card.username} />
                    </div>
                    <ScoutingMetricsPanel card={card} />
                </main>
                <Footer />
            </div>
        );
    } catch {
        return (
            <div className="min-h-screen bg-neutral-900">
                <TopBar showBack />
                <main className="flex flex-col items-center justify-center gap-4 p-8 text-center text-white">
                    <p className="text-lg text-white/80">
                        Player <span className="font-semibold">&quot;{username}&quot;</span> not found.
                    </p>
                    <p className="text-sm text-white/50">Check the chess.com username and try again.</p>
                    <a
                        href="/"
                        className="mt-2 rounded-md bg-emerald-500 px-5 py-2 font-semibold text-black hover:bg-emerald-400"
                    >
                        Back to search
                    </a>
                </main>
                <Footer />
            </div>
        );
    }
}
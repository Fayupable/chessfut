import type { Metadata } from "next";
import Link from "next/link";
import { getDetailedCard } from "@/lib/api";
import { PlayerCard } from "@/components/card/PlayerCard";
import { CardActions } from "@/components/card/CardActions";
import { ScoutingMetricsPanel } from "@/components/panels/ScoutingMetricsPanel";
import { TopOpeningsPanel } from "@/components/panels/TopOpeningsPanel";
import { TopBar } from "@/components/layout/TopBar";
import { Footer } from "@/components/layout/Footer";
import type { Card } from "@/types/card.types";

export async function generateMetadata({
    params,
}: {
    params: Promise<{ username: string }>;
}): Promise<Metadata> {
    const { username } = await params;
    const card = await getDetailedCard(username).catch(() => null);

    if (!card) {
        return { title: `${username} not found` };
    }

    const displayName = card.name || card.username;
    const title = `${displayName} — ${card.ovr} OVR`;
    const description = `${displayName}'s Chessfut player card: ${card.ovr} OVR, ${card.position ?? ""} — PAC ${card.attributes.pac} SHO ${card.attributes.sho} PAS ${card.attributes.pas} DRI ${card.attributes.dri} DEF ${card.attributes.def} PHY ${card.attributes.phy}.`;

    return {
        title,
        description,
        openGraph: { title, description },
        twitter: { title, description },
    };
}

export default async function PlayerPage({
    params,
}: {
    params: Promise<{ username: string }>;
}) {
    const { username } = await params;
    const card: Card | null = await getDetailedCard(username).catch(() => null);

    if (!card) {
        return (
            <div className="min-h-screen bg-neutral-900">
                <TopBar showBack />
                <main className="flex flex-col items-center justify-center gap-4 p-8 text-center text-white">
                    <p className="text-lg text-white/80">
                        Player <span className="font-semibold">&quot;{username}&quot;</span> not found.
                    </p>
                    <p className="text-sm text-white/50">Check the chess.com username and try again.</p>
                    <Link
                        href="/"
                        className="mt-2 rounded-md bg-emerald-500 px-5 py-2 font-semibold text-black hover:bg-emerald-400"
                    >
                        Back to search
                    </Link>
                </main>
                <Footer />
            </div>
        );
    }

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
                <div className="flex flex-col gap-6">
                    <ScoutingMetricsPanel card={card} />
                    <TopOpeningsPanel card={card} />
                </div>
            </main>
            <Footer />
        </div>
    );
}
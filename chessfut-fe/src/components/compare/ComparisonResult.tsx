import type { Card } from "@/types/card.types";
import { PlayerCard } from "@/components/card/PlayerCard";
import { ComparisonTable } from "./ComparisonTable";

export function ComparisonResult({ cardA, cardB }: { cardA: Card; cardB: Card }) {
    return (
        <>
            <h1 className="text-2xl font-bold">
                {cardA.username} vs {cardB.username}
            </h1>
            <div className="flex flex-col items-center gap-8 md:flex-row md:items-start">
                <div className="w-[280px]">
                    <PlayerCard card={cardA} />
                </div>
                <ComparisonTable cardA={cardA} cardB={cardB} />
                <div className="w-[280px]">
                    <PlayerCard card={cardB} />
                </div>
            </div>
        </>
    );
}
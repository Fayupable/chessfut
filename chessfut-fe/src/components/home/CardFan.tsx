"use client";

import { useState } from "react";
import Link from "next/link";
import type { Card } from "@/types/card.types";
import { PlayerCard } from "@/components/card/PlayerCard";

const OFFSETS = [
    { rotate: -8, translateX: -130 },
    { rotate: -3, translateX: -45 },
    { rotate: 3, translateX: 45 },
    { rotate: 8, translateX: 130 },
];

export function CardFan({ cards }: { cards: Card[] }) {
    const [activeIndex, setActiveIndex] = useState<number | null>(null);

    return (
        <div className="relative flex h-[380px] w-[420px] items-center justify-center">
            {cards.slice(0, 4).map((card, i) => {
                const offset = OFFSETS[i] ?? OFFSETS[0];
                const isActive = activeIndex === i;

                return (
                    <Link
                        key={card.username}
                        href={`/${card.username}`}
                        onMouseEnter={() => setActiveIndex(i)}
                        onFocus={() => setActiveIndex(i)}
                        className="absolute w-[180px] transition-transform duration-200 hover:-translate-y-2"
                        style={{
                            transform: `translateX(${offset.translateX}px) rotate(${offset.rotate}deg)`,
                            zIndex: isActive ? 50 : i,
                        }}
                    >
                        <PlayerCard card={card} />
                    </Link>
                );
            })}
        </div>
    );
}
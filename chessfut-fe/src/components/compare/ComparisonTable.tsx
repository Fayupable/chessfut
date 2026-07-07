import type { Card } from "@/types/card.types";

function Row({ label, a, b }: { label: string; a: number; b: number }) {
    const aWins = a > b;
    const bWins = b > a;

    return (
        <div className="grid grid-cols-3 items-center gap-4 py-1 text-sm">
            <span className={`text-right ${aWins ? "font-bold text-emerald-400" : "text-white/70"}`}>{a}</span>
            <span className="text-center text-white/50">{label}</span>
            <span className={`text-left ${bWins ? "font-bold text-emerald-400" : "text-white/70"}`}>{b}</span>
        </div>
    );
}

function TextRow({ label, a, b }: { label: string; a: string; b: string }) {
    return (
        <div className="grid grid-cols-3 items-center gap-4 py-1 text-sm">
            <span className="text-right text-white/70">{a || "—"}</span>
            <span className="text-center text-white/50">{label}</span>
            <span className="text-left text-white/70">{b || "—"}</span>
        </div>
    );
}

function fideDisplay(card: Card): string {
    if (card.fide_source === "verified" && card.fide_rating) {
        return String(card.fide_rating);
    }
    if (card.fide_source === "title_default" && card.effective_fide_rating) {
        return `${card.effective_fide_rating}*`;
    }
    return "—";
}

export function ComparisonTable({ cardA, cardB }: { cardA: Card; cardB: Card }) {
    return (
        <div className="flex w-full max-w-xs flex-col gap-1 rounded-xl border border-white/10 bg-neutral-800 p-4">
            <TextRow label="Title" a={cardA.title} b={cardB.title} />
            <TextRow label="FIDE" a={fideDisplay(cardA)} b={fideDisplay(cardB)} />
            <div className="my-2 h-px bg-white/10" />
            <Row label="OVR" a={cardA.ovr} b={cardB.ovr} />
            <Row label="PAC" a={cardA.attributes.pac} b={cardB.attributes.pac} />
            <Row label="SHO" a={cardA.attributes.sho} b={cardB.attributes.sho} />
            <Row label="PAS" a={cardA.attributes.pas} b={cardB.attributes.pas} />
            <Row label="DRI" a={cardA.attributes.dri} b={cardB.attributes.dri} />
            <Row label="DEF" a={cardA.attributes.def} b={cardB.attributes.def} />
            <Row label="PHY" a={cardA.attributes.phy} b={cardB.attributes.phy} />
            <div className="my-2 h-px bg-white/10" />
            <Row label="Bullet" a={cardA.bullet.rating} b={cardB.bullet.rating} />
            <Row label="Blitz" a={cardA.blitz.rating} b={cardB.blitz.rating} />
            <Row label="Rapid" a={cardA.rapid.rating} b={cardB.rapid.rating} />
            <Row label="Daily" a={cardA.daily.rating} b={cardB.daily.rating} />
            {(cardA.fide_source === "title_default" || cardB.fide_source === "title_default") && (
                <p className="mt-2 text-center text-[10px] text-white/30">* FIDE not linked — title minimum assumed</p>
            )}
        </div>
    );
}
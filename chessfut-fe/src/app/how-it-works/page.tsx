import type { Metadata } from "next";
import { TopBar } from "@/components/layout/TopBar";
import { Footer } from "@/components/layout/Footer";

export const metadata: Metadata = {
    title: "How OVR is calculated",
    description: "How Chessfut turns your Chess.com stats into a player card rating.",
};

export default function HowItWorksPage() {
    return (
        <div className="min-h-screen bg-neutral-900">
            <TopBar showBack />
            <main className="mx-auto flex max-w-2xl flex-col gap-6 p-8 text-white">
                <h1 className="text-2xl font-bold">How OVR is calculated</h1>

                <section className="flex flex-col gap-2">
                    <h2 className="text-lg font-semibold text-emerald-400">The anchor: FIDE first, then Chess.com</h2>
                    <p className="text-sm text-white/70">
                        Every card starts from a single &quot;how strong is this player, really&quot; number. If a
                        real FIDE rating is linked to the account, it dominates — FIDE is independently verified and
                        far more tightly banded than online ratings. If the account has a title (GM, IM, FM, CM) but
                        no linked FIDE rating, we assume the title&apos;s minimum norm rating instead of treating the
                        player as unrated. Without a title or a FIDE rating, Chess.com&apos;s own ratings (bullet,
                        blitz, rapid) drive the score directly — a genuinely strong untitled player can still land a
                        high OVR.
                    </p>
                </section>

                <section className="flex flex-col gap-2">
                    <h2 className="text-lg font-semibold text-emerald-400">The six attributes</h2>
                    <p className="text-sm text-white/70">
                        PAC, SHO, PAS, DRI, DEF and PHY are each derived from a different signal — speed ratings,
                        blitz win rate, rapid rating, tactics/puzzle performance and opening variety, non-loss rate,
                        and total games played — then shaped around the anchor score so a player&apos;s own relative
                        strengths and weaknesses show through.
                    </p>
                </section>

                <section className="flex flex-col gap-2">
                    <h2 className="text-lg font-semibold text-emerald-400">Position</h2>
                    <p className="text-sm text-white/70">
                        Your position (ST, CAM, CM, CB) is assigned from your attribute profile and playing style —
                        aggressive, decisive, short games push toward attack; a high draw rate and longer games push
                        toward defense. The final OVR is a weighted blend of your six attributes, with the weights
                        depending on your assigned position — similar to how a real FIFA card&apos;s rating emphasizes
                        different stats for a striker versus a defender.
                    </p>
                </section>

                <section className="flex flex-col gap-2">
                    <h2 className="text-lg font-semibold text-emerald-400">A note on data reliability</h2>
                    <p className="text-sm text-white/70">
                        Some ratings are based on very few recent games and can shift quickly — we flag this on a
                        player&apos;s card when it applies. An account with no recorded games in any format is
                        floored well below any account with real game history, regardless of what else is on file.
                    </p>
                </section>

                <p className="text-xs text-white/40">
                    This is a fun, unofficial project — it has no affiliation with Chess.com or FIDE, and OVR is not
                    an official rating of any kind.
                </p>
            </main>
            <Footer />
        </div>
    );
}
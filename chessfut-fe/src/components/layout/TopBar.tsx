import Link from "next/link";

export function TopBar({ showBack = false }: { showBack?: boolean }) {
    return (
        <header className="flex w-full items-center justify-between px-8 py-4 text-white">
            <div>
                {showBack ? (
                    <Link href="/" className="text-sm text-white/60 hover:text-white">
                        ← Back
                    </Link>
                ) : (
                    <Link href="/" className="text-sm font-bold tracking-wide">
                        CHESSFUT
                    </Link>
                )}
            </div>

            <div className="flex items-center gap-4 text-sm">
                <Link href="/leaderboard" className="text-white/60 hover:text-white">
                    Leaderboard
                </Link>
                <Link href="/compare" className="text-white/60 hover:text-white">
                    Compare
                </Link>
                <a
                    href="https://github.com/fayupable/chessfut"
                    target="_blank"
                    rel="noreferrer"
                    className="rounded-full border border-white/10 bg-white/5 px-3 py-1 text-white/80 hover:bg-white/10"
                >
                    Contribute on GitHub
                </a>
            </div>
        </header>
    );
}
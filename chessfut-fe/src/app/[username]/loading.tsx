import { TopBar } from "@/components/layout/TopBar";

export default function Loading() {
    return (
        <div className="min-h-screen bg-neutral-900">
            <TopBar showBack />
            <main className="flex flex-col items-center justify-center gap-4 p-8 text-center text-white">
                <div className="h-10 w-10 animate-spin rounded-full border-4 border-emerald-500 border-t-transparent" />
                <p className="text-lg text-white/80">Building your card...</p>
                <p className="text-sm text-white/50">First-time lookups can take a few seconds while we pull fresh stats from chess.com.</p>
            </main>
        </div>
    );
}
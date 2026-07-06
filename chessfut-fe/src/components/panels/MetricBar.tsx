export function MetricBar({ label, value, max, display }: { label: string; value: number; max: number; display: string }) {
    const percent = Math.min(100, Math.round((value / max) * 100));

    return (
        <div className="flex flex-col gap-1">
            <div className="flex justify-between text-sm text-white/80">
                <span>{label}</span>
                <span className="font-semibold text-white">{display}</span>
            </div>
            <div className="h-1.5 w-full rounded-full bg-white/10">
                <div className="h-1.5 rounded-full bg-emerald-400" style={{ width: `${percent}%` }} />
            </div>
        </div>
    );
}
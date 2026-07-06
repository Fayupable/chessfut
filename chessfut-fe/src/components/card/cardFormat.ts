export const pad2 = (n: number) => String(Math.round(n)).padStart(2, "0");

export function formatDisplayName(fullName: string): string {
    const trimmed = fullName.trim();
    const short = trimmed.length <= 9 ? trimmed : trimmed.split(" ").slice(-1)[0];
    return short.toUpperCase();
}
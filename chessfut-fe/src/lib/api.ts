import type { Card, LeaderboardResponse } from "@/types/card.types";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

async function fetchJSON<T>(path: string): Promise<T> {
    const res = await fetch(`${API_URL}${path}`, { cache: "no-store" });
    if (!res.ok) {
        throw new Error(`Request failed: ${res.status}`);
    }
    return res.json();
}

export function getFastCard(username: string): Promise<Card> {
    return fetchJSON<Card>(`/player/${username}`);
}

export function getDetailedCard(username: string): Promise<Card> {
    return fetchJSON<Card>(`/player/${username}/detailed`);
}

export function getLeaderboard(limit = 20, offset = 0): Promise<LeaderboardResponse> {
    return fetchJSON<LeaderboardResponse>(`/leaderboard?limit=${limit}&offset=${offset}`);
}

export function getStats(): Promise<{ total_cards: number }> {
    return fetchJSON<{ total_cards: number }>("/stats");
}
export function searchPlayers(query: string, limit = 10, offset = 0): Promise<{ players: string[] }> {
    return fetchJSON<{ players: string[] }>(`/search?q=${encodeURIComponent(query)}&limit=${limit}&offset=${offset}`);
}
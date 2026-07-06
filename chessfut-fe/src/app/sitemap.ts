import type { MetadataRoute } from "next";
import { getLeaderboard } from "@/lib/api";

const BASE_URL = "https://chessfut.com";

export default async function sitemap(): Promise<MetadataRoute.Sitemap> {
  const staticRoutes: MetadataRoute.Sitemap = [
    { url: BASE_URL, changeFrequency: "daily", priority: 1 },
    { url: `${BASE_URL}/leaderboard`, changeFrequency: "hourly", priority: 0.8 },
    { url: `${BASE_URL}/compare`, changeFrequency: "daily", priority: 0.6 },
  ];

  const playerRoutes: MetadataRoute.Sitemap = await getLeaderboard(100)
    .then(({ leaderboard }) =>
      leaderboard.map((card) => ({
        url: `${BASE_URL}/${card.username}`,
        changeFrequency: "daily" as const,
        priority: 0.5,
      }))
    )
    .catch(() => []);

  return [...staticRoutes, ...playerRoutes];
}
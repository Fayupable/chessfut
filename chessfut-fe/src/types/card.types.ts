export type CardTier = "standard" | "titled";
export type CardType = "fast" | "detailed";
export type PlayStyle = "aggressive" | "positional" | "defensive" | "balanced";
export type Position = "ST" | "CM" | "CB" | "CAM";
export type WorkRateLevel = "High" | "Med" | "Low";
export type FideSource = "verified" | "title_default" | "none";

export interface TimeControlStats {
  rating: number;
  highest: number;
  wins: number;
  losses: number;
  draws: number;
  win_rate: number;
  rd: number;
}

export interface OpeningStat {
  eco: string;
  name: string;
  count: number;
  win_rate: number;
}

export interface Attributes {
  pac: number;
  sho: number;
  pas: number;
  dri: number;
  def: number;
  phy: number;
}

export interface WorkRate {
  attack: WorkRateLevel;
  defense: WorkRateLevel;
}

export interface Card {
  username: string;
  name: string;
  title: string;
  avatar: string;
  followers: number;
  country_code: string;
  joined_date: string;
  chesscom_url?: string;
  last_online?: string;
  card_type: CardType;
  tier: CardTier;
  ovr: number;
  play_style?: PlayStyle;
  position?: Position;
  tactics_rating?: number;
  puzzle_rush_accuracy?: number;
  fide_rating?: number;
  effective_fide_rating?: number;
  fide_source?: FideSource;
  games_snapshot: number;
  attributes: Attributes;
  work_rate: WorkRate;
  badges: string[];
  bullet: TimeControlStats;
  blitz: TimeControlStats;
  rapid: TimeControlStats;
  daily: TimeControlStats;
  top_openings?: OpeningStat[];
  last_updated: string;
}

export interface LeaderboardResponse {
  leaderboard: Card[];
}
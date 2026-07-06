import type { CSSProperties } from "react";
import type { Card } from "@/types/card.types";
import { resolveCardFrame } from "./tierStyles";
import { CardFlag } from "./CardFlag";
import { pad2, formatDisplayName } from "./cardFormat";

const STAT_CELLS: { key: keyof Card["attributes"]; label: string; vx: number; lx: number; vy: number; ly: number }[] = [
  { key: "pac", label: "PAC", vx: 21.3, lx: 32.41, vy: 64.63, ly: 65.24 },
  { key: "dri", label: "DRI", vx: 56.48, lx: 67.59, vy: 64.63, ly: 65.24 },
  { key: "sho", label: "SHO", vx: 21.3, lx: 32.41, vy: 72.2, ly: 72.8 },
  { key: "def", label: "DEF", vx: 56.48, lx: 67.59, vy: 72.2, ly: 72.8 },
  { key: "pas", label: "PAS", vx: 21.3, lx: 32.41, vy: 79.76, ly: 80.37 },
  { key: "phy", label: "PHY", vx: 56.48, lx: 67.59, vy: 79.76, ly: 80.37 },
];

const H_LINES: [number, number, number][] = [
  [19.44, 31.1, 10.19],
  [19.44, 40.85, 10.19],
  [16.67, 64.02, 66.67],
  [44.44, 89.63, 11.11],
];

const at = (left: number, top: number): CSSProperties => ({
  position: "absolute",
  left: `${left}%`,
  top: `${top}%`,
});

export function PlayerCard({ card }: { card: Card }) {
  const displayName = formatDisplayName(card.name || card.username);
  const frame = resolveCardFrame(card.ovr);

  const wrap: CSSProperties = {
    containerType: "inline-size",
    position: "relative",
    width: "100%",
    aspectRatio: "540 / 820",
    filter: "drop-shadow(0 7cqw 10cqw rgba(0,0,0,.5))",
    userSelect: "none",
    WebkitUserSelect: "none",
    WebkitTouchCallout: "none",
  };

  return (
    <div style={wrap}>
      {/* tier background art carries the card silhouette */}
      <img
        src={frame}
        alt=""
        aria-hidden
        style={{ position: "absolute", inset: 0, width: "100%", height: "100%", objectFit: "fill" }}
      />

      {/* avatar clipped to the card silhouette so it can never spill past the frame */}
      <div
        style={{
          position: "absolute",
          inset: 0,
          WebkitMaskImage: `url("${frame}")`,
          maskImage: `url("${frame}")`,
          WebkitMaskSize: "100% 100%",
          maskSize: "100% 100%",
        }}
      >
        <div
          style={{
            position: "absolute",
            left: "30cqw",
            top: "13cqw",
            width: "50cqw",
            height: "55cqw",
            WebkitMaskImage:
              "radial-gradient(ellipse 66% 88% at 52% 40%, #000 56%, transparent 80%), linear-gradient(180deg, transparent 1%, #000 22%)",
            maskImage:
              "radial-gradient(ellipse 66% 88% at 52% 40%, #000 56%, transparent 80%), linear-gradient(180deg, transparent 1%, #000 22%)",
            WebkitMaskComposite: "source-in",
            maskComposite: "intersect",
          }}
        >
          <img
            src={
              card.avatar ||
              "data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='320' height='320'><circle cx='160' cy='130' r='60' fill='%23999'/><rect x='60' y='210' width='200' height='140' rx='70' fill='%23999'/></svg>"
            }
            alt={card.username}
            style={{ width: "100%", height: "100%", objectFit: "cover", objectPosition: "center 20%" }}
          />
        </div>
      </div>

      {/* separator lines */}
      {H_LINES.map(([l, top, w], i) => (
        <div key={i} style={{ ...at(l, top), width: `${w}%`, height: "0.3cqw", background: "#000", opacity: 0.5 }} />
      ))}
      <div style={{ ...at(50, 66.46), width: "0.3cqw", height: "20.12%", background: "#000", opacity: 0.5 }} />

      {/* OVR */}
      <div style={{ ...at(16.3, 9.76), fontSize: "22.2cqw", fontWeight: 500, lineHeight: 1, color: "#000" }}>
        {pad2(card.ovr)}
      </div>

      {/* position */}
      <div
        style={{
          ...at(25, 23.78),
          transform: "translateX(-50%)",
          fontSize: "9.3cqw",
          fontWeight: 500,
          letterSpacing: ".02em",
          color: "#000",
        }}
      >
        {card.position ?? "—"}
      </div>

      <CardFlag countryCode={card.country_code} />

      {/* name */}
      <div
        style={{
          ...at(50, 53.66),
          transform: "translateX(-50%)",
          fontSize: "13cqw",
          fontWeight: 700,
          whiteSpace: "nowrap",
          color: "#000",
        }}
      >
        {displayName}
      </div>

      {/* six attributes */}
      {STAT_CELLS.map((c) => (
        <div key={c.key}>
          <span style={{ ...at(c.vx, c.vy), fontSize: "10.2cqw", fontWeight: 700, color: "#000" }}>
            {pad2(card.attributes[c.key])}
          </span>
          <span
            style={{
              ...at(c.lx, c.ly),
              fontSize: "9.3cqw",
              fontWeight: 500,
              letterSpacing: ".02em",
              marginLeft: "1cqw",
              color: "#000",
            }}
          >
            {c.label}
          </span>
        </div>
      ))}
    </div>
  );
}
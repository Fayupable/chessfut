export function resolveCardFrame(ovr: number): string {
  if (ovr >= 90) return "/cards/legend.png";
  if (ovr >= 80) return "/cards/gold.png";
  if (ovr >= 65) return "/cards/silver.png";
  return "/cards/bronze.png";
}
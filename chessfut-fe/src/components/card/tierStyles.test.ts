import { describe, it, expect } from "vitest";
import { resolveCardFrame } from "./tierStyles";

describe("resolveCardFrame", () => {
    it("returns bronze for low OVR", () => {
        expect(resolveCardFrame(50)).toBe("/cards/bronze.png");
    });

    it("returns silver for mid OVR", () => {
        expect(resolveCardFrame(70)).toBe("/cards/silver.png");
    });

    it("returns gold for high OVR", () => {
        expect(resolveCardFrame(85)).toBe("/cards/gold.png");
    });

    it("returns legend for elite OVR", () => {
        expect(resolveCardFrame(92)).toBe("/cards/legend.png");
    });

    it("uses boundary values correctly", () => {
        expect(resolveCardFrame(65)).toBe("/cards/silver.png");
        expect(resolveCardFrame(80)).toBe("/cards/gold.png");
        expect(resolveCardFrame(90)).toBe("/cards/legend.png");
    });
});
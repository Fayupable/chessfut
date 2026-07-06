import { describe, it, expect } from "vitest";
import { pad2, formatDisplayName } from "./cardFormat";

describe("pad2", () => {
  it("pads single-digit numbers with a leading zero", () => {
    expect(pad2(5)).toBe("05");
  });

  it("leaves two-digit numbers unchanged", () => {
    expect(pad2(88)).toBe("88");
  });

  it("rounds decimal values before padding", () => {
    expect(pad2(7.6)).toBe("08");
  });
});

describe("formatDisplayName", () => {
  it("keeps short full names as-is, uppercased", () => {
    expect(formatDisplayName("Magnus")).toBe("MAGNUS");
  });

  it("keeps names exactly 9 characters long in full", () => {
    expect(formatDisplayName("Hikaru N.")).toBe("HIKARU N.");
  });

  it("falls back to the last word for long names", () => {
    expect(formatDisplayName("Magnus Carlsen")).toBe("CARLSEN");
  });

  it("trims surrounding whitespace before measuring length", () => {
    expect(formatDisplayName("  Magnus  ")).toBe("MAGNUS");
  });
});
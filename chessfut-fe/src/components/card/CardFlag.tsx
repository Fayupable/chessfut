"use client";

export function CardFlag({ countryCode }: { countryCode: string }) {
    if (!countryCode) return null;

    return (
        <img
            alt={countryCode}
            src={`/badges/flags/${countryCode.toLowerCase()}.png`}
            onError={(e) => {
                e.currentTarget.style.visibility = "hidden";
            }}
            className="absolute object-contain"
            style={{ left: "17.59%", top: "33.17%", width: "14.81%", height: "5.73%" }}
        />
    );
}
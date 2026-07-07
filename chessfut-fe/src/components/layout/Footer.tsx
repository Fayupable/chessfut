import Link from "next/link";

export function Footer() {
    return (
        <footer className="flex w-full items-center justify-center gap-4 py-8 text-sm text-white/40">
            <Link href="/how-it-works" className="hover:text-white">
                How OVR is calculated
            </Link>
            <span className="text-white/20">·</span>
            <span>Built by</span>
            <a
                href="https://github.com/Fayupable"
                target="_blank"
                rel="noreferrer"
                className="font-semibold text-white/70 hover:text-white"
            >
                @Fayupable
            </a>
        </footer>
    );
}
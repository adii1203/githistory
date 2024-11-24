import type { Metadata } from "next";
import "./globals.css";
import { Button } from "@/components/ui/button";
import Image from "next/image";
import Link from "next/link";
import localFont from "next/font/local";
import { Toaster } from "@/components/ui/sonner";

export const metadata: Metadata = {
  title: "gitgraph",
  description: "Generated beautiful github stars graph",
};

const Space = localFont({
  src: "./fonts/SpaceGrotesk.ttf",
});

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={`${Space.className} antialiased`}>
        <div className="flex min-h-screen flex-col">
          <main className="flex-1 flex">
            <div className="container max-w-2xl mx-auto px-4  text-center">
              <section>{children}</section>
            </div>
            <Toaster />
          </main>
          <footer className="border-t">
            <div className="flex items-center justify-center py-2">
              <div className="flex items-center justify-center gap-1">
                <p className="flex text-center text-sm leading-loose text-muted-foreground">
                  Created with ❤️ by{" "}
                </p>
                <ProfileButton />
              </div>
            </div>
          </footer>
        </div>
      </body>
    </html>
  );
}

function ProfileButton() {
  return (
    <Button
      variant={"ghost"}
      className="rounded-full h-7 py-0 ps-0 items-center justify-center">
      <Link target="_blank" href={"https://github.com/adii1203"}>
        <div className="flex items-center">
          <div className="w-8 aspect-square h-full p-1.5">
            <Image
              className="h-auto w-full rounded-full"
              src="https://avatars.githubusercontent.com/u/114096753?v=4"
              alt="Profile image"
              width={20}
              height={20}
              aria-hidden="true"
            />
          </div>
          <p>adii</p>
        </div>
      </Link>
    </Button>
  );
}

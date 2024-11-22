"use client";

import { Input } from "@/components/ui/input";
import Image from "next/image";
import * as htmlToImage from "html-to-image";

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
  ResponsiveContainer,
} from "recharts";
import { Button } from "@/components/ui/button";
import { useRef } from "react";
import { Download } from "lucide-react";

const data = [
  {
    date: "Aug 2022",
    stars: 0,
  },
  {
    date: "Sep 2022",
    stars: 2460,
  },
  {
    date: "Nov 2022",
    stars: 3720,
  },
  {
    date: "Jan 2023",
    stars: 4980,
  },
  {
    date: "May 2023",
    stars: 7500,
  },
  {
    date: "Jun 2023",
    stars: 8760,
  },
  {
    date: "Aug 2023",
    stars: 10050,
  },
  {
    date: "Sep 2023",
    stars: 11310,
  },
  {
    date: "Nov 2023",
    stars: 12570,
  },
  {
    date: "Jan 2024",
    stars: 13830,
  },
  {
    date: "Mar 2024",
    stars: 15090,
  },
  {
    date: "Jun 2024",
    stars: 16350,
  },
  {
    date: "Sep 2024",
    stars: 17610,
  },
  {
    date: "Nov 2024",
    stars: 18948,
  },
  {
    date: "Dec 2024",
    stars: 19000,
  },
];

export default function Home() {
  const divRef = useRef<HTMLDivElement>(null);
  const exportAsImage = async (el: HTMLDivElement | null, fileName: string) => {
    if (!el) return;
    const blob = await htmlToImage.toPng(el);
    downloadImage({ blob, fileName });
  };

  const downloadImage = ({
    blob,
    fileName,
  }: {
    blob: string;
    fileName: string;
  }) => {
    const a = document.createElement("a");
    a.download = `${fileName}_gitgraph.png`;
    a.href = blob;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    a.remove();
  };

  return (
    <div>
      <div>
        <div className="space-y-2 w-[20rem] sm:w-[28rem] mx-auto">
          <div className="flex rounded-lg shadow-sm shadow-black/5">
            <Input
              id="input-21"
              className="-me-px flex-1 rounded-e-none shadow-none focus-visible:z-10"
              placeholder="dubinc/dub"
              type="search"
            />
            <button className="inline-flex items-center rounded-e-lg border border-input bg-background px-3 text-sm font-medium text-foreground outline-offset-2 transition-colors hover:bg-accent hover:text-foreground focus:z-10 focus-visible:outline focus-visible:outline-2 focus-visible:outline-ring/70 disabled:cursor-not-allowed disabled:opacity-50">
              Search
            </button>
          </div>
        </div>
      </div>
      <div className="relative">
        <div className="p-4 relative" ref={divRef}>
          <div className="bg-white space-y-6 mx-auto mt-10 p-6 rounded-lg shadow-2xl w-[100%] sm:w-[600px]">
            <div className="flex items-stretch space-y-0 border-b p-0">
              <div className="flex flex-1 flex-col justify-center gap-1 px-2 py-4">
                <div className="flex items-center gap-2">
                  <div className="w-8 h-8 rounded-full overflow-hidden">
                    <Image
                      src="https://avatars.githubusercontent.com/u/153106492?s=200&v=4"
                      alt="repo image"
                      width={40}
                      height={40}
                    />
                  </div>
                  <h3 className="text-2xl font-semibold leading-none tracking-tight">
                    dubinc/dub
                  </h3>
                </div>
              </div>
              <div className="flex border-l">
                <div className="flex flex-1 flex-col justify-center gap-1 px-6 py-4 text-left sm:px-8 sm:py-4">
                  <span className="text-xs">Total Stars</span>
                  <span className="text-lg font-bold">19,000</span>
                </div>
              </div>
            </div>
            <div className="">
              <div
                className="flex
      w-full
      
        justify-center
        text-xs
      [&_.recharts-cartesian-axis-tick_text]:fill-gray-400
      [&_.recharts-cartesian-grid_line[stroke='#ccc']]:stroke-black/30 
      [&_.recharts-curve.recharts-tooltip-cursor]:stroke-black/20 
        [&_.recharts-dot[stroke='#fff']]:stroke-transparent 
        [&_.recharts-layer]:outline-none 
        [&_.recharts-radial-bar-background-sector]:fill-muted 
        [&_.recharts-rectangle.recharts-tooltip-cursor]:fill-muted 
        [&_.recharts-reference-line_[stroke='#ccc']]:stroke-border 
        [&_.recharts-sector[stroke='#fff']]:stroke-transparent 
        [&_.recharts-sector]:outline-none 
        [&_.recharts-surface]:outline-none">
                <ResponsiveContainer width="100%" height={300}>
                  <LineChart data={data}>
                    <XAxis
                      minTickGap={20}
                      padding={{ left: 10, right: 40 }}
                      axisLine={false}
                      tickLine={false}
                      dataKey="date"
                    />
                    <YAxis
                      axisLine={false}
                      tickLine={false}
                      dataKey={"stars"}
                    />
                    <Tooltip />
                    <Line
                      type={"natural"}
                      legendType="line"
                      dot={false}
                      strokeWidth={2}
                      strokeLinecap="round"
                      dataKey={"stars"}
                      stroke="#ffc073"
                    />
                    <CartesianGrid
                      vertical={false}
                      stroke="#ccc"
                      opacity={0.3}
                    />
                  </LineChart>
                </ResponsiveContainer>
              </div>
            </div>
          </div>
        </div>
        <div>
          <Button
            variant={"ghost"}
            size={"icon"}
            className="rounded-full absolute right-10 top-16"
            onClick={() => exportAsImage(divRef.current, "dub")}>
            <Download size={16} />
          </Button>
        </div>
      </div>
    </div>
  );
}

"use client";

import Image from "next/image";
import React, { useRef } from "react";
import {
  CartesianGrid,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { Button } from "./ui/button";
import { Download } from "lucide-react";
import * as htmlToImage from "html-to-image";

type ChartData = {
  name: string;
  logo_url: string;
  total_stars: number;
  data: { date: string; stars: number }[];
};

const Graph = ({ data }: { data: ChartData }) => {
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
    <div className="relative">
      <div className="p-4 relative" ref={divRef}>
        <div className="bg-white space-y-6 mx-auto mt-10 p-6 rounded-lg shadow-2xl w-[100%] sm:w-[600px]">
          <div className="flex items-stretch space-y-0 border-b p-0">
            <div className="flex flex-1 flex-col justify-center gap-1 px-2 py-4">
              <div className="flex items-center gap-2">
                <div className="w-8 h-8 rounded-full overflow-hidden">
                  <Image
                    src={data.logo_url}
                    alt="repo image"
                    width={40}
                    height={40}
                  />
                </div>
                <h3 className="text-2xl font-semibold leading-none tracking-tight">
                  {data?.name}
                </h3>
              </div>
            </div>
            <div className="flex border-l">
              <div className="flex flex-1 flex-col justify-center gap-1 px-6 py-4 text-left sm:px-8 sm:py-4">
                <span className="text-xs">Total Stars</span>
                <span className="text-lg font-bold">{data?.total_stars}</span>
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
                <LineChart data={data.data}>
                  <XAxis
                    minTickGap={20}
                    padding={{ left: 10, right: 40 }}
                    axisLine={false}
                    tickLine={false}
                    dataKey="date"
                  />
                  <YAxis axisLine={false} tickLine={false} dataKey={"stars"} />
                  <Tooltip />
                  <Line
                    type={"natural"}
                    legendType="line"
                    strokeWidth={3}
                    strokeLinecap="round"
                    dataKey={"stars"}
                    stroke="#47c98f"
                  />
                  <CartesianGrid vertical={false} stroke="#ccc" opacity={0.3} />
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
  );
};

export default Graph;

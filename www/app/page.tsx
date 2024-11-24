"use client";

import { useEffect, useState } from "react";
import axios from "axios";
import { toast } from "sonner";
import {
  AddTokenDialog,
  RemoveTokenButton,
} from "@/components/add-remove-token";
import Graph from "@/components/graph";
import Search from "@/components/search";
import SkeletonCard from "@/components/card-loader";

type ChartData = {
  name: string;
  logo_url: string;
  total_stars: number;
  data: { date: string; stars: number }[];
};

export default function Home() {
  const [token, setToken] = useState<string>();

  useEffect(() => {
    const token = localStorage.getItem("github_token");
    if (token) {
      setToken(token);
    }
  }, []);

  const [data, setData] = useState<ChartData>({
    data: [
      {
        date: "Dec 2014",
        stars: 0,
      },
      {
        date: "Jul 2017",
        stars: 5250,
      },
      {
        date: "Jul 2016",
        stars: 7920,
      },
      {
        date: "Mar 2018",
        stars: 10590,
      },
      {
        date: "Jan 2017",
        stars: 13260,
      },
      {
        date: "Oct 2016",
        stars: 15930,
      },
      {
        date: "Mar 2016",
        stars: 18600,
      },
      {
        date: "Jul 2015",
        stars: 21270,
      },
      {
        date: "Feb 2015",
        stars: 23940,
      },
      {
        date: "Dec 2017",
        stars: 26610,
      },
      {
        date: "May 2018",
        stars: 29280,
      },
      {
        date: "Jul 2018",
        stars: 31950,
      },
      {
        date: "Apr 2017",
        stars: 34620,
      },
      {
        date: "Nov 2015",
        stars: 37290,
      },
      {
        date: "Oct 2017",
        stars: 39960,
      },
      {
        date: "Nov 2024",
        stars: 124227,
      },
    ],
    logo_url: "https://avatars.githubusercontent.com/u/4314092?v=4",
    name: "golang/go",
    total_stars: 124227,
  });

  const [loading, setLoading] = useState(false);

  const getData = async (repo: string) => {
    try {
      setLoading(true);
      const res = await axios.get(
        `http://localhost:5000/history?repo=${repo}`,
        {
          // withCredentials: true,
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      setData(res.data);
      console.log(res.data);
    } catch (error) {
      toast.error(error?.response.data.message || "Something went wrong");
      console.log(error);
      setLoading(false);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <div className="space-y-2">
        <Search getData={getData} loading={loading} />
        <div>
          {token ? (
            <RemoveTokenButton setToken={setToken} />
          ) : (
            <AddTokenDialog setToken={setToken} />
          )}
        </div>
      </div>
      {data && (
        <div>
          {loading ? (
            <div>
              <div className="p-4">
                <div className="space-y-6 mx-auto mt-10 p-6 rounded-lg w-[100%] sm:w-[600px]">
                  <SkeletonCard />
                </div>
              </div>
            </div>
          ) : (
            <Graph data={data} />
          )}
        </div>
      )}
    </div>
  );
}

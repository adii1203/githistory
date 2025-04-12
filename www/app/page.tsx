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

    getData("golang/go");
  }, []);

  const [data, setData] = useState<ChartData>();

  const [loading, setLoading] = useState(false);

  const getData = async (repo: string) => {
    try {
      setLoading(true);
      const res = await axios.get(
        `https://githistory-production-92ae.up.railway.app/history?repo=${repo}`,
        {
          // withCredentials: true,
          headers: {
            Authorization: token && `Bearer ${token}`,
          },
        }
      );
      setData(res.data);
      console.log(res.data);
    } catch (error) {
      const err = error as { response: { data: { message: string } } };
      const msg = err.response?.data?.message;
      toast.error(msg || "Something went wrong");
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

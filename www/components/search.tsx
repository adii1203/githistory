import React, { useState } from "react";
import { Input } from "./ui/input";
import { LoaderIcon } from "lucide-react";

const Search = ({
  loading,
  getData,
}: {
  loading: boolean;
  getData: (repo: string) => void;
}) => {
  const [repo, setRepo] = useState("");

  return (
    <div className="space-y-2 mt-10 w-[20rem] sm:w-[28rem] mx-auto">
      <div className="flex rounded-lg shadow-sm shadow-black/5">
        <Input
          id="input-21"
          className="-me-px flex-1 rounded-e-none shadow-none focus-visible:z-10"
          placeholder="golang/go"
          type="search"
          value={repo}
          onChange={(e) => setRepo(e.target.value)}
        />
        <button
          onClick={() => getData(repo)}
          className="inline-flex items-center rounded-e-lg border border-input bg-background px-3 text-sm font-medium text-foreground outline-offset-2 transition-colors hover:bg-accent hover:text-foreground focus:z-10 focus-visible:outline focus-visible:outline-2 focus-visible:outline-ring/70 disabled:cursor-not-allowed disabled:opacity-50">
          {loading ? (
            <LoaderIcon className="animate-spin" size={16} />
          ) : (
            "Search"
          )}
        </button>
      </div>
    </div>
  );
};

export default Search;

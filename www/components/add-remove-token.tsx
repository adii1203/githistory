import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import Link from "next/link";
import { useState } from "react";
import { toast } from "sonner";

export function AddTokenDialog({
  setToken,
}: {
  setToken: (token: string) => void;
}) {
  const [value, setValue] = useState("");
  const [open, setOpen] = useState(false);

  const addTokenToLocalStorage = () => {
    localStorage.setItem("github_token", value);
    setOpen(false);
    toast.success("Token added successfully");
    setToken(value);
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant={"outline"} className="h-8">
          Add token
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Add github token</DialogTitle>
          <DialogDescription>
            {`gitgraph uses the GitHub API to fetch repo data. You've hit the API
            rate limit. To continue, provide a personal access token (no special
            permissions needed).`}
            <Link
              className="text-blue-600 underline"
              href={"https://github.com/settings/tokens/new"}
              target="_blank">
              {" "}
              Generate{" "}
            </Link>
            one and paste it below. Your token will be stored in local storage.
            you can remove it any time.
          </DialogDescription>
        </DialogHeader>
        <div className="flex items-center space-x-2">
          <div className="grid flex-1 gap-2">
            <Label htmlFor="token" className="sr-only">
              Token
            </Label>
            <Input
              id="token"
              value={value}
              onChange={(e) => setValue(e.target.value)}
            />
          </div>
          <Button onClick={addTokenToLocalStorage} size="sm" className="px-3">
            Add
          </Button>
        </div>
        <DialogFooter className="sm:justify-start">
          <DialogClose asChild>
            <Button type="button" variant="secondary">
              Close
            </Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

export const RemoveTokenButton = ({
  setToken,
}: {
  setToken: (token: string) => void;
}) => {
  return (
    <Button
      onClick={() => {
        localStorage.removeItem("github_token");
        toast.success("Token removed successfully");
        setToken("");
      }}
      variant={"outline"}
      className="h-8">
      Remove token
    </Button>
  );
};

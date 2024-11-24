import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";

export default function SkeletonCard() {
  return (
    <Card className="w-full max-w-3xl">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-7">
        <div className="flex items-center space-x-4">
          <Skeleton className="h-10 w-10 rounded-full" />
          <Skeleton className="h-6 w-32" />
        </div>
        <div className="flex flex-col items-start space-y-2">
          <Skeleton className="h-5 w-24" />
          <Skeleton className="h-8 w-24" />
        </div>
      </CardHeader>
      <CardContent>
        {/* Graph area */}
        <div className="h-[300px]">
          <div className="ml-0 h-full space-y-2">
            <Skeleton className="h-full w-full" />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

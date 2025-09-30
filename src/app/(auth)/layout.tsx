import { cn } from "@/lib/utils";
import { Card } from "@/components/ui/card";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center gap-6 p-6 md:p-10">
      <div className="flex w-full max-w-sm flex-col gap-6">
        <div className="flex flex-col gap-6">
          <Card>{children}</Card>
        </div>
      </div>
    </div>
  );
}

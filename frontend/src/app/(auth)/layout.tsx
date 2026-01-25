import { Card, CardContent } from "@/components/ui/card";
import { FieldDescription } from "@/components/ui/field";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <div className="bg-muted flex min-h-svh flex-col items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm md:max-w-4xl">
        <div className="flex flex-col gap-6">
          <Card className="overflow-hidden p-0">{children}</Card>
          <FieldDescription className="px-6 text-center">
            Created by{" "}
            <a href="https://github.com/insanXYZ/whatsupp">insanXYZ</a>.
          </FieldDescription>
        </div>
      </div>
    </div>
  );
}

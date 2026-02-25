import { SidebarProvider } from "@/components/ui/sidebar";
import { AccountProvider } from "@/provider/account-provider";
import { ReactNode } from "react";

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <SidebarProvider
      style={
        {
          "--sidebar-width": "350px",
        } as React.CSSProperties
      }
    >
      <AccountProvider>
        {children}
      </AccountProvider>
    </SidebarProvider>
  );
}

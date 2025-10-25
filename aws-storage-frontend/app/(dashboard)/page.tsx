import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/ui/app-sidebar"
import "./globals.css"
import { NavigationMenu, NavigationMenuItem, NavigationMenuLink } from "@/components/ui/navigation-menu"

export const metadata = {
  title: "AWS Storage Dashboard",
  description: "Demo Dashboard using shadcn/ui + Next.js",
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (

    // <SidebarProvider>
    //   <AppSidebar />
    //   <main>
    //     <SidebarTrigger />
    //     {children}
    //   </main>
    // </SidebarProvider>    
    <h1></h1>
  )
}

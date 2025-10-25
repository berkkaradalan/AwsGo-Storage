"use client"
import {
    Sidebar,
    SidebarContent,
    SidebarGroup,
    SidebarGroupContent,
    SidebarGroupLabel,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarFooter,
    SidebarHeader
  } from "@/components/ui/sidebar"

  import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"

  import { Calendar, ContainerIcon, GithubIcon, Home, Inbox, Search, Settings } from "lucide-react"

  const items = [
    {
      title: "Home",
      url: "/",
      icon: Home,
    },
    {
      title: "Storage",
      url: "/storage",
      icon: ContainerIcon,
    },
    {
      title: "About",
      url: "/about",
      icon: GithubIcon,
    },
  ]
  
  export function AppSidebar() {
    return (
      <Sidebar>
        <SidebarHeader>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton size="lg" asChild>
                <a href="/" className="flex items-center gap-3">
                  <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-orange-400 to-orange-600 shadow-sm">
                    <img src="/logo.png" alt="AWS Logo" className="h-8 w-8 object-contain" />
                  </div>
                  <div className="flex flex-col gap-0.5">
                    <span className="font-semibold">AWS Storage</span>
                    <span className="text-xs text-muted-foreground">Cloud Manager</span>
                  </div>
                </a>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarHeader>


        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Application</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {items.map((item) => (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton asChild>
                      <a href={item.url}>
                        <item.icon />
                        <span>{item.title}</span>
                      </a>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>


        <SidebarFooter>
            <SidebarMenu>
                <SidebarMenuItem>
                    <SidebarMenuButton asChild>
                        <a href="#profile" className="flex items-center gap-2">
                            <Avatar className="h-8 w-8">
                                <AvatarImage src="https://github.com/shadcn.png"></AvatarImage>
                                <AvatarFallback>CN</AvatarFallback>
                            </Avatar>
                            <div className="flex flex-col flex-1 text-left text-sm">
                                <span className="font-medium">Berk Karadalan</span>
                                <span className="text-xs text-muted-foreground">berkkaradalan@gmail.com</span>
                            </div>
                        </a>
                    </SidebarMenuButton>
                </SidebarMenuItem>
            </SidebarMenu>
        </SidebarFooter>
      </Sidebar>
    )
  }
  
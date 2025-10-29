"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"
import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/ui/app-sidebar"
import { useAuth } from "@/app/contexts/AuthContext"
import { StorageProvider } from "@/app/contexts/StorageContext"
import { Toaster } from "sonner"


export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, loading } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (!loading && !isAuthenticated) {
      router.push("/login")
    }
  }, [isAuthenticated, loading, router])

  // Loading state
  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="text-muted-foreground">Loading...</div>
      </div>
    )
  }

  // Not authenticated
  if (!isAuthenticated) {
    return null
  }

  // Authenticated - show dashboard
  return (
    <>
    <StorageProvider>
      <SidebarProvider>
        <AppSidebar />
        <main className="flex-1 overflow-y-auto">
          <SidebarTrigger />
          {children}
        </main>
      </SidebarProvider>
    </StorageProvider>
    <Toaster position="top-right" richColors />
    </>
  )
}
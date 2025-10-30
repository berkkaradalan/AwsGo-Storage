"use client"

import { createContext, useContext, useState, useEffect, ReactNode } from "react"
import { toast } from "sonner"

interface MonthlyUsage {
  month: string
  monthName: string
  totalSize: number
  fileCount: number
  sizeInMB: number
}

interface DashboardSummary {
  totalSizeInBytes: number
  totalSizeInMB: number
  totalSizeInGB: number
  totalFiles: number
}

interface DashboardData {
  months: MonthlyUsage[]
  summary: DashboardSummary
}

interface DashboardResponse {
  success: boolean
  message: string
  data: DashboardData
  error?: string
}

interface DashboardContextType {
  dashboardData: DashboardData | null
  loading: boolean
  error: string | null
  refetch: () => Promise<void>
}

const DashboardContext = createContext<DashboardContextType | undefined>(undefined)

export function DashboardProvider({ children }: { children: ReactNode }) {
  const [dashboardData, setDashboardData] = useState<DashboardData | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [mounted, setMounted] = useState(false)

  const fetchDashboard = async () => {
    if (typeof window === 'undefined') return
    
    try {
      setLoading(true)
      setError(null)

      const token = localStorage.getItem("auth_token")
      if (!token) {
        throw new Error("No authentication token found")
      }

      const response = await fetch("http://localhost:8080/api/v1/storage/dashboard", {
        method: "GET",
        headers: {
          "Authorization": `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.error || `HTTP error! status: ${response.status}`)
      }

      const data: DashboardResponse = await response.json()
      
      if (!data.success) {
        throw new Error(data.error || "Failed to fetch dashboard data")
      }

      setDashboardData(data.data)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to fetch dashboard data"
      setError(errorMessage)
      console.error("Dashboard fetch error:", err)
      toast.error("Dashboard Error", {
        description: errorMessage,
      })
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    setMounted(true)
  }, [])

  useEffect(() => {
    if (mounted) {
      fetchDashboard()
    }
  }, [mounted])

  return (
    <DashboardContext.Provider value={{ dashboardData, loading, error, refetch: fetchDashboard }}>
      {children}
    </DashboardContext.Provider>
  )
}

export function useDashboard() {
  const context = useContext(DashboardContext)
  if (context === undefined) {
    throw new Error("useDashboard must be used within a DashboardProvider")
  }
  return context
}

export default useDashboard
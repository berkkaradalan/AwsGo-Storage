"use client"

import { Bar, BarChart, CartesianGrid, XAxis } from "recharts"
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart"

import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

import { useDashboard } from "../contexts/DashboardContext"
import { useStorage } from "@/app/contexts/StorageContext"
import { Skeleton } from "@/components/ui/skeleton"
import { AlertCircle } from "lucide-react"

const chartConfig = {
  storage: {
    label: "Storage (MB)",
    color: "#fc6c00",
  },
} satisfies ChartConfig

export default function DashboardPage() {
  const { dashboardData, loading: dashboardLoading, error: dashboardError } = useDashboard()
  const { files, loading: filesLoading } = useStorage()

  const chartData = dashboardData?.months.map(month => ({
    month: month.monthName,
    storage: month.sizeInMB,
  })) || []

  const lastFiles = files
    .sort((a, b) => new Date(b.UploadedAt).getTime() - new Date(a.UploadedAt).getTime())
    .slice(0, 12)

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
  }

  const getFileType = (contentType: string): string => {
    if (contentType.includes('pdf')) return 'PDF'
    if (contentType.includes('excel') || contentType.includes('spreadsheet')) return 'Excel'
    if (contentType.includes('word') || contentType.includes('document')) return 'Word'
    if (contentType.includes('powerpoint') || contentType.includes('presentation')) return 'PowerPoint'
    if (contentType.includes('image')) return 'Image'
    return 'File'
  }

  const formatDate = (dateString: string): string => {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', { year: 'numeric', month: '2-digit', day: '2-digit' })
  }

  if (dashboardLoading || filesLoading) {
    return (
      <div className="container mx-auto p-6">
        <div className="rounded-lg border bg-card p-6">
          <Skeleton className="h-8 w-64 mb-4" />
          <Skeleton className="h-[300px] w-full mb-6" />
          <Skeleton className="h-[400px] w-full" />
        </div>
      </div>
    )
  }

  if (dashboardError || !dashboardData) {
    return (
      <div className="container mx-auto p-6">
        <div className="rounded-lg border bg-card p-6">
          <div className="flex items-center justify-center h-[300px]">
            <div className="text-center">
              <AlertCircle className="mx-auto h-12 w-12 text-muted-foreground mb-4" />
              <h3 className="text-lg font-semibold mb-2">Unable to load dashboard</h3>
              <p className="text-sm text-muted-foreground">
                Dashboard data could not be loaded. Please try again later.
              </p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="container mx-auto p-6">
      <div className="rounded-lg border bg-card p-6">
        {/* Storage Usage Chart */}
        <div className="mb-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold">Storage Usage By Months</h2>
            <div className="text-sm text-muted-foreground">
              Total: {dashboardData.summary.totalSizeInMB.toFixed(2)} MB
            </div>
          </div>
          
          {chartData.length > 0 ? (
            <ChartContainer config={chartConfig} className="h-[300px] w-full">
              <BarChart data={chartData} barSize={50}>
                <CartesianGrid vertical={false} />
                <XAxis
                  dataKey="month"
                  tickLine={false}
                  tickMargin={10}
                  axisLine={false}
                  tickFormatter={(value) => value}
                />
                <ChartTooltip 
                  content={<ChartTooltipContent />}
                  formatter={(value) => `${Number(value).toFixed(2)} MB`}
                />
                <Bar dataKey="storage" fill="var(--color-storage)" radius={8} />
              </BarChart>
            </ChartContainer>
          ) : (
            <div className="flex items-center justify-center h-[300px] text-muted-foreground">
              No storage data available
            </div>
          )}
        </div>

        {/* Summary Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
          <div className="p-4 rounded-lg border bg-muted/50">
            <p className="text-sm text-muted-foreground">Total Files</p>
            <p className="text-2xl font-bold">{dashboardData.summary.totalFiles}</p>
          </div>
          <div className="p-4 rounded-lg border bg-muted/50">
            <p className="text-sm text-muted-foreground">Total Storage</p>
            <p className="text-2xl font-bold">{dashboardData.summary.totalSizeInMB.toFixed(2)} MB</p>
          </div>
          <div className="p-4 rounded-lg border bg-muted/50">
            <p className="text-sm text-muted-foreground">Storage (GB)</p>
            <p className="text-2xl font-bold">{dashboardData.summary.totalSizeInGB.toFixed(3)} GB</p>
          </div>
        </div>

        {/* Recent Files Table */}
        <div>
          <h2 className="text-xl font-semibold mb-4">Recent Files</h2>
          <Table>
            <TableCaption>A list of your recent files.</TableCaption>
            <TableHeader>
              <TableRow>
                <TableHead>File Name</TableHead>
                <TableHead>Type</TableHead>
                <TableHead>Size</TableHead>
                <TableHead className="text-right">Upload Date</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {lastFiles.length > 0 ? (
                lastFiles.map((file) => (
                  <TableRow key={file.ObjectID}>
                    <TableCell className="font-medium">{file.FileName}</TableCell>
                    <TableCell>{getFileType(file.ContentType)}</TableCell>
                    <TableCell>{formatFileSize(file.FileSize)}</TableCell>
                    <TableCell className="text-right">{formatDate(file.UploadedAt)}</TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell colSpan={4} className="text-center text-muted-foreground">
                    No files uploaded yet
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>
      </div>
    </div>
  )
}
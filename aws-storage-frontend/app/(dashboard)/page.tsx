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
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

// Mock data
const lastFiles = [
  {
    fileName: "report_january.pdf",
    fileType: "PDF",
    fileSize: "1.2 MB",
    uploadDate: "2025-01-15",
  },
  {
    fileName: "invoice_february.xlsx",
    fileType: "Excel",
    fileSize: "850 KB",
    uploadDate: "2025-02-10",
  },
  {
    fileName: "contract_march.docx",
    fileType: "Word",
    fileSize: "2.3 MB",
    uploadDate: "2025-03-05",
  },
  {
    fileName: "receipt_april.pdf",
    fileType: "PDF",
    fileSize: "1.0 MB",
    uploadDate: "2025-04-20",
  },
  {
    fileName: "report_may.xlsx",
    fileType: "Excel",
    fileSize: "1.5 MB",
    uploadDate: "2025-05-18",
  },
  {
    fileName: "summary_june.docx",
    fileType: "Word",
    fileSize: "950 KB",
    uploadDate: "2025-06-12",
  },
  {
    fileName: "invoice_july.pdf",
    fileType: "PDF",
    fileSize: "1.1 MB",
    uploadDate: "2025-07-07",
  },
  {
    fileName: "presentation_august.pptx",
    fileType: "PowerPoint",
    fileSize: "2.0 MB",
    uploadDate: "2025-08-03",
  },
  {
    fileName: "summary_september.docx",
    fileType: "Word",
    fileSize: "1.3 MB",
    uploadDate: "2025-09-14",
  },
  {
    fileName: "report_october.pdf",
    fileType: "PDF",
    fileSize: "1.4 MB",
    uploadDate: "2025-10-10",
  },
  {
    fileName: "invoice_november.xlsx",
    fileType: "Excel",
    fileSize: "1.0 MB",
    uploadDate: "2025-11-05",
  },
  {
    fileName: "contract_december.docx",
    fileType: "Word",
    fileSize: "2.1 MB",
    uploadDate: "2025-12-20",
  },
]


const chartData = [
  { month: "January", storage: 186 },
  { month: "February", storage: 305 },
  { month: "March", storage: 237 },
  { month: "April", storage: 273 },
  { month: "May", storage: 209 },
  { month: "June", storage: 214 },
  { month: "July", storage: 298 },
  { month: "August", storage: 256 },
  { month: "September", storage: 332 },
  { month: "October", storage: 287 },
  { month: "November", storage: 240 },
  { month: "December", storage: 310 },
]


const chartConfig = {
  storage: {
    label: "Storage (GB)",
    color: "#fc6c00",
  },
} satisfies ChartConfig

export default function DashboardPage() {
  return (
    <div className="container mx-auto p-6">
      
      <div className="rounded-lg border bg-card p-6">
        <h2 className="text-xl font-semibold mb-4">Storage Usage By Months</h2>
        <ChartContainer config={chartConfig} className="h-[300px] w-full">
          <BarChart data={chartData} barSize={50}>
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="month"
              tickLine={false}
              tickMargin={10}
              axisLine={false}
              tickFormatter={(value) => value.slice(0, 5)}
            />
            <ChartTooltip content={<ChartTooltipContent />} />
            <Bar dataKey="storage" fill="var(--color-storage)" radius={8} />
          </BarChart>
        </ChartContainer>

        <Table>
          <TableCaption>A list of your recent files.</TableCaption>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[100px]">Invoice</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Method</TableHead>
              <TableHead className="text-right">Amount</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {lastFiles.map((lastFiles) => (
              <TableRow key={lastFiles.fileName}>
                <TableCell>{lastFiles.fileName}</TableCell>
                <TableCell className="font-medium">{lastFiles.fileSize}</TableCell>
                <TableCell>{lastFiles.fileType}</TableCell>
                <TableCell className="text-right">{lastFiles.uploadDate}</TableCell>
              </TableRow>
            ))}
          </TableBody>
      </Table>
      </div>
    </div>
  )
}
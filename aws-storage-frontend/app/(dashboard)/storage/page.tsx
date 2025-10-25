"use client"

import { useState, useEffect } from "react"
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Download, Trash2, FileImage, Eye } from "lucide-react"

const ITEMS_PER_PAGE = 12

// todo - Mock data
const mockFiles = Array.from({ length: 47 }, (_, i) => ({
  id: `file-${i + 1}`,
  name: `image-${i + 1}.jpg`,
  url: `https://picsum.photos/seed/${i + 1}/400/300`,
  size: `${1000 + (i * 100)} KB`,
  uploadDate: `2024-01-${String((i % 30) + 1).padStart(2, '0')}`,
  status: i % 10 === 0 ? "processing" : "active"
}))

export default function StoragePage() {
  const [currentPage, setCurrentPage] = useState(1)
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  const totalPages = Math.ceil(mockFiles.length / ITEMS_PER_PAGE)
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE
  const endIndex = startIndex + ITEMS_PER_PAGE
  const currentFiles = mockFiles.slice(startIndex, endIndex)

  if (!mounted) {
    return (
      <div className="flex flex-col gap-6 p-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold">Storage</h1>
            <p className="text-muted-foreground mt-2">Loading...</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-6 p-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Storage</h1>
          <p className="text-muted-foreground mt-2">
            Total {mockFiles.length} images • Page {currentPage} of {totalPages}
          </p>
        </div>
        <Button>
          <FileImage className="mr-2 h-4 w-4" />
          Upload Image
        </Button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {currentFiles.map((file) => (
          <Card key={file.id} className="overflow-hidden">
            <div className="aspect-video relative bg-muted">
              <img 
                src={file.url} 
                alt={file.name}
                className="w-full h-full object-cover"
              />
              <Badge 
                variant={file.status === "active" ? "default" : "secondary"}
                className="absolute top-2 right-2"
              >
                {file.status}
              </Badge>
            </div>
            <CardHeader className="p-4">
              <CardTitle className="text-sm truncate">{file.name}</CardTitle>
              <CardDescription className="text-xs">
                {file.size} • {file.uploadDate}
              </CardDescription>
            </CardHeader>
            <CardFooter className="p-4 pt-0 flex gap-2">
              <Button variant="outline" size="sm" className="flex-1">
                <Eye className="h-4 w-4 mr-1" />
                View
              </Button>
              <Button variant="outline" size="sm">
                <Download className="h-4 w-4" />
              </Button>
              <Button variant="outline" size="sm">
                <Trash2 className="h-4 w-4" />
              </Button>
            </CardFooter>
          </Card>
        ))}
      </div>

      <Pagination>
        <PaginationContent>
          <PaginationItem>
            <PaginationPrevious 
              onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
              className={currentPage === 1 ? "pointer-events-none opacity-50" : "cursor-pointer"}
            />
          </PaginationItem>
          
          {[...Array(Math.min(5, totalPages))].map((_, i) => {
            const pageNum = i + 1
            return (
              <PaginationItem key={pageNum}>
                <PaginationLink
                  onClick={() => setCurrentPage(pageNum)}
                  isActive={currentPage === pageNum}
                  className="cursor-pointer"
                >
                  {pageNum}
                </PaginationLink>
              </PaginationItem>
            )
          })}
          
          <PaginationItem>
            <PaginationNext 
              onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
              className={currentPage === totalPages ? "pointer-events-none opacity-50" : "cursor-pointer"}
            />
          </PaginationItem>
        </PaginationContent>
      </Pagination>
    </div>
  )
}
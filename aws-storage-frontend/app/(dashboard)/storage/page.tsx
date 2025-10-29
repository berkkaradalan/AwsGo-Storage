"use client"

import { useState, useEffect } from "react"
import { useStorage } from "@/app/contexts/StorageContext"
import { UploadDialog } from "@/components/ui/upload-dialog"
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Download, Trash2, FileImage, Eye, Loader2 } from "lucide-react"
import { toast } from "sonner"

const ITEMS_PER_PAGE = 12

export default function StoragePage() {
  const { files, loading, error, refetch, deleteFile, downloadFile } = useStorage()
  const [currentPage, setCurrentPage] = useState(1)
  const [mounted, setMounted] = useState(false)
  const [uploadDialogOpen, setUploadDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [fileToDelete, setFileToDelete] = useState<{ id: string, name: string } | null>(null)
  const [deleting, setDeleting] = useState(false)
  const [downloadingId, setDownloadingId] = useState<string | null>(null)

  useEffect(() => {
    setMounted(true)
  }, [])

  const handleDelete = async () => {
    if (!fileToDelete) return

    try {
        setDeleting(true)
        await deleteFile(fileToDelete.id)
        toast.success(`${fileToDelete.name} deleted successfully`)
        setDeleteDialogOpen(false)
        setFileToDelete(null)
      } catch (error) {
        toast.error(error instanceof Error ? error.message : "Failed to delete file")
      } finally {
        setDeleting(false)
      }
  }

  const handleDownload = async (fileId: string, fileName: string) => {
    try {
      setDownloadingId(fileId)
      await downloadFile(fileId, fileName)
      toast.success(`${fileName} downloaded successfully`)
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to download file")
    } finally {
      setDownloadingId(null)
    }
  }

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`
    return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', { 
      year: 'numeric', 
      month: 'short', 
      day: 'numeric' 
    })
  }

  const totalPages = Math.ceil(files.length / ITEMS_PER_PAGE)
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE
  const endIndex = startIndex + ITEMS_PER_PAGE
  const currentFiles = files.slice(startIndex, endIndex)

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

  if (loading) {
    return (
      <div className="flex flex-col gap-6 p-6">
        <div className="flex items-center justify-center h-64">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex flex-col gap-6 p-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold">Storage</h1>
            <p className="text-red-500 mt-2">Error: {error}</p>
          </div>
          <Button onClick={refetch}>
            Retry
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-6 p-6">
      <UploadDialog open={uploadDialogOpen} onOpenChange={setUploadDialogOpen} />
      
      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently delete <span className="font-semibold text-foreground">{fileToDelete?.name}</span> from your storage.
              This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={deleting}>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDelete}
              disabled={deleting}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {deleting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Deleting...
                </>
              ) : (
                'Delete'
              )}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
      
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Storage</h1>
          <p className="text-muted-foreground mt-2">
            {files.length === 0 
              ? "No files uploaded yet" 
              : `Total ${files.length} files${totalPages > 0 ? ` • Page ${currentPage} of ${totalPages}` : ''}`
            }
          </p>
        </div>
        <Button onClick={() => setUploadDialogOpen(true)}>
          <FileImage className="mr-2 h-4 w-4" />
          Upload File
        </Button>
      </div>

      {files.length === 0 ? (
        <div className="flex flex-col items-center justify-center h-64 border-2 border-dashed rounded-lg">
          <FileImage className="h-12 w-12 text-muted-foreground mb-4" />
          <p className="text-muted-foreground">No files uploaded yet</p>
          <Button className="mt-4" variant="outline" onClick={() => setUploadDialogOpen(true)}>
            Upload your first file
          </Button>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            {currentFiles.map((file) => (
              <Card key={file.ObjectID} className="overflow-hidden flex flex-col h-full">
                <div className="aspect-video relative bg-muted flex-shrink-0 overflow-hidden">
                  {file.ContentType === 'application/pdf' ? (
                    <div className="w-full h-full flex items-center justify-center bg-gradient-to-br from-red-50 to-red-100 dark:from-red-950 dark:to-red-900">
                      <FileImage className="h-20 w-20 text-red-500" />
                    </div>
                  ) : (
                    <img 
                      src={file.previewUrl} 
                      alt={file.FileName}
                      className="w-full h-full object-cover"
                      loading="lazy"
                      onError={(e) => {
                        const target = e.currentTarget
                        target.style.display = 'none'
                        const parent = target.parentElement
                        if (parent) {
                          parent.innerHTML = `
                            <div class="w-full h-full flex items-center justify-center bg-gradient-to-br from-gray-100 to-gray-200 dark:from-gray-800 dark:to-gray-900">
                              <svg class="h-20 w-20 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                              </svg>
                            </div>
                          `
                        }
                      }}
                    />
                  )}
                  <Badge 
                    variant="default"
                    className="absolute top-2 right-2"
                  >
                    {file.ContentType.split('/')[1].toUpperCase()}
                  </Badge>
                </div>
                <div className="flex flex-col flex-1 min-h-0">
                  <CardHeader className="p-4 flex-shrink-0">
                    <CardTitle className="text-sm truncate" title={file.FileName}>
                      {file.FileName}
                    </CardTitle>
                    <CardDescription className="text-xs">
                      {formatFileSize(file.FileSize)} • {formatDate(file.UploadedAt)}
                    </CardDescription>
                    {file.Description && (
                      <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
                        {file.Description}
                      </p>
                    )}
                  </CardHeader>
                  <CardFooter className="p-4 pt-0 flex gap-2 mt-auto flex-shrink-0">
                    <Button 
                      variant="outline" 
                      size="sm" 
                      className="flex-1"
                      onClick={() => window.open(file.previewUrl, '_blank')}
                    >
                      <Eye className="h-4 w-4 mr-1" />
                      View
                    </Button>
                    <Button 
                      variant="outline" 
                      size="sm"
                      onClick={() => handleDownload(file.ObjectID, file.FileName)}
                      disabled={downloadingId === file.ObjectID}
                    >
                      {downloadingId === file.ObjectID ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        <Download className="h-4 w-4" />
                      )}
                    </Button>
                    <Button 
                      variant="outline" 
                      size="sm"
                      onClick={() => {
                        setFileToDelete({ id: file.ObjectID, name: file.FileName })
                        setDeleteDialogOpen(true)
                      }}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </CardFooter>
                </div>
              </Card>
            ))}
          </div>

          {totalPages > 1 && (
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
          )}
        </>
      )}
    </div>
  )
}
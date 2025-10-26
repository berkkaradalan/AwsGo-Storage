"use client"

import { useState } from "react"
import { useStorage } from "@/app/contexts/StorageContext"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Upload, Loader2, FileImage } from "lucide-react"
import { toast } from "sonner"

interface UploadDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function UploadDialog({ open, onOpenChange }: UploadDialogProps) {
  const { uploadFile, uploading } = useStorage()
  const [file, setFile] = useState<File | null>(null)
  const [description, setDescription] = useState("")
  const [preview, setPreview] = useState<string | null>(null)

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0]
    if (selectedFile) {
      setFile(selectedFile)
      
      if (selectedFile.type.startsWith('image/')) {
        const reader = new FileReader()
        reader.onloadend = () => {
          setPreview(reader.result as string)
        }
        reader.readAsDataURL(selectedFile)
      } else {
        setPreview(null)
      }
    }
  }

  const handleUpload = async () => {
    if (!file) {
      toast.error("Please select a file")
      return
    }

    try {
      await uploadFile(file, description || undefined)
      toast.success("File uploaded successfully!")
      
      setFile(null)
      setDescription("")
      setPreview(null)
      onOpenChange(false)
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Upload failed")
    }
  }

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return `${bytes} B`
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`
    return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Upload File</DialogTitle>
          <DialogDescription>
            Upload images (JPEG, PNG, GIF, WebP) or PDF files. Max size: 50MB
          </DialogDescription>
        </DialogHeader>
        
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="file">File</Label>
            <Input
              id="file"
              type="file"
              accept="image/jpeg,image/png,image/gif,image/webp,application/pdf"
              onChange={handleFileChange}
              disabled={uploading}
            />
          </div>

          {file && (
            <div className="rounded-lg border p-4 space-y-3">
              {preview && (
                <div className="aspect-video relative bg-muted rounded-md overflow-hidden">
                  <img 
                    src={preview} 
                    alt="Preview"
                    className="w-full h-full object-contain"
                  />
                </div>
              )}
              {!preview && (
                <div className="aspect-video relative bg-muted rounded-md overflow-hidden flex items-center justify-center">
                  <FileImage className="h-16 w-16 text-muted-foreground" />
                </div>
              )}
              <div className="text-sm space-y-1">
                <p className="font-medium truncate">{file.name}</p>
                <p className="text-muted-foreground">
                  {formatFileSize(file.size)} • {file.type}
                </p>
              </div>
            </div>
          )}

          <div className="grid gap-2">
            <Label htmlFor="description">Description (Optional)</Label>
            <Textarea
              id="description"
              placeholder="Add a description for your file..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              disabled={uploading}
              rows={3}
            />
          </div>
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={uploading}
          >
            Cancel
          </Button>
          <Button onClick={handleUpload} disabled={!file || uploading}>
            {uploading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Uploading...
              </>
            ) : (
              <>
                <Upload className="mr-2 h-4 w-4" />
                Upload
              </>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
"use client"

import { createContext, useContext, useState, useEffect, ReactNode } from "react"

interface StorageFile {
  ObjectID: string
  UserID: string
  FileName: string
  FileSize: number
  ContentType: string
  S3Key: string
  S3Bucket: string
  UploadedAt: string
  UpdatedAt: string
  Description: string | null
  previewUrl: string
}

interface StorageResponse {
  success: boolean
  message: string
  data: StorageFile[]
  count: number
}

interface StorageContextType {
  files: StorageFile[]
  loading: boolean
  error: string | null
  uploading: boolean
  refetch: () => Promise<void>
  uploadFile: (file: File, description?: string) => Promise<void>
}

const StorageContext = createContext<StorageContextType | undefined>(undefined)

export function StorageProvider({ children }: { children: ReactNode }) {
  const [files, setFiles] = useState<StorageFile[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [mounted, setMounted] = useState(false)
  const [uploading, setUploading] = useState(false)

  const fetchFiles = async () => {
    if (typeof window === 'undefined') return
    
    try {
      setLoading(true)
      setError(null)

      const token = localStorage.getItem("auth_token")
      if (!token) {
        throw new Error("No authentication token found")
      }

      //todo - add dynamic url via env
      const response = await fetch("http://localhost:8080/api/v1/storage/files", {
        method: "GET",
        headers: {
          "Authorization": `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        console.error("Error response:", errorData)
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data: StorageResponse = await response.json()
      setFiles(data.data || [])
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch files")
      console.error("Storage fetch error:", err)
    } finally {
      setLoading(false)
    }
  }

  const uploadFile = async (file: File, description?: string) => {
    if (typeof window === 'undefined') return

    try {
      setUploading(true)
      setError(null)

      const token = localStorage.getItem("auth_token")
      if (!token) {
        throw new Error("No authentication token found")
      }

      const formData = new FormData()
      formData.append("file", file)
      if (description) {
        formData.append("description", description)
      }

      const response = await fetch("http://localhost:8080/api/v1/storage/upload", {
        method: "POST",
        headers: {
          "Authorization": `Bearer ${token}`,
        },
        body: formData,
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.error || `Upload failed with status: ${response.status}`)
      }

      const data = await response.json()
      console.log("Upload successful:", data)

      await fetchFiles()
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to upload file")
      console.error("Upload error:", err)
      throw err
    } finally {
      setUploading(false)
    }
  }

  useEffect(() => {
    setMounted(true)
  }, [])

  useEffect(() => {
    if (mounted) {
      fetchFiles()
    }
  }, [mounted])

  return (
    <StorageContext.Provider value={{ files, loading, error, uploading, refetch: fetchFiles, uploadFile }}>
      {children}
    </StorageContext.Provider>
  )
}

export function useStorage() {
  const context = useContext(StorageContext)
  if (context === undefined) {
    throw new Error("useStorage must be used within a StorageProvider")
  }
  return context
}

export default useStorage
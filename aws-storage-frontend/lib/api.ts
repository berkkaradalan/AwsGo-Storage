//todo - add dynamic url via env
const API_BASE_URL = "http://localhost:8080/api/v1"

async function apiFetch(endpoint: string, options: RequestInit = {}) {
  const token = localStorage.getItem("token")
  
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(token && { Authorization: `Bearer ${token}` }),
      ...options.headers,
    },
  })

  if (!response.ok) {
    const error = await response.json()
    throw error
  }

  return response.json()
}

// Auth endpoints
export const authAPI = {
  register: (data: RegisterData) => 
    apiFetch("/user/register", {
      method: "POST",
      body: JSON.stringify({
        user_name: data.username,
        user_email: data.email,
        user_password: data.password,
      }),
    }),

  login: (data: LoginData) =>
    apiFetch("/user/login", {
      method: "POST",
      body: JSON.stringify({
        user_email: data.email,
        user_password: data.password,
      }),
    }),

  getProfile: () => apiFetch("/user/profile"),
}

export const storageAPI = {
  getFiles: (page: number) => apiFetch(`/files?page=${page}`),
  uploadFile: (file: File) => {
    const formData = new FormData()
    formData.append("file", file)
    return apiFetch("/files/upload", {
      method: "POST",
      body: formData,
      headers: {},
    })
  },
  deleteFile: (id: string) => 
    apiFetch(`/files/${id}`, { method: "DELETE" }),
}

// Types
export interface RegisterData {
  username: string
  email: string
  password: string
}

export interface LoginData {
  email: string
  password: string
}
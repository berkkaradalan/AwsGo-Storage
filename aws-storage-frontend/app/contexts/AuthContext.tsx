"use client"

import { createContext, useContext, useState, useEffect } from "react"
import { authStorage } from "@/lib/auth"
import { authAPI } from "@/lib/api"

interface User {
  user_id: string
  user_name: string
  user_email: string
}

interface AuthContextType {
  user: User | null
  loading: boolean
  login: (email: string, password: string) => Promise<void>
  logout: () => void
  isAuthenticated: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    // Sayfa yüklendiğinde token kontrolü
    const token = authStorage.getToken()
    const savedUser = authStorage.getUser()
    
    if (token && savedUser) {
      setUser(savedUser)
    }
    setLoading(false)
  }, [])

  const login = async (email: string, password: string) => {
    const response = await authAPI.login({ email, password })
    
    // Token ve user'ı kaydet
    authStorage.setToken(response.token)
    authStorage.setUser(response.user)
    setUser(response.user)
  }

  const logout = () => {
    authStorage.clear()
    setUser(null)
  }

  return (
    <AuthContext.Provider 
      value={{ 
        user, 
        loading, 
        login, 
        logout, 
        isAuthenticated: !!user 
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider")
  }
  return context
}
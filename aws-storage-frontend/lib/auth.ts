const TOKEN_KEY = "auth_token"
const USER_KEY = "auth_user"

export const authStorage = {
  setToken: (token: string) => {
    localStorage.setItem(TOKEN_KEY, token)
  },
  
  getToken: () => {
    return localStorage.getItem(TOKEN_KEY)
  },
  
  removeToken: () => {
    localStorage.removeItem(TOKEN_KEY)
  },
  
  setUser: (user: any) => {
    localStorage.setItem(USER_KEY, JSON.stringify(user))
  },
  
  getUser: () => {
    const user = localStorage.getItem(USER_KEY)
    return user ? JSON.parse(user) : null
  },
  
  removeUser: () => {
    localStorage.removeItem(USER_KEY)
  },
  
  clear: () => {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
  },
}
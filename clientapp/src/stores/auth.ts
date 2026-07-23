import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export interface User {
  id: string
  username: string
  email: string
  role: 'admin' | 'author' | 'player'
  team_id?: string
}

interface AuthState {
  token: string | null
  user: User | null
  setToken: (token: string | null) => void
  setUser: (user: User | null) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      user: null,
      setToken: (token: string | null) => set({ token }),
      setUser: (user: User | null) => set({ user }),
      logout: () => set({ token: null, user: null }),
    }),
    { name: 'yctf-auth' }
  )
)

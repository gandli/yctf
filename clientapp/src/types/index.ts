export interface User {
  id: string
  username: string
  email: string
  role: 'admin' | 'author' | 'player'
  team_id?: string
  team?: Team
}

export interface Team {
  id: string
  name: string
  captain_id?: string
  score: number
  invite_code?: string
}

export interface Challenge {
  id: string
  title: string
  description: string
  category: 'web' | 'pwn' | 'crypto' | 're' | 'misc' | 'forensics' | 'osint'
  points: number
  solves: number
  is_solved?: boolean
  container_image?: string
}

export interface Instance {
  id: string
  challenge_id: string
  team_id: string
  container_id: string
  image: string
  status: string
  port: number
  flag: string
  expires_at: string
  created_at: string
}

export interface Submission {
  id: string
  user_id: string
  team_id: string
  challenge_id: string
  flag_submitted: string
  is_correct: boolean
  submitted_at: string
}

export interface ScoreboardEntry {
  rank: number
  team_id: string
  team_name: string
  score: number
  solves: number
}

export interface Writeup {
  id: string
  challenge_id: string
  team_id: string
  user_id: string
  url?: string
  content?: string
  is_approved: boolean
  score: number
  reviewed_by?: string
}

export interface AdminStats {
  users: number
  teams: number
  challenges: number
  submissions: number
  containers: number
  writeups: number
  uptime: string
}

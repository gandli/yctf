import { AppShell, Group, Text, NavLink, ActionIcon, useMantineColorScheme } from '@mantine/core'
import { IconSun, IconMoon, IconFlag, IconTrophy, IconUser, IconShield } from '@tabler/icons-react'
import { useNavigate, useLocation, Outlet } from 'react-router-dom'
import { useAuthStore } from '../stores/auth'

export function Layout() {
  const navigate = useNavigate()
  const location = useLocation()
  const { toggleColorScheme } = useMantineColorScheme()
  const user = useAuthStore((s) => s.user)
  const logout = useAuthStore((s) => s.logout)

  const navItems = [
    { label: 'Challenges', icon: <IconFlag size={18} />, path: '/challenges' },
    { label: 'Scoreboard', icon: <IconTrophy size={18} />, path: '/scoreboard' },
    { label: 'Profile', icon: <IconUser size={18} />, path: '/profile' },
  ]

  if (user?.role === 'admin') {
    navItems.push({ label: 'Admin', icon: <IconShield size={18} />, path: '/admin' })
  }

  return (
    <AppShell header={{ height: 60 }} navbar={{ width: 250, breakpoint: 'sm' }} padding="md">
      <AppShell.Header>
        <Group h="100%" px="md" justify="space-between">
          <Text size="xl" fw={700} onClick={() => navigate('/')} style={{ cursor: 'pointer' }}>
            YCTF
          </Text>
          <Group>
            <Text size="sm" c="dimmed">{user?.username}</Text>
            <ActionIcon variant="subtle" onClick={toggleColorScheme}>
              <IconSun size={18} />
            </ActionIcon>
            <ActionIcon variant="subtle" color="red" onClick={() => { logout(); navigate('/login') }}>
              Logout
            </ActionIcon>
          </Group>
        </Group>
      </AppShell.Header>

      <AppShell.Navbar p="md">
        {navItems.map((item) => (
          <NavLink
            key={item.path}
            label={item.label}
            leftSection={item.icon}
            active={location.pathname.startsWith(item.path)}
            onClick={() => navigate(item.path)}
          />
        ))}
      </AppShell.Navbar>

      <AppShell.Main>
        <Outlet />
      </AppShell.Main>
    </AppShell>
  )
}

import { useQuery } from '@tanstack/react-query'
import { Card, Text, Stack, Title, Table, Badge, Group } from '@mantine/core'
import { api } from '../utils/api'

export function AdminPage() {
  const { data: stats } = useQuery({
    queryKey: ['admin-stats'],
    queryFn: async () => {
      const res = await api.get('/admin/stats')
      return res.data
    },
  })

  const { data: users } = useQuery({
    queryKey: ['admin-users'],
    queryFn: async () => {
      const res = await api.get('/admin/users')
      return res.data.users || []
    },
  })

  return (
    <Stack>
      <Title order={2}>Admin Panel</Title>
      
      <Group grow>
        <Card withBorder p="md">
          <Text size="sm" c="dimmed">Users</Text>
          <Text size="xl" fw={700}>{stats?.users ?? 0}</Text>
        </Card>
        <Card withBorder p="md">
          <Text size="sm" c="dimmed">Challenges</Text>
          <Text size="xl" fw={700}>{stats?.challenges ?? 0}</Text>
        </Card>
        <Card withBorder p="md">
          <Text size="sm" c="dimmed">Submissions</Text>
          <Text size="xl" fw={700}>{stats?.submissions ?? 0}</Text>
        </Card>
        <Card withBorder p="md">
          <Text size="sm" c="dimmed">Uptime</Text>
          <Text size="xl" fw={700}>{stats?.uptime ?? 'N/A'}</Text>
        </Card>
      </Group>

      <Title order={3}>Users</Title>
      <Table>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Username</Table.Th>
            <Table.Th>Email</Table.Th>
            <Table.Th>Role</Table.Th>
            <Table.Th>Status</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {users.map((user: any) => (
            <Table.Tr key={user.id}>
              <Table.Td>{user.username}</Table.Td>
              <Table.Td>{user.email}</Table.Td>
              <Table.Td><Badge size="sm">{user.role}</Badge></Table.Td>
              <Table.Td>
                <Badge color={user.is_banned ? 'red' : 'green'} size="sm">
                  {user.is_banned ? 'Banned' : 'Active'}
                </Badge>
              </Table.Td>
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>
    </Stack>
  )
}

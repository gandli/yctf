import { useQuery } from '@tanstack/react-query'
import { Card, Text, Stack, Title, Table } from '@mantine/core'
import { useAuthStore } from '../stores/auth'
import { api } from '../utils/api'

export function ProfilePage() {
  const user = useAuthStore((s) => s.user)

  const { data: submissions } = useQuery({
    queryKey: ['my-submissions'],
    queryFn: async () => {
      const res = await api.get('/users/me/submissions')
      return res.data.submissions || []
    },
  })

  return (
    <Stack>
      <Title order={2}>Profile</Title>
      <Card withBorder p="md">
        <Text><strong>Username:</strong> {user?.username}</Text>
        <Text><strong>Email:</strong> {user?.email}</Text>
        <Text><strong>Role:</strong> {user?.role}</Text>
      </Card>
      <Title order={3}>My Submissions</Title>
      <Table>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Challenge</Table.Th>
            <Table.Th>Result</Table.Th>
            <Table.Th>Time</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {submissions?.map((sub: any) => (
            <Table.Tr key={sub.id}>
              <Table.Td>{sub.challenge_id}</Table.Td>
              <Table.Td>{sub.is_correct ? '✅ Correct' : '❌ Wrong'}</Table.Td>
              <Table.Td>{sub.submitted_at}</Table.Td>
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>
    </Stack>
  )
}

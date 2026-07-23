import { useQuery } from '@tanstack/react-query'
import { Table, Text, Badge, Stack, Title, Group } from '@mantine/core'
import { IconTrophy } from '@tabler/icons-react'
import { api } from '../utils/api'
import type { ScoreboardEntry } from '../types'

export function ScoreboardPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['scoreboard'],
    queryFn: async () => {
      const res = await api.get('/scoreboard')
      return res.data.teams as ScoreboardEntry[]
    },
    refetchInterval: 5000,
  })

  if (isLoading) return <Text>Loading...</Text>

  const rows = data?.map((entry) => (
    <Table.Tr key={entry.team_id}>
      <Table.Td>
        <Group gap="xs">
          {entry.rank <= 3 && (
            <IconTrophy size={16} color={entry.rank === 1 ? '#ffd700' : entry.rank === 2 ? '#c0c0c0' : '#cd7f32'} />
          )}
          <Text fw={entry.rank <= 3 ? 700 : 400}>#{entry.rank}</Text>
        </Group>
      </Table.Td>
      <Table.Td>{entry.team_name}</Table.Td>
      <Table.Td>
        <Badge variant="light" color="blue">{entry.score}</Badge>
      </Table.Td>
      <Table.Td>{entry.solves}</Table.Td>
    </Table.Tr>
  ))

  return (
    <Stack>
      <Title order={2}>Scoreboard</Title>
      <Table>
        <Table.Thead>
          <Table.Tr>
            <Table.Th>Rank</Table.Th>
            <Table.Th>Team</Table.Th>
            <Table.Th>Score</Table.Th>
            <Table.Th>Solves</Table.Th>
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>{rows}</Table.Tbody>
      </Table>
    </Stack>
  )
}

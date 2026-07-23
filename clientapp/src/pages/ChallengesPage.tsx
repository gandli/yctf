import { useQuery } from '@tanstack/react-query'
import { Card, Text, Badge, Group, Stack, Title, SimpleGrid } from '@mantine/core'
import { IconFlag } from '@tabler/icons-react'
import { api } from '../utils/api'
import type { Challenge } from '../types'

export function ChallengesPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['challenges'],
    queryFn: async () => {
      const res = await api.get('/challenges')
      return res.data.challenges as Challenge[]
    },
  })

  if (isLoading) return <Text>Loading...</Text>

  return (
    <Stack>
      <Title order={2}>Challenges</Title>
      <SimpleGrid cols={{ base: 1, sm: 2, lg: 3 }}>
        {data?.map((chal) => (
          <Card key={chal.id} shadow="sm" padding="lg" radius="md" withBorder>
            <Group justify="space-between" mb="xs">
              <Text fw={500}>{chal.title}</Text>
              <Badge color={getCategoryColor(chal.category)} size="sm">
                {chal.category}
              </Badge>
            </Group>
            <Text size="sm" c="dimmed" lineClamp={2}>{chal.description}</Text>
            <Group justify="space-between" mt="md">
              <Text size="sm" fw={600}>{chal.points} pts</Text>
              <Text size="sm" c="dimmed">{chal.solves} solves</Text>
            </Group>
          </Card>
        ))}
      </SimpleGrid>
    </Stack>
  )
}

function getCategoryColor(cat: string) {
  const colors: Record<string, string> = {
    web: 'blue', pwn: 'red', crypto: 'violet', re: 'orange',
    misc: 'gray', forensics: 'cyan', osint: 'green',
  }
  return colors[cat] || 'gray'
}

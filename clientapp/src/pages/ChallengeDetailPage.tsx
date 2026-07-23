import { useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { Card, Text, Badge, Group, Stack, Title, TextInput, Button, Code } from '@mantine/core'
import { useForm } from '@mantine/form'
import { notifications } from '@mantine/notifications'
import { api } from '../../utils/api'
import type { Challenge } from '../../types'

export function ChallengeDetailPage() {
  const { id } = useParams<{ id: string }>()

  const { data: challenge, isLoading } = useQuery({
    queryKey: ['challenge', id],
    queryFn: async () => {
      const res = await api.get(`/challenges/${id}`)
      return res.data as Challenge
    },
  })

  const form = useForm({ initialValues: { flag: '' } })

  const handleSubmit = async (values: { flag: string }) => {
    try {
      const res = await api.post('/submit', { challenge_id: id, flag: values.flag })
      notifications.show({
        title: res.data.correct ? 'Correct!' : 'Incorrect',
        message: res.data.message,
        color: res.data.correct ? 'green' : 'red',
      })
      form.reset()
    } catch (err: any) {
      notifications.show({ title: 'Error', message: err.response?.data?.message, color: 'red' })
    }
  }

  if (isLoading || !challenge) return <Text>Loading...</Text>

  return (
    <Stack>
      <Group justify="space-between">
        <Title order={2}>{challenge.title}</Title>
        <Group>
          <Badge color="blue" size="lg">{challenge.category}</Badge>
          <Badge color="green" size="lg">{challenge.points} pts</Badge>
        </Group>
      </Group>

      <Text c="dimmed">{challenge.solves} solves</Text>

      <Card withBorder p="md">
        <Text style={{ whiteSpace: 'pre-wrap' }}>{challenge.description}</Text>
      </Card>

      {challenge.container_image && (
        <Card withBorder p="md">
          <Text fw={600} mb="sm">Environment</Text>
          <Code>Container: {challenge.container_image}</Code>
          <Button mt="md" variant="light">Start Container</Button>
        </Card>
      )}

      <Card withBorder p="md">
        <Text fw={600} mb="sm">Submit Flag</Text>
        <form onSubmit={form.onSubmit(handleSubmit)}>
          <Group>
            <TextInput
              placeholder="flag{...}"
              style={{ flex: 1 }}
              {...form.getInputProps('flag')}
              required
            />
            <Button type="submit">Submit</Button>
          </Group>
        </form>
      </Card>
    </Stack>
  )
}

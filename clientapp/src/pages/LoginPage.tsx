import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { TextInput, PasswordInput, Button, Paper, Title, Text, Stack } from '@mantine/core'
import { useForm } from '@mantine/form'
import { notifications } from '@mantine/notifications'
import { api } from '../utils/api'
import { useAuthStore } from '../stores/auth'

export function LoginPage() {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const setToken = useAuthStore((s) => s.setToken)
  const setUser = useAuthStore((s) => s.setUser)

  const form = useForm({
    initialValues: { username: '', password: '' },
    validate: {
      username: (v) => (v.length < 3 ? 'Username must be at least 3 characters' : null),
      password: (v) => (v.length < 8 ? 'Password must be at least 8 characters' : null),
    },
  })

  const handleSubmit = async (values: typeof form.values) => {
    setLoading(true)
    try {
      const { data } = await api.post('/auth/login', values)
      setToken(data.token)
      setUser(data.user)
      notifications.show({ title: 'Success', message: 'Logged in!', color: 'green' })
      navigate('/')
    } catch (err: any) {
      notifications.show({
        title: 'Error',
        message: err.response?.data?.message || 'Login failed',
        color: 'red',
      })
    } finally {
      setLoading(false)
    }
  }

  return (
    <Paper p="xl" maw={400} mx="auto" mt="xl">
      <Title order={2} ta="center" mb="md">Login</Title>
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack>
          <TextInput label="Username" {...form.getInputProps('username')} required />
          <PasswordInput label="Password" {...form.getInputProps('password')} required />
          <Button type="submit" loading={loading}>Login</Button>
          <Text size="sm" ta="center" c="dimmed">
            Don't have an account? <a href="/register">Register</a>
          </Text>
        </Stack>
      </form>
    </Paper>
  )
}

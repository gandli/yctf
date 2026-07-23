import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { TextInput, PasswordInput, Button, Paper, Title, Text, Stack } from '@mantine/core'
import { useForm } from '@mantine/form'
import { notifications } from '@mantine/notifications'
import { api } from '../../utils/api'
import { useAuthStore } from '../../stores/auth'

export function RegisterPage() {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const setToken = useAuthStore((s) => s.setToken)
  const setUser = useAuthStore((s) => s.setUser)

  const form = useForm({
    initialValues: { username: '', email: '', password: '' },
    validate: {
      username: (v) => (v.length < 3 ? 'Username must be at least 3 characters' : null),
      email: (v) => (!/\S+@\S+\.\S+/.test(v) ? 'Invalid email' : null),
      password: (v) => (v.length < 8 ? 'Password must be at least 8 characters' : null),
    },
  })

  const handleSubmit = async (values: typeof form.values) => {
    setLoading(true)
    try {
      const { data } = await api.post('/auth/register', values)
      setToken(data.token)
      setUser(data.user)
      notifications.show({ title: 'Success', message: 'Account created!', color: 'green' })
      navigate('/')
    } catch (err: any) {
      notifications.show({
        title: 'Error',
        message: err.response?.data?.message || 'Registration failed',
        color: 'red',
      })
    } finally {
      setLoading(false)
    }
  }

  return (
    <Paper p="xl" maw={400} mx="auto" mt="xl">
      <Title order={2} ta="center" mb="md">Register</Title>
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack>
          <TextInput label="Username" {...form.getInputProps('username')} required />
          <TextInput label="Email" {...form.getInputProps('email')} required />
          <PasswordInput label="Password" {...form.getInputProps('password')} required />
          <Button type="submit" loading={loading}>Register</Button>
          <Text size="sm" ta="center" c="dimmed">
            Already have an account? <a href="/login">Login</a>
          </Text>
        </Stack>
      </form>
    </Paper>
  )
}

import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

export default function Login() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await login(email, password)
      navigate('/')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Login failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center px-4 pt-20 pb-12">
      <div className="w-full max-w-md">
        <div className="bg-stark-surface border border-stark-border p-8">
          <div className="section-tag justify-center">Authentication</div>
          <h1 className="text-3xl font-display font-extrabold text-stark-text text-center mb-8 tracking-[-0.03em]">
            Welcome Back
          </h1>

          {error && (
            <div className="bg-red-500/10 border border-red-500/30 text-red-400 font-mono text-sm p-3 mb-6">{error}</div>
          )}

          <form onSubmit={handleSubmit} className="space-y-5">
            <div>
              <label className="block font-mono text-xs text-stark-muted uppercase tracking-[0.1em] mb-2">Email</label>
              <input
                type="email"
                value={email}
                onChange={e => setEmail(e.target.value)}
                required
                placeholder="you@example.com"
                className="w-full bg-stark-bg border border-stark-border text-stark-text font-mono text-sm p-3 focus:outline-none focus:border-stark-accent transition-colors placeholder:text-stark-muted/50"
              />
            </div>
            <div>
              <label className="block font-mono text-xs text-stark-muted uppercase tracking-[0.1em] mb-2">Password</label>
              <input
                type="password"
                value={password}
                onChange={e => setPassword(e.target.value)}
                required
                placeholder="Enter your password"
                className="w-full bg-stark-bg border border-stark-border text-stark-text font-mono text-sm p-3 focus:outline-none focus:border-stark-accent transition-colors placeholder:text-stark-muted/50"
              />
            </div>
            <button type="submit" disabled={loading} className="btn-primary w-full justify-center">
              {loading ? 'Signing in...' : 'Sign In'}
            </button>
          </form>

          <p className="font-mono text-xs text-stark-muted text-center mt-6">
            Don't have an account?{' '}
            <a href="/signup" onClick={(e) => { e.preventDefault(); navigate('/signup'); }} className="text-stark-accent hover:text-stark-accent-2">
              Create one
            </a>
          </p>
        </div>
      </div>
    </div>
  )
}

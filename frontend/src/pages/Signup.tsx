import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

export default function Signup() {
  const { signup } = useAuth()
  const navigate = useNavigate()
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [acceptedTerms, setAcceptedTerms] = useState(false)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    if (password !== confirmPassword) {
      setError('Passwords do not match')
      return
    }
    if (password.length < 8) {
      setError('Password must be at least 8 characters')
      return
    }
    if (!acceptedTerms) {
      setError('You must accept the Privacy Policy and Terms of Use')
      return
    }
    setLoading(true)
    try {
      await signup(username, email, password, acceptedTerms)
      navigate('/')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Signup failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center px-4 pt-20 pb-12">
      <div className="w-full max-w-md">
        <div className="bg-stark-surface border border-stark-border p-8">
          <div className="section-tag justify-center">Get Started</div>
          <h1 className="text-3xl font-display font-extrabold text-stark-text text-center mb-8 tracking-[-0.03em]">
            Create Account
          </h1>

          {error && (
            <div className="bg-red-500/10 border border-red-500/30 text-red-400 font-mono text-sm p-3 mb-6">{error}</div>
          )}

          <form onSubmit={handleSubmit} className="space-y-5">
            <div>
              <label className="block font-mono text-xs text-stark-muted uppercase tracking-[0.1em] mb-2">Username</label>
              <input
                type="text"
                value={username}
                onChange={e => setUsername(e.target.value)}
                required
                minLength={3}
                maxLength={30}
                placeholder="Your username"
                className="w-full bg-stark-bg border border-stark-border text-stark-text font-mono text-sm p-3 focus:outline-none focus:border-stark-accent transition-colors placeholder:text-stark-muted/50"
              />
            </div>
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
                minLength={8}
                placeholder="Min. 8 characters"
                className="w-full bg-stark-bg border border-stark-border text-stark-text font-mono text-sm p-3 focus:outline-none focus:border-stark-accent transition-colors placeholder:text-stark-muted/50"
              />
            </div>
            <div>
              <label className="block font-mono text-xs text-stark-muted uppercase tracking-[0.1em] mb-2">Confirm Password</label>
              <input
                type="password"
                value={confirmPassword}
                onChange={e => setConfirmPassword(e.target.value)}
                required
                placeholder="Confirm your password"
                className="w-full bg-stark-bg border border-stark-border text-stark-text font-mono text-sm p-3 focus:outline-none focus:border-stark-accent transition-colors placeholder:text-stark-muted/50"
              />
            </div>
            <label className="flex items-start gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={acceptedTerms}
                onChange={e => setAcceptedTerms(e.target.checked)}
                className="mt-1.5 accent-stark-accent"
              />
              <span className="font-mono text-xs text-stark-muted leading-relaxed">
                I have read and agree to the{' '}
                <Link to="/privacy" className="text-stark-accent hover:text-stark-accent-2 underline">Privacy Policy</Link>
                {' '}and{' '}
                <Link to="/terms" className="text-stark-accent hover:text-stark-accent-2 underline">Terms of Use</Link>.
                I understand that violation of these terms may result in account suspension.
              </span>
            </label>
            <button type="submit" disabled={loading} className="btn-primary w-full justify-center">
              {loading ? 'Creating account...' : 'Create Account'}
            </button>
          </form>

          <p className="font-mono text-xs text-stark-muted text-center mt-6">
            Already have an account?{' '}
            <a href="/login" onClick={(e) => { e.preventDefault(); navigate('/login'); }} className="text-stark-accent hover:text-stark-accent-2">
              Sign in
            </a>
          </p>
        </div>
      </div>
    </div>
  )
}

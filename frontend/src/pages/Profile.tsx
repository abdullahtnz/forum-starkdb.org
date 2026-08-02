import { useState } from 'react'
import { api } from '../api'
import { useAuth } from '../contexts/AuthContext'
import ProtectedRoute from '../components/auth/ProtectedRoute'
import type { Post } from '../types'

export default function Profile() {
  return (
    <ProtectedRoute>
      <ProfilePage />
    </ProtectedRoute>
  )
}

function ProfilePage() {
  const { user, refreshUser } = useAuth()
  const [username, setUsername] = useState(user?.username || '')
  const [bio, setBio] = useState(user?.bio || '')
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [saving, setSaving] = useState(false)

  const handleUpdate = async () => {
    if (!username.trim() || username.length < 3) {
      setError('Username must be at least 3 characters')
      return
    }
    setSaving(true)
    setError('')
    setMessage('')
    try {
      await api.put('/users/me', { username, bio })
      await refreshUser()
      setMessage('Profile updated successfully')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Update failed')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="min-h-screen pt-20 pb-12">
      <div className="max-w-3xl mx-auto px-4 sm:px-8">
        <div className="card mb-6">
          <div className="section-tag">Settings</div>
          <h1 className="text-3xl font-display font-extrabold text-stark-text tracking-[-0.03em] mb-8">Profile Settings</h1>

          {message && (
            <div className="bg-green-500/10 border border-green-500/30 text-green-400 font-mono text-sm p-3 mb-6">{message}</div>
          )}
          {error && (
            <div className="bg-red-500/10 border border-red-500/30 text-red-400 font-mono text-sm p-3 mb-6">{error}</div>
          )}

          <div className="space-y-5">
            <div>
              <label className="block font-mono text-xs text-stark-muted uppercase tracking-[0.1em] mb-2">Username</label>
              <input
                type="text"
                value={username}
                onChange={e => setUsername(e.target.value)}
                minLength={3}
                maxLength={30}
                className="w-full bg-stark-bg border border-stark-border text-stark-text font-mono text-sm p-3 focus:outline-none focus:border-stark-accent transition-colors"
              />
            </div>
            <div>
              <label className="block font-mono text-xs text-stark-muted uppercase tracking-[0.1em] mb-2">Email</label>
              <input
                type="text"
                value={user?.email || ''}
                disabled
                className="w-full bg-stark-bg border border-stark-border text-stark-muted font-mono text-sm p-3 opacity-50"
              />
            </div>
            <div>
              <label className="block font-mono text-xs text-stark-muted uppercase tracking-[0.1em] mb-2">Bio</label>
              <textarea
                value={bio}
                onChange={e => setBio(e.target.value)}
                rows={3}
                placeholder="Tell the community about yourself..."
                className="w-full bg-stark-bg border border-stark-border text-stark-text font-mono text-sm p-3 focus:outline-none focus:border-stark-accent transition-colors placeholder:text-stark-muted/50 resize-y"
              />
            </div>
            <button onClick={handleUpdate} disabled={saving} className="btn-primary">
              {saving ? 'Saving...' : 'Save Changes'}
            </button>
          </div>
        </div>

        <div className="card">
          <h2 className="font-mono text-sm text-stark-accent uppercase tracking-[0.1em] mb-2">Account Info</h2>
          <div className="font-mono text-xs text-stark-muted space-y-1">
            <p>Member since: {user ? new Date(user.created_at).toLocaleDateString() : 'N/A'}</p>
            <p>Last login: {user?.last_login ? new Date(user.last_login).toLocaleDateString() : 'N/A'}</p>
            <p>Role: {user?.is_admin ? 'Administrator' : 'Member'}</p>
          </div>
        </div>
      </div>
    </div>
  )
}

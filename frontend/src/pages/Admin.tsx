import { useState, useEffect } from 'react'
import { api } from '../api'
import { useAuth } from '../contexts/AuthContext'
import type { UserPublic } from '../types'

export default function Admin() {
  const { user } = useAuth()
  const [activeTab, setActiveTab] = useState<'users' | 'badwords'>('users')
  const [users, setUsers] = useState<UserPublic[]>([])
  const [usersTotal, setUsersTotal] = useState(0)
  const [badWords, setBadWords] = useState<string[]>([])
  const [newWord, setNewWord] = useState('')
  const [page, setPage] = useState(1)

  if (!user?.is_admin) {
    return (
      <div className="min-h-screen flex items-center justify-center pt-20">
        <div className="text-center">
          <div className="section-tag justify-center">Access Denied</div>
          <h1 className="text-3xl font-display font-extrabold text-stark-text mb-4">Admin Only</h1>
          <p className="font-mono text-sm text-stark-muted">You do not have administrator privileges.</p>
        </div>
      </div>
    )
  }

  const loadUsers = async () => {
    try {
      const data = await api.get<{ users: UserPublic[]; total_count: number; page: number }>(`/admin/users?page=${page}&page_size=20`)
      setUsers(data.users || [])
      setUsersTotal(data.total_count || 0)
    } catch {}
  }

  const loadBadWords = async () => {
    try {
      const data = await api.get<{ words: string[] }>('/admin/bad-words')
      setBadWords(data.words || [])
    } catch {}
  }

  useEffect(() => {
    if (activeTab === 'users') loadUsers()
    else loadBadWords()
  }, [activeTab, page])

  const handleDeleteUser = async (userId: string) => {
    if (!confirm('Delete this user?')) return
    try {
      await api.delete(`/admin/users/${userId}`)
      loadUsers()
    } catch {}
  }

  const handleToggleAdmin = async (userId: string, makeAdmin: boolean) => {
    try {
      await api.put(`/admin/users/${userId}/admin`, { make_admin: makeAdmin })
      loadUsers()
    } catch {}
  }

  const handleAddBadWord = async () => {
    if (!newWord.trim()) return
    try {
      await api.post('/admin/bad-words', { word: newWord })
      setNewWord('')
      loadBadWords()
    } catch {}
  }

  const handleRemoveBadWord = async (word: string) => {
    try {
      await api.delete('/admin/bad-words', false)
      await fetch('/api/admin/bad-words', {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
        },
        body: JSON.stringify({ word }),
      })
      loadBadWords()
    } catch {}
  }

  return (
    <div className="min-h-screen pt-20 pb-12">
      <div className="max-w-5xl mx-auto px-4 sm:px-8">
        <div className="section-tag">Administration</div>
        <h1 className="text-3xl font-display font-extrabold text-stark-text tracking-[-0.03em] mb-8">Admin Panel</h1>

        <div className="flex gap-0 border-b border-stark-border mb-6">
          <button
            onClick={() => setActiveTab('users')}
            className={`font-mono text-xs px-6 py-3 uppercase tracking-[0.08em] border-b-2 transition-colors ${
              activeTab === 'users' ? 'text-stark-accent border-stark-accent' : 'text-stark-muted border-transparent hover:text-stark-accent'
            }`}
          >
            Users
          </button>
          <button
            onClick={() => setActiveTab('badwords')}
            className={`font-mono text-xs px-6 py-3 uppercase tracking-[0.08em] border-b-2 transition-colors ${
              activeTab === 'badwords' ? 'text-stark-accent border-stark-accent' : 'text-stark-muted border-transparent hover:text-stark-accent'
            }`}
          >
            Bad Words
          </button>
        </div>

        {activeTab === 'users' && (
          <div>
            <div className="card overflow-x-auto">
              <table className="w-full font-mono text-xs">
                <thead>
                  <tr className="border-b border-stark-border">
                    <th className="text-left p-3 text-stark-accent uppercase tracking-[0.1em]">Username</th>
                    <th className="text-left p-3 text-stark-accent uppercase tracking-[0.1em]">Role</th>
                    <th className="text-left p-3 text-stark-accent uppercase tracking-[0.1em]">Joined</th>
                    <th className="text-right p-3 text-stark-accent uppercase tracking-[0.1em]">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {users.map(u => (
                    <tr key={u.id} className="border-b border-stark-border hover:bg-white/[0.02]">
                      <td className="p-3">
                        <div className="flex items-center gap-2">
                          <div className="w-6 h-6 rounded bg-stark-accent/10 border border-stark-accent/20 flex items-center justify-center text-[0.6rem] text-stark-accent font-bold">
                            {u.username[0]?.toUpperCase()}
                          </div>
                          <span className="text-stark-text">{u.username}</span>
                        </div>
                      </td>
                      <td className="p-3">
                        <span className={`text-xs px-2 py-0.5 ${u.is_admin ? 'text-stark-accent bg-stark-accent/10 border border-stark-accent/20' : 'text-stark-muted bg-stark-surface-2 border border-stark-border'}`}>
                          {u.is_admin ? 'Admin' : 'Member'}
                        </span>
                      </td>
                      <td className="p-3 text-stark-muted">{new Date(u.created_at).toLocaleDateString()}</td>
                      <td className="p-3 text-right">
                        <div className="flex gap-2 justify-end">
                          {!u.is_admin && (
                            <button
                              onClick={() => handleToggleAdmin(u.id, true)}
                              className="font-mono text-[0.6rem] text-stark-accent hover:text-stark-accent-2 uppercase tracking-[0.08em] transition-colors"
                            >
                              Make Admin
                            </button>
                          )}
                          {u.is_admin && u.id !== user?.id && (
                            <button
                              onClick={() => handleToggleAdmin(u.id, false)}
                              className="font-mono text-[0.6rem] text-stark-accent-3 hover:text-stark-accent uppercase tracking-[0.08em] transition-colors"
                            >
                              Remove Admin
                            </button>
                          )}
                          {u.id !== user?.id && !u.is_admin && (
                            <button
                              onClick={() => handleDeleteUser(u.id)}
                              className="font-mono text-[0.6rem] text-red-400 hover:text-stark-accent-3 uppercase tracking-[0.08em] transition-colors"
                            >
                              Delete
                            </button>
                          )}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            {usersTotal > 20 && (
              <div className="flex justify-center gap-4 mt-4">
                <button disabled={page <= 1} onClick={() => setPage(p => p - 1)} className="btn-ghost text-xs py-2 px-4 disabled:opacity-30">Previous</button>
                <button disabled={page >= Math.ceil(usersTotal / 20)} onClick={() => setPage(p => p + 1)} className="btn-ghost text-xs py-2 px-4 disabled:opacity-30">Next</button>
              </div>
            )}
          </div>
        )}

        {activeTab === 'badwords' && (
          <div>
            <div className="card mb-6">
              <h2 className="font-mono text-xs text-stark-accent uppercase tracking-[0.1em] mb-4">Add Bad Word</h2>
              <div className="flex gap-3">
                <input
                  type="text"
                  value={newWord}
                  onChange={e => setNewWord(e.target.value)}
                  placeholder="Enter word to filter..."
                  className="flex-1 bg-stark-bg border border-stark-border text-stark-text font-mono text-sm p-3 focus:outline-none focus:border-stark-accent transition-colors"
                />
                <button onClick={handleAddBadWord} className="btn-primary text-sm py-3 px-6">Add</button>
              </div>
            </div>

            <div className="card">
              <h2 className="font-mono text-xs text-stark-accent uppercase tracking-[0.1em] mb-4">Filtered Words ({badWords.length})</h2>
              {badWords.length === 0 ? (
                <p className="font-mono text-xs text-stark-muted">No bad words configured.</p>
              ) : (
                <div className="flex flex-wrap gap-2">
                  {badWords.map(word => (
                    <span key={word} className="flex items-center gap-2 font-mono text-xs px-3 py-1.5 bg-stark-surface-2 border border-stark-border text-stark-muted">
                      {word}
                      <button onClick={() => handleRemoveBadWord(word)} className="text-stark-accent-3 hover:text-stark-accent transition-colors">×</button>
                    </span>
                  ))}
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

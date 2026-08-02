import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { api } from '../api'
import { useAuth } from '../contexts/AuthContext'
import type { Post, Tag, PostListResponse } from '../types'

function timeAgo(date: string) {
  const seconds = Math.floor((Date.now() - new Date(date).getTime()) / 1000)
  if (seconds < 60) return 'just now'
  const mins = Math.floor(seconds / 60)
  if (mins < 60) return `${mins}m ago`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}h ago`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days}d ago`
  const months = Math.floor(days / 30)
  if (months < 12) return `${months}mo ago`
  return `${Math.floor(months / 12)}y ago`
}

export default function Home() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [posts, setPosts] = useState<Post[]>([])
  const [tags, setTags] = useState<Tag[]>([])

  if (!Array.isArray(posts)) return null
  if (!Array.isArray(tags)) return null
  const [search, setSearch] = useState('')
  const [activeTag, setActiveTag] = useState('')
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const pageSize = 20

  useEffect(() => {
    api.get<Tag[]>('/tags').then(data => setTags(data || [])).catch(() => {})
  }, [])

  useEffect(() => {
    setLoading(true)
    const params = new URLSearchParams()
    params.set('page', String(page))
    params.set('page_size', String(pageSize))
    params.set('sort', 'newest')
    if (search) params.set('q', search)
    if (activeTag) params.set('tag', activeTag)

        api.get<PostListResponse>(`/posts?${params.toString()}`)
          .then(data => {
            setPosts(data.posts || []);
            setTotal(data.total_count);
          })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [page, search, activeTag])

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    setPage(1)
  }

  return (
    <div className="min-h-screen pt-20">
      <div className="max-w-5xl mx-auto px-4 sm:px-8">
        <div className="page-header text-center pb-8 border-b border-stark-border">
          <div className="section-tag justify-center">Community</div>
          <h1 className="text-4xl sm:text-5xl font-display font-extrabold text-stark-text tracking-[-0.04em] mb-4">
            STARKDB <span className="text-stark-accent">Forum</span>
          </h1>
          <p className="font-mono text-sm text-stark-muted max-w-xl mx-auto">
            Ask questions, share solutions, and discuss the lightweight B-Tree key-value database.
          </p>
          {user && (
            <button
              onClick={() => navigate('/ask')}
              className="btn-primary mt-6 text-sm"
            >
              Ask Question
            </button>
          )}
        </div>

        <div className="py-6">
          <form onSubmit={handleSearch} className="flex gap-3 mb-6">
            <input
              type="text"
              value={search}
              onChange={e => { setSearch(e.target.value); setPage(1); }}
              placeholder="Search posts..."
              className="flex-1 bg-stark-surface border border-stark-border text-stark-text font-mono text-sm p-3 focus:outline-none focus:border-stark-accent transition-colors placeholder:text-stark-muted/50"
            />
            <button type="submit" className="btn-primary text-sm py-3 px-6">
              Search
            </button>
          </form>

          {tags.length > 0 && (
            <div className="flex gap-2 flex-wrap mb-6">
              <button
                onClick={() => { setActiveTag(''); setPage(1); }}
                className={`font-mono text-xs px-3 py-1.5 border transition-colors ${
                  activeTag === '' ? 'bg-stark-accent/10 border-stark-accent/30 text-stark-accent' : 'bg-stark-surface border-stark-border text-stark-muted hover:border-stark-accent/30 hover:text-stark-accent'
                }`}
              >
                All
              </button>
              {tags.map(tag => (
                <button
                  key={tag.id}
                  onClick={() => { setActiveTag(activeTag === tag.slug ? '' : tag.slug); setPage(1); }}
                  className={`font-mono text-xs px-3 py-1.5 border transition-colors ${
                    activeTag === tag.slug ? 'bg-stark-accent/10 border-stark-accent/30 text-stark-accent' : 'bg-stark-surface border-stark-border text-stark-muted hover:border-stark-accent/30 hover:text-stark-accent'
                  }`}
                >
                  {tag.name}
                </button>
              ))}
            </div>
          )}

          {loading ? (
            <div className="text-center py-20 text-stark-muted font-mono text-sm">Loading...</div>
          ) : posts.length === 0 ? (
            <div className="text-center py-20">
              <p className="text-stark-muted font-mono text-sm mb-4">No posts found.</p>
              {user && (
                <button onClick={() => navigate('/ask')} className="btn-primary text-sm">
                  Ask the first question
                </button>
              )}
            </div>
          ) : (
            <div className="grid gap-4">
              {posts.map(post => (
                <Link
                  key={post.id}
                  to={`/posts/${post.id}`}
                  className="card block hover:border-stark-accent/20 transition-all group"
                >
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex-1 min-w-0">
                      <h2 className="text-lg font-display font-bold text-stark-text group-hover:text-stark-accent transition-colors mb-2 tracking-[-0.01em] line-clamp-2">
                        {post.is_pinned && <span className="text-stark-accent mr-2">[PINNED]</span>}
                        {post.title}
                      </h2>
                      <p className="font-mono text-xs text-stark-muted leading-relaxed line-clamp-2 mb-3">
                        {post.content.replace(/<[^>]*>/g, '').substring(0, 300)}
                      </p>
                      <div className="flex flex-wrap items-center gap-2">
                        {post.tags?.map(tag => (
                          <span key={tag.id} className="font-mono text-[0.6rem] px-2 py-0.5 bg-stark-accent/5 border border-stark-accent/15 text-stark-accent uppercase tracking-[0.1em]">
                            {tag.name}
                          </span>
                        ))}
                      </div>
                    </div>
                    <div className="flex items-center gap-4 text-right shrink-0">
                      <div className="text-center min-w-[40px]">
                        <div className="font-display font-bold text-stark-text text-lg">{post.like_count}</div>
                        <div className="font-mono text-[0.6rem] text-stark-muted uppercase tracking-[0.1em]">likes</div>
                      </div>
                      <div className="text-center min-w-[40px]">
                        <div className="font-display font-bold text-stark-text text-lg">{post.comment_count}</div>
                        <div className="font-mono text-[0.6rem] text-stark-muted uppercase tracking-[0.1em]">replies</div>
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center justify-between mt-4 pt-3 border-t border-stark-border">
                    <div className="flex items-center gap-2">
                      <div className="w-6 h-6 rounded bg-stark-accent/10 border border-stark-accent/20 flex items-center justify-center text-[0.6rem] font-mono text-stark-accent font-bold">
                        {post.username[0]?.toUpperCase()}
                      </div>
                      <span className="font-mono text-xs text-stark-muted">{post.username}</span>
                    </div>
                    <span className="font-mono text-xs text-stark-muted">{timeAgo(post.created_at)}</span>
                  </div>
                </Link>
              ))}
            </div>
          )}

          {total > pageSize && (
            <div className="flex justify-center items-center gap-4 mt-8 pb-8">
              <button
                disabled={page <= 1}
                onClick={() => setPage(p => p - 1)}
                className="btn-ghost text-xs py-2 px-4 disabled:opacity-30 disabled:cursor-not-allowed"
              >
                Previous
              </button>
              <span className="font-mono text-xs text-stark-muted">
                Page {page} of {Math.ceil(total / pageSize)}
              </span>
              <button
                disabled={page >= Math.ceil(total / pageSize)}
                onClick={() => setPage(p => p + 1)}
                className="btn-ghost text-xs py-2 px-4 disabled:opacity-30 disabled:cursor-not-allowed"
              >
                Next
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

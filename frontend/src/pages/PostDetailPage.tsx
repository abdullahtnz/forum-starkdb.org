import { useState, useEffect } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { api } from '../api'
import { useAuth } from '../contexts/AuthContext'
import CommentItem from '../components/comments/CommentItem'
import type { Post, Comment } from '../types'

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

function formatSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

export default function PostDetailPage() {
  const { id } = useParams<{ id: string }>()
  const { user } = useAuth()
  const navigate = useNavigate()
  const [post, setPost] = useState<Post | null>(null)
  const [comments, setComments] = useState<Comment[]>([])
  const [commentText, setCommentText] = useState('')
  const [loading, setLoading] = useState(true)
  const [commentLoading, setCommentLoading] = useState(false)
  const [liked, setLiked] = useState(false)
  const [likeCount, setLikeCount] = useState(0)

  const loadPost = async () => {
    if (!id) return
    try {
      const data = await api.get<Post>(`/posts/${id}`)
      setPost(data)
      setLiked(data.liked)
      setLikeCount(data.like_count)
    } catch {
      navigate('/')
    } finally {
      setLoading(false)
    }
  }

  const loadComments = async () => {
    if (!id) return
    try {
      const data = await api.get<Comment[]>(`/posts/${id}/comments`)
      setComments(data || [])
    } catch {}
  }

  useEffect(() => { loadPost(); loadComments(); }, [id])

  const handleLike = async () => {
    if (!user || !id) return
    try {
      const res = await api.post<{ liked: boolean }>(`/posts/${id}/like`, {})
      setLiked(res.liked)
      setLikeCount(c => res.liked ? c + 1 : c - 1)
    } catch {}
  }

  const handleComment = async () => {
    if (!user || !id || !commentText.trim()) return
    setCommentLoading(true)
    try {
      await api.post(`/posts/${id}/comments`, { content: commentText })
      setCommentText('')
      loadComments()
    } catch {} finally {
      setCommentLoading(false)
    }
  }

  const handleDelete = async () => {
    if (!id || !confirm('Delete this post?')) return
    try {
      await api.delete(`/posts/${id}`)
      navigate('/')
    } catch {}
  }

  if (loading) return <div className="min-h-screen flex items-center justify-center text-stark-muted font-mono text-sm">Loading...</div>
  if (!post) return null

  return (
    <div className="min-h-screen pt-20 pb-12">
      <div className="max-w-4xl mx-auto px-4 sm:px-8">
        <div className="card mb-6">
          <div className="flex flex-wrap items-center gap-2 mb-3">
            {post.tags?.map(tag => (
              <span key={tag.id} className="font-mono text-[0.6rem] px-2 py-0.5 bg-stark-accent/5 border border-stark-accent/15 text-stark-accent uppercase tracking-[0.1em]">
                {tag.name}
              </span>
            ))}
            {post.is_pinned && <span className="font-mono text-[0.6rem] px-2 py-0.5 bg-stark-accent-3/10 border border-stark-accent-3/20 text-stark-accent-3 uppercase tracking-[0.1em]">Pinned</span>}
            {post.is_closed && <span className="font-mono text-[0.6rem] px-2 py-0.5 bg-stark-accent-3/10 border border-stark-accent-3/20 text-stark-accent-3 uppercase tracking-[0.1em]">Closed</span>}
          </div>

          <h1 className="text-2xl sm:text-3xl font-display font-extrabold text-stark-text tracking-[-0.03em] mb-4">
            {post.title}
          </h1>

          <div className="flex items-center gap-3 mb-6 pb-4 border-b border-stark-border">
            <Link to={`/users/${post.user_id}`} className="flex items-center gap-2 hover:text-stark-accent transition-colors">
              <div className="w-7 h-7 rounded bg-stark-accent/10 border border-stark-accent/20 flex items-center justify-center text-xs font-mono text-stark-accent font-bold">
                {post.username[0]?.toUpperCase()}
              </div>
              <span className="font-mono text-sm text-stark-text">{post.username}</span>
            </Link>
            <span className="font-mono text-xs text-stark-muted">· {timeAgo(post.created_at)}</span>
            <span className="font-mono text-xs text-stark-muted">· {post.view_count} views</span>
          </div>

          <div className="prose prose-invert max-w-none font-mono text-sm leading-relaxed text-stark-muted mb-6" dangerouslySetInnerHTML={{ __html: post.content }} />

          {post.attachments && post.attachments.length > 0 && (
            <div className="mb-6 p-4 bg-stark-surface-2 border border-stark-border">
              <h3 className="font-mono text-xs text-stark-accent uppercase tracking-[0.1em] mb-3">Attachments</h3>
              <div className="flex flex-wrap gap-3">
                {post.attachments.map(att => (
                  <a
                    key={att.id}
                    href={att.file_path}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-2 font-mono text-xs text-stark-muted bg-stark-bg border border-stark-border p-2 hover:border-stark-accent/30 hover:text-stark-accent transition-all"
                  >
                    <span>⬇</span>
                    <div>
                      <div className="text-stark-text">{att.file_name}</div>
                      <div className="text-[0.6rem]">{formatSize(att.file_size)}</div>
                    </div>
                  </a>
                ))}
              </div>
            </div>
          )}

          <div className="flex items-center justify-between pt-4 border-t border-stark-border">
            <div className="flex items-center gap-4">
              <button
                onClick={handleLike}
                disabled={!user}
                className={`btn-ghost text-sm py-2 px-4 ${liked ? 'border-stark-accent text-stark-accent' : ''}`}
              >
                {liked ? '♥' : '♡'} {likeCount}
              </button>
              <span className="font-mono text-xs text-stark-muted">{post.comment_count} comments</span>
            </div>
            <div className="flex gap-2">
              {user && post.is_owner && (
                <>
                  <button onClick={handleDelete} className="font-mono text-xs text-stark-accent-3 hover:text-stark-accent transition-colors uppercase tracking-[0.08em]">
                    Delete
                  </button>
                </>
              )}
            </div>
          </div>
        </div>

        <div className="card mb-6">
          <h2 className="font-mono text-sm text-stark-accent uppercase tracking-[0.1em] mb-4">
            {comments.length} {comments.length === 1 ? 'Comment' : 'Comments'}
          </h2>

          {user && !post.is_closed && (
            <div className="mb-6 flex gap-3">
              <input
                type="text"
                value={commentText}
                onChange={e => setCommentText(e.target.value)}
                placeholder="Write a comment..."
                className="flex-1 bg-stark-bg border border-stark-border text-stark-text font-mono text-sm p-3 focus:outline-none focus:border-stark-accent transition-colors placeholder:text-stark-muted/50"
              />
              <button onClick={handleComment} disabled={commentLoading || !commentText.trim()} className="btn-primary text-sm py-3 px-6">
                {commentLoading ? '...' : 'Post'}
              </button>
            </div>
          )}

          {post.is_closed && (
            <div className="mb-6 p-3 bg-stark-accent-3/5 border border-stark-accent-3/20 text-stark-accent-3 font-mono text-xs">
              This post is closed to new comments.
            </div>
          )}

          {comments.length === 0 && !user && (
            <p className="font-mono text-xs text-stark-muted">
              No comments yet.{' '}
              <Link to="/login" className="text-stark-accent hover:text-stark-accent-2">Login</Link>
              {' '}to add the first comment.
            </p>
          )}

          {comments.length === 0 && user && !post.is_closed && (
            <p className="font-mono text-xs text-stark-muted">Be the first to comment.</p>
          )}

          <div className="space-y-1">
            {comments.map(comment => (
              <CommentItem key={comment.id} comment={comment} postId={post.id} onRefresh={loadComments} />
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}
